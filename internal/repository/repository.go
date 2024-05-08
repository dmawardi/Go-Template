package repository

import (
	"github.com/dmawardi/Go-Template/internal/config"
	corerepositories "github.com/dmawardi/Go-Template/internal/repository/core"
)

var app *config.AppConfig

func SetAppConfig(appConfig *config.AppConfig) {
	app = appConfig
	// Set app config in repository
	corerepositories.SetAppConfig(app)
}
