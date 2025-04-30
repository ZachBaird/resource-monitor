package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"resource-monitor/config"
	"resource-monitor/types"
	"resource-monitor/utils"
	"sync"
	"time"
)

var (
	status     = &types.AppStatus{LastUpdate: time.Now()}
	templates  = template.Must(template.ParseFiles("templates/index.html", "templates/dashboard.html"))
	updateLock = sync.Mutex{}
	httpClient = http.Client{}
)

func main() {
	f, err := os.Open("./apps.json")
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		os.Exit(1)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v", err)
		}
	}(f)

	apps, err := parseJsonFile(f)
	if err != nil {
		fmt.Printf("Failure parsing file: %v", err)
		os.Exit(1)
	}

	status.Apps = apps
	startStatusChecker(apps, 30*time.Second)

	// Set up HTTP routes
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/dashboard", handleDashboard)
	http.HandleFunc("/timestamp", handleTimestamp)
	http.HandleFunc("/style.css", handleCSS)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	status.Mu.RLock()
	defer status.Mu.RUnlock()

	utils.ServeTemplateFile("index.html", w, status, templates)
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	status.Mu.RLock()
	defer status.Mu.RUnlock()

	utils.ServeTemplateFile("dashboard.html", w, status, templates)
}

func handleTimestamp(w http.ResponseWriter, r *http.Request) {
	status.Mu.RLock()
	defer status.Mu.RUnlock()

	// timestamp template
	tmpl := template.Must(template.New("timestamp").Parse("{{.Format \"Jan 02, 2006 15:04:05\"}}"))
	err := tmpl.Execute(w, status.LastUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	http.ServeFile(w, r, "./static/style.css")
}

func parseJsonFile(f *os.File) ([]types.Application, error) {
	var apps types.Applications

	bytes, _ := io.ReadAll(f)
	err := json.Unmarshal(bytes, &apps)
	if err != nil {
		return nil, err
	}

	return apps.Applications, nil
}

func execDryRun(apps []types.Application) []types.Application {
	var results []types.Application
	for _, app := range apps {
		app.Status = make(map[string]bool)
		for _, server := range app.Servers {
			app.Status[server] = (rand.Intn(2)+1)%2 == 0
		}
		results = append(results, app)
	}

	return results
}

func pingServers(apps []types.Application, secrets map[string]string) []types.Application {
	var results []types.Application
	for _, app := range apps {
		app.Status = make(map[string]bool)
		for _, server := range app.Servers {
			reqUrl := fmt.Sprintf("%s/%s", server, app.TestUrl)
			req, err := http.NewRequest(app.HttpMethod, reqUrl, nil)
			if err != nil {
				app.Status[server] = false
				continue
			}

			if app.ApiKey != "" {
				authHeader := "Authorization"
				if app.Header != "" {
					authHeader = app.Header
				}
				req.Header.Set(authHeader, app.ApiKey)
			}
			res, err := httpClient.Do(req)

			app.Status[server] = !(err != nil && res.StatusCode != http.StatusOK)

		}
		results = append(results, app)
	}

	return results
}

// startStatusChecker runs the pingServers function in a goroutine
func startStatusChecker(apps []types.Application, interval time.Duration) {
	isDryRun := config.GetDryRunConfig()
	appSecrets := make(map[string]string)

	for _, app := range apps {
		appSecrets[app.Name] = config.GetSecretConfig(app.Name)
	}
	go func() {
		for {
			updateLock.Lock()
			status.Mu.Lock()
			if isDryRun {
				status.Apps = execDryRun(apps)
			} else {
				status.Apps = pingServers(apps, appSecrets)
			}

			status.LastUpdate = time.Now()
			status.Mu.Unlock()

			updateLock.Unlock()

			time.Sleep(interval)
		}
	}()
}
