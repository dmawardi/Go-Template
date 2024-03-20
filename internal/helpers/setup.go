package helpers

import (
	"github.com/dmawardi/Go-Template/internal/models"
	"gorm.io/gorm"
)

// MODULE STANDARDIZATION
//

// Returns module set for each basic module to be set up with key as module name
func SetupBasicModules(basicModulesToSetup []models.EntityConfig, client *gorm.DB) map[string]models.ModuleSet {
	modules := make(map[string]models.ModuleSet)

	for _, config := range basicModulesToSetup {
		repo := config.NewRepo(client)
		service := config.NewService(repo)
		controller := config.NewController(service)

		// Add controllers to a map for later use (e.g., in building the API)
		// Use type or name as a key as appropriate for your application

		modules[config.Name] = models.ModuleSet{Repo: repo, Service: service, Controller: controller}
	}
	return modules
}
