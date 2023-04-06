package service

import (
	"github.com/dmawardi/Go-Template/internal/config"
)

// Repository used by handler package
var app *config.AppConfig

// // Repository type
// type Repository struct {
// 	App *config.AppConfig
// }

// Create new service repository
func BuildServiceState(a *config.AppConfig) {
	app = a
}

type Services struct {
	User *UserService
}
