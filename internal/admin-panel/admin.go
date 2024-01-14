package adminpanel

import (
	"fmt"
	"reflect"

	"github.com/dmawardi/Go-Template/internal/config"
)

// init state for db access
var app *config.AppConfig

// Function called in main.go to connect app state to current file
func SetStateInAdminPanel(a *config.AppConfig) {
	app = a
}

// Build item list for sidebar (Add for every module)
var sidebar = AdminSideBar{
	Main: []sidebarItem{
		// This list is filled upon runtime by GenerateAndSetAdminSidebar
	},
	Auth: BuildAuthSidebarSection(),
}

// Default Records Displayed on find all pages
var recordsPerPage = []int{10, 25, 50, 100}

// Admin controller (used in API)
type AdminController struct {
	Base AdminBaseController
	User AdminUserController
	Post AdminPostController
	Auth AdminAuthPolicyController
}

// Constructor
func NewAdminController(base AdminBaseController, users AdminUserController, posts AdminPostController, authPolicies AdminAuthPolicyController) AdminController {
	return AdminController{base, users, posts, authPolicies}
}

// Used for rendering admin sidebar
type sidebarItem struct {
	Name        string
	FindAllLink string
	AddLink     string
}

type AdminSideBar struct {
	Main []sidebarItem
	Auth []sidebarItem
}

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

// Constructor
func NewBasicAdminController(url, schema, schemaPlural string, obtainFields func() BasicAdminController) BasicAdminController {
	return basicAdminController{url, schema, schemaPlural}
}

func (b basicAdminController) ObtainFields() basicAdminController {
	return basicAdminController{b.AdminHomeUrl, b.SchemaName, b.PluralSchemaName}
}

// Interface for all schemas (used for Admin panel) (Add receiver functions for every schema)
type AdminPanelSchema interface {
	// Returns ID of record
	GetID() string
	// Returns value of schema field
	ObtainValue(keyValue string) string
}

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

// Generate and set sidebar list
// Accepts current Admin controller and generates sidebar list based on controllers
func GenerateAndSetAdminSidebar(adminCont AdminController) {
	// Iterate through all controllers and add to sidebar list
	// Get the reflect.Value of the struct.
	valueOfCont := reflect.ValueOf(adminCont)

	// Iterate through the struct fields.
	for i := 0; i < valueOfCont.NumField(); i++ {
		// Get the field name and value.
		fieldName := valueOfCont.Type().Field(i).Name
		fieldValue := valueOfCont.Field(i).Interface()

		// If not base controller, add to sidebar list
		if fieldName != "Base" && fieldName != "Auth" {
			currentController := ObtainFieldsForAnyType(fieldValue)
			// Create sidebar item
			item := sidebarItem{
				Name:        currentController.PluralSchemaName,
				AddLink:     fmt.Sprintf("%s/create", currentController.AdminHomeUrl),
				FindAllLink: currentController.AdminHomeUrl,
			}

			// append to sidebar list
			sidebar.Main = append(sidebar.Main, item)
		}
	}
}

func BuildAuthSidebarSection() []sidebarItem {
	return []sidebarItem{
		{
			Name:        "Policies",
			FindAllLink: "/admin/policy",
			AddLink:     "/admin/policy/create",
		},
		{
			Name:        "Roles",
			FindAllLink: "/admin/policy/roles",
			AddLink:     "/admin/policy/create-role",
		},
		{
			Name:        "Inheritance",
			FindAllLink: "/admin/policy/inheritance",
			AddLink:     "/admin/policy/create-inheritance",
		},
	}
}
