package types

import (
	"sync"
	"time"
)

type Applications struct {
	Applications []Application `json:"applications"`
}
type Application struct {
	Name    string   `json:"appName"`
	Servers []string `json:"servers"`
	TestUrl string   `json:"testUrl"`
	Status  map[string]bool
}

type AppStatus struct {
	Apps       []Application
	LastUpdate time.Time
	Mu         sync.RWMutex
}
