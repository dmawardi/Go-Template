package adminpanel

import (
	"html/template"
	"reflect"

	"github.com/dmawardi/Go-Template/internal/config"
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
type AdminController struct {
	// Basic modules
	Base AdminBaseController
	User AdminUserController
	Auth AdminAuthPolicyController
	// ADD BASIC MODULES HERE
	Post AdminPostController
}

// Constructor
func NewAdminController(base AdminBaseController, users AdminUserController, authPolicies AdminAuthPolicyController, posts AdminPostController) AdminController {
	return AdminController{base, users, authPolicies, posts}
}

// Interface for all basic admin controllers (used for Admin panel to dynamically generate sidebar)
type BasicAdminController interface {
	ObtainFields() basicAdminController
}
type basicAdminController struct {
	// For links
	AdminHomeUrl string
	// For HTML text rendering
	SchemaName       string
	PluralSchemaName string
}

// RECEIVER FUNCTIONS
func (b basicAdminController) ObtainFields() basicAdminController {
	return basicAdminController{b.AdminHomeUrl, b.SchemaName, b.PluralSchemaName}
}

// ADMIN SIDEBAR CREATION
//
// Uses the ObtainFields method to get the fields of any Basic Admin Controller type
func ObtainFieldsForAnyType(input interface{}) basicAdminController {
	// Use reflection to call ObtainFields method if it exists.
	value := reflect.ValueOf(input)
	// ObtainFields method
	method := value.MethodByName("ObtainFields")
	if !method.IsValid() {
		return basicAdminController{}
	}

	// Call ObtainFields method
	result := method.Call(nil)
	// Check if result is valid
	if len(result) == 1 {

		resultFields := result[0].Interface()

		controllerFields := basicAdminController{
			AdminHomeUrl:     resultFields.(basicAdminController).AdminHomeUrl,
			SchemaName:       resultFields.(basicAdminController).SchemaName,
			PluralSchemaName: resultFields.(basicAdminController).PluralSchemaName,
		}
		return controllerFields
	}

	return basicAdminController{}
}

// Interface for all schemas (used for Admin panel) (Add receiver functions for every schema)
type AdminPanelSchema interface {
	// Returns ID of record
	GetID() string
	// Returns value of schema field
	ObtainValue(keyValue string) string
}
