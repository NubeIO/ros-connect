package model

type ServiceFile struct {
	Name                        string   `json:"name"`
	Version                     string   `json:"version"`
	ServiceDescription          string   `json:"service_description"`
	RunAsUser                   string   `json:"run_as_user"`
	ServiceWorkingDirectory     string   `json:"service_working_directory"`        // /data/installer/apps/install/rubix-os/v0.6.1/
	ExecStart                   string   `json:"exec_start"`                       // app -p 1660 -g <data_dir> -d data -prod
	AttachWorkingDirOnExecStart bool     `json:"attach_working_dir_on_exec_start"` // true, false
	EnvironmentVars             []string `json:"environment_vars"`                 // Environment="g=/data/bacnet-server-c"
}

type InstalledApps struct {
	AppName              string  `json:"appName,omitempty"`
	Version              string  `json:"version,omitempty"`
	MinVersion           string  `json:"minVersion,omitempty"`
	MaxVersion           string  `json:"maxVersion,omitempty"`
	ServiceName          string  `json:"serviceName,omitempty"`
	IsInstalled          bool    `json:"isInstalled"`
	Message              string  `json:"message,omitempty"`
	Match                bool    `json:"match,omitempty"`
	DowngradeRequired    bool    `json:"downgradeRequired,omitempty"`
	UpgradeRequired      bool    `json:"upgradeRequired,omitempty"`
	State                string  `json:"state,omitempty"`
	ActiveState          string  `json:"activeState,omitempty"`
	SubState             string  `json:"subState,omitempty"`
	ActiveEnterTimestamp string  `json:"activeEnterTimestamp"`
	RestartExpression    *string `json:"restartExpression,omitempty"`
}

type AppsAvailableForInstall struct {
	AppName     string `json:"appName,omitempty"`
	MinVersion  string `json:"minVersion,omitempty"`
	MaxVersion  string `json:"maxVersion,omitempty"`
	Description string `json:"description,omitempty"`
}

type RunningServices struct {
	Name                 string `json:"name,omitempty"`
	ServiceName          string `json:"serviceName,omitempty"`
	State                string `json:"state,omitempty"`
	ActiveState          string `json:"activeState,omitempty"`
	SubState             string `json:"subState,omitempty"`
	ActiveEnterTimestamp string `json:"activeEnterTimestamp"`
}

type EdgeAppsInfo struct {
	InstalledApps           []InstalledApps           `json:"installedApps"`
	AppsAvailableForInstall []AppsAvailableForInstall `json:"appsAvailableForInstall"`
	RunningServices         []RunningServices         `json:"runningServices"`
}
