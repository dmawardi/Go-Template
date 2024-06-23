package modules

import (
	adminpanel "github.com/dmawardi/Go-Template/internal/admin-panel"
	"github.com/dmawardi/Go-Template/internal/models"
	"gorm.io/gorm"
)

// Returns module map that contains structs with modules (controller, service, & repo) using module name as key
func SetupModules(modulesToSetup []EntityConfig, client *gorm.DB, selectorService adminpanel.SelectorService) models.ModuleMap {
	// Init
	moduleMap := make(map[string]models.ModuleSet)

	for _, module := range modulesToSetup {
		// Create repo, service, and controller using the client
		repo := module.NewRepo(client)
		service := module.NewService(repo)
		controller := module.NewController(service)

		// Assign admin controller
		adminController := module.NewAdminController
		// If admin controller is not nil, add it to the module map
		if adminController != nil {

			// Create admin controller using the service
			adminController := adminController(service, selectorService)

			// Add module set including admin controller to the map
			moduleMap[module.Name] = models.ModuleSet{
				Repo:            repo,
				Service:         service,
				Controller:      controller,
				AdminController: adminController,
			}
		} else {
			// Add module set without admin controller to the map
			moduleMap[module.Name] = models.ModuleSet{
				Repo:            repo,
				Service:         service,
				Controller:      controller,
				AdminController: nil,
			}
		}
	}
	return moduleMap
}

// Used in API setup to standardize the array of setup configurations
type EntityConfig struct {
	Name               string
	NewRepo            func(*gorm.DB) interface{}
	NewService         func(interface{}) interface{}
	NewController      func(interface{}) interface{}
	NewAdminController func(interface{}, interface{}) interface{}
}
