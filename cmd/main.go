package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"resource-monitor/types"
	"time"
)

func main() {
	f, err := os.Open("./apps.json")
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		os.Exit(1)
	}
	fmt.Println("Application configuration accessed..")

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

	http.HandleFunc("/", servePage("index.html"))
	//err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		results, _ := pingServers(apps, true)
		for _, r := range results {
			fmt.Printf("%v\n", r)
		}

		fmt.Printf("\n\n\n")
		time.Sleep(5 * time.Second)
	}
}

func servePage(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("html", filename))
	}
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

func pingServers(apps []types.Application, dryRun bool) ([]string, error) {
	results := make([]string, 0)
	for _, app := range apps {
		for _, server := range app.Servers {
			if dryRun {
				random := rand.Intn(2) + 1
				if random%2 == 0 {
					results = append(results, fmt.Sprintf("%s: %s is down", app.Name, server))
				}

				results = append(results, fmt.Sprintf("%s: %s", app.Name, server))
				continue
			}

			reqUrl := fmt.Sprintf("%s/%s", server, app.TestUrl)
			req, err := http.Get(reqUrl)
			if err != nil || req.StatusCode != http.StatusOK {
				results = append(results, fmt.Sprintf("%s: %s is down", app.Name, server))
			}

			results = append(results, fmt.Sprintf("%s: %s", app.Name, server))
		}
	}

	return results, nil
}
