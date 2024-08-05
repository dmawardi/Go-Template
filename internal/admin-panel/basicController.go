package adminpanel

import (
	"reflect"

	"github.com/dmawardi/Go-Template/internal/service"
)

// Interface for all basic admin controllers (used for Admin panel to dynamically generate sidebar)
type BasicAdminController interface {
	ObtainUrlDetails() URLDetails
}
type basicAdminController struct {
	service service.BasicModuleService
	// For links
	AdminHomeUrl string
	// For HTML text rendering
	SchemaName       string
	PluralSchemaName string
	// Custom table headers
	tableHeaders  []TableHeader
	formSelectors SelectorService
}

// RECEIVER FUNCTIONS
func (b basicAdminController) ObtainUrlDetails() URLDetails {
	return URLDetails{
		AdminHomeUrl:     b.AdminHomeUrl,
		SchemaName:       b.SchemaName,
		PluralSchemaName: b.PluralSchemaName,
	}
}

// ADMIN SIDEBAR CREATION
//
// Uses the ObtainUrlDetails method to get the sidebar details of any Basic Admin Controller type
func ObtainUrlDetailsForBasicAdminController(input interface{}) URLDetails {
	// Use reflection to call ObtainUrlDetails method if it exists.
	value := reflect.ValueOf(input)
	// ObtainUrlDetails method
	method := value.MethodByName("ObtainUrlDetails")
	if !method.IsValid() {
		return URLDetails{}
	}

	// Call ObtainUrlDetails method
	result := method.Call(nil)

	// Check if result is valid (if it is a BasicAdminController)
	// If it has a result
	if len(result) == 1 {
		// Assign the result as an interface to resultFields
		interfaceFields := result[0].Interface()
		// Assign the fields of the resultFields to sidebarDetails
		sidebarDetails := URLDetails{
			AdminHomeUrl:     interfaceFields.(URLDetails).AdminHomeUrl,
			SchemaName:       interfaceFields.(URLDetails).SchemaName,
			PluralSchemaName: interfaceFields.(URLDetails).PluralSchemaName,
		}
		return sidebarDetails
	}

	return URLDetails{}
}

type URLDetails struct {
	AdminHomeUrl string
	SchemaName   string
	PluralSchemaName string
}