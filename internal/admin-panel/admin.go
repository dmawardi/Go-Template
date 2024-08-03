package adminpanel

import (
	"html/template"
	"reflect"

	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/models"
)

// init state for db access
var app *config.AppConfig

// Links for admin header
var header = HeaderSection{
	HomeUrl:           "#",
	ViewSiteUrl:       "#",
	LogOutUrl:         "#",
	ChangePasswordUrl: "#",
}

// Function called in main.go to connect app state to current file
func SetStateInAdminPanel(a *config.AppConfig) {
	app = a
	// Set header urls after setting state
	header.HomeUrl = template.URL("http://" + app.BaseURL + "/admin/home")
	header.ChangePasswordUrl = template.URL("http://" + app.BaseURL + "/admin/change-password")
	header.LogOutUrl = template.URL("http://" + app.BaseURL + "/admin/logout")
	header.ViewSiteUrl = template.URL("http://" + app.BaseURL + "/swagger/index.html")
}

// Admin controller (used in API)
type AdminPanelController struct {
	// Basic modules
	Base AdminCoreController
	User AdminUserController
	Auth AdminAuthPolicyController
	// Additional modules contained in module map
	ModuleMap models.ModuleMap
}

// Constructor
func NewAdminPanelController(base AdminCoreController, users AdminUserController, authPolicies AdminAuthPolicyController, moduleMap models.ModuleMap) AdminPanelController {
	return AdminPanelController{base, users, authPolicies, moduleMap}
}

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

// Interface for all schemas that makes it compatible with admin panel (Add receiver functions for every schema)
type AdminPanelSchema interface {
	// Returns ID of record
	GetID() string
	// Returns value of schema field
	ObtainValue(keyValue string) string
}
