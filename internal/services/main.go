package services

import "github.com/dmawardi/Go-Template/internal/config"

// Repository used by handler package
var Repo *Repository

// Repository type
type Repository struct {
	App *config.AppConfig
}

// Create new service repository
func BuildServiceRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}
