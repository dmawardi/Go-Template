package service

import (
	"github.com/dmawardi/Go-Template/internal/config"
	coreservices "github.com/dmawardi/Go-Template/internal/service/core"
)

// Repository used by handler package
var app *config.AppConfig

// Create new service repository
func SetAppConfig(a *config.AppConfig) {
	// Set app state in core services
	coreservices.SetAppConfig(a)
	app = a
}
