package service

import "github.com/dmawardi/Go-Template/internal/config"

// Repository used by handler package
var app *config.AppConfig

// Create new service repository
func BuildServiceState(a *config.AppConfig) {
	app = a
}
