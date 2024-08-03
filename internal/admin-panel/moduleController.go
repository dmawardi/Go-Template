package adminpanel

import "github.com/dmawardi/Go-Template/internal/service"

// Interface for all module controllers
type adminModuleController struct {
	service service.BasicModuleService
	// For links
	adminHomeUrl string
	// For HTML text rendering
	schemaName       string
	pluralSchemaName string
	// Custom table headers
	tableHeaders  []TableHeader
	formSelectors SelectorService
}