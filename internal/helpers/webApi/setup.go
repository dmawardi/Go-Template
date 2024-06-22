package webapi

import (
	"bytes"
	"html/template"

	"github.com/dmawardi/Go-Template/internal/models"
	"gorm.io/gorm"
)

// MODULE STANDARDIZATION
//

// Returns module set (controller, service, & repo) to be set up with key as module name
func SetupBasicModules(basicModulesToSetup []models.EntityConfig, client *gorm.DB) map[string]models.ModuleSet {
	moduleMap := make(map[string]models.ModuleSet)

	for _, module := range basicModulesToSetup {
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
