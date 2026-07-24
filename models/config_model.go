package models

type CustomInstallation struct {
	Instructions []string `json:"instructions"`
	Installed    bool     `json:"installed"`
}

type ConfigFile struct {
	PacmanInstallations []string                      `json:"pacman"`
	CustomInstalations  map[string]CustomInstallation `json:"custom"`
}
