package types

import (
	"sync"
	"time"
)

type Applications struct {
	Applications []Application `json:"applications"`
}
type Application struct {
	Name       string   `json:"appName"`
	Servers    []string `json:"servers"`
	TestUrl    string   `json:"testUrl"`
	HttpMethod string   `json:"httpMethod"`
	ApiKey     string   `json:"apiKey"`
	Header     string   `json:"header"`

	Status map[string]bool
}

type AppStatus struct {
	Apps       []Application
	LastUpdate time.Time
	Mu         sync.RWMutex
}
