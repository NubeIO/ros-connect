package model

type Plugin struct {
	UUID         string   `json:"uuid"`
	Name         string   `json:"name"`
	ModulePath   string   `json:"modulePath"`
	Author       string   `json:"author"`
	Website      string   `json:"website"`
	Enabled      bool     `json:"enabled"`
	HasNetwork   bool     `json:"hasNetwork"`
	Capabilities []string `json:"capabilities"`
}

type AvailablePlugin struct {
	Name        string `json:"name"`
	IsInstalled bool   `json:"isInstalled"`
	Description string `json:"description"`
}
