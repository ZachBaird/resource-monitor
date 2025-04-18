package types

type Applications struct {
	Applications []Application `json:"applications"`
}
type Application struct {
	Name    string   `json:"appName"`
	Servers []string `json:"servers"`
	TestUrl string   `json:"testUrl"`
}
