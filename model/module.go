package model

type Module struct {
	Name          string   `json:"name"`
	Arch          string   `json:"arch"`
	Version       string   `json:"version"`
	IsInstalled   bool     `json:"isInstalled"`
	Repo          string   `json:"repo"`
	Description   string   `json:"description"`
	Products      []string `json:"products"`
	AppDependency []string `json:"appDependency"`
}
