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
	"sync"
	"time"
)

var (
	status     = &types.AppStatus{LastUpdate: time.Now()}
	templates  = template.Must(template.ParseFiles("templates/index.html", "templates/dashboard.html"))
	updateLock = sync.Mutex{}
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

	err := templates.ExecuteTemplate(w, "index.html", status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	status.Mu.RLock()
	defer status.Mu.RUnlock()

	err := templates.ExecuteTemplate(w, "dashboard.html", status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleTimestamp(w http.ResponseWriter, r *http.Request) {
	status.Mu.RLock()
	defer status.Mu.RUnlock()

	// Create a template for just the timestamp
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

func pingServers(apps []types.Application, dryRun bool) ([]types.Application, error) {
	var results []types.Application
	for _, app := range apps {
		app.Status = make(map[string]bool)
		for _, server := range app.Servers {
			if dryRun {
				app.Status[server] = (rand.Intn(2)+1)%2 == 0
			} else {
				reqUrl := fmt.Sprintf("%s/%s", server, app.TestUrl)
				req, err := http.Get(reqUrl)

				app.Status[server] = !(err != nil && req.StatusCode != http.StatusOK)
			}
		}
		results = append(results, app)
	}

	return results, nil
}

// startStatusChecker runs the pingServers function in a goroutine
func startStatusChecker(apps []types.Application, interval time.Duration) {
	isDryRun := config.GetDryRunConfig()
	go func() {
		for {
			updateLock.Lock()
			updatedApps, _ := pingServers(apps, isDryRun)

			status.Mu.Lock()
			status.Apps = updatedApps
			status.LastUpdate = time.Now()
			status.Mu.Unlock()

			updateLock.Unlock()

			time.Sleep(interval)
		}
	}()
}
