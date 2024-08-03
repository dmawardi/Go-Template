package adminpanel

import "reflect"

// Interface for all basic admin controllers (used for Admin panel to dynamically generate sidebar)
type BasicAdminController interface {
	ObtainUrlDetails() basicAdminController
}
type basicAdminController struct {
	// For links
	AdminHomeUrl string
	// For HTML text rendering
	SchemaName       string
	PluralSchemaName string
}

// RECEIVER FUNCTIONS
func (b basicAdminController) ObtainUrlDetails() basicAdminController {
	return basicAdminController{b.AdminHomeUrl, b.SchemaName, b.PluralSchemaName}
}

// ADMIN SIDEBAR CREATION
//
// Uses the ObtainUrlDetails method to get the sidebar details of any Basic Admin Controller type
func ObtainUrlDetailsForBasicAdminController(input interface{}) basicAdminController {
	// Use reflection to call ObtainUrlDetails method if it exists.
	value := reflect.ValueOf(input)
	// ObtainUrlDetails method
	method := value.MethodByName("ObtainUrlDetails")
	if !method.IsValid() {
		return basicAdminController{}
	}

	// Call ObtainUrlDetails method
	result := method.Call(nil)

	// Check if result is valid (if it is a BasicAdminController)
	// If it has a result
	if len(result) == 1 {
		// Assign the result as an interface to resultFields
		interfaceFields := result[0].Interface()
		// Assign the fields of the resultFields to sidebarDetails
		sidebarDetails := basicAdminController{
			AdminHomeUrl:     interfaceFields.(basicAdminController).AdminHomeUrl,
			SchemaName:       interfaceFields.(basicAdminController).SchemaName,
			PluralSchemaName: interfaceFields.(basicAdminController).PluralSchemaName,
		}
		return sidebarDetails
	}

	return basicAdminController{}
}