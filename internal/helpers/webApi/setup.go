package webapi

import (
	"bytes"
	"html/template"

	"github.com/dmawardi/Go-Template/internal/models"
	"gorm.io/gorm"
)

// MODULE STANDARDIZATION
//

// Helper function to create a new repository. Takes a repository creation function and returns a function that takes a database connection and returns an interface
func NewRepository[T any](repoFunc func(*gorm.DB) T) func(*gorm.DB) interface{} {
	return func(db *gorm.DB) interface{} {
		return repoFunc(db)
	}
}

// Helper function to create a new service. Takes a service creation function and returns a function that takes an interface and returns an interface
func NewService[T any, S any](serviceFunc func(T) S) func(interface{}) interface{} {
	return func(repoInterface interface{}) interface{} {
		repo, ok := repoInterface.(T)
		if !ok {
			panic("Incorrect repository type")
		}
		return serviceFunc(repo)
	}
}

// Helper function to create a new controller. Takes a controller creation function and returns a function that takes an interface and returns an interface
func NewController[T any, C any](controllerFunc func(T) C) func(interface{}) interface{} {
	return func(serviceInterface interface{}) interface{} {
		service, ok := serviceInterface.(T)
		if !ok {
			panic("Incorrect service type")
		}
		return controllerFunc(service)
	}
}

// Returns module set (controller, service, & repo) to be set up with key as module name
func SetupModules(modulesToSetup []models.EntityConfig, client *gorm.DB) map[string]models.ModuleSet {
	moduleMap := make(map[string]models.ModuleSet)

	for _, module := range modulesToSetup {
		repo := module.NewRepo(client)
		service := module.NewService(repo)
		controller := module.NewController(service)

		// Add controllers to a map for later use (e.g., in building the API)
		// Use name as a key
		moduleMap[module.Name] = models.ModuleSet{Repo: repo, Service: service, Controller: controller}
	}
	return moduleMap
}

// LoadTemplate parses an HTML template, executes it with the provided data, and returns the result as a string.
func LoadTemplate(templateFilePath string, data interface{}) (string, error) {
	// Parse the template file
	t, err := template.ParseFiles(templateFilePath)
	if err != nil {
		return "", err
	}

	// Build the template with the injected data
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
