package service

import (
	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/models"
	coreservices "github.com/dmawardi/Go-Template/internal/service/core"
	moduleservices "github.com/dmawardi/Go-Template/internal/service/module"
)

// Repository used by handler package
var app *config.AppConfig

// Create new service repository
func SetAppConfig(a *config.AppConfig) {
	// Set app state in core services
	coreservices.SetAppConfig(a)
	// Set app state in module services
	moduleservices.SetAppConfig(a)

	app = a
}

// Basic Module Service Interface
type BasicModuleService interface {
	FindAll(limit int, offset int, order string, conditions []models.QueryConditionParameters) (*models.BasicPaginatedResponse, error)
	FindById(int) (*struct{}, error)
	Create(entity *struct{}) (*struct{}, error)
	Update(int, *struct{}) (*struct{}, error)
	Delete(int) error
	BulkDelete([]int) error
}