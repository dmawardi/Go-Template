package adminpanel

import (
	"fmt"
	"reflect"
)

// Build item list for sidebar (Add for every module)
var sidebar = AdminSideBar{
	Main: []sidebarItem{
		// This list is filled upon runtime by GenerateAndSetAdminSidebar
	},
	Auth: BuildAuthSidebarSection(),
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

// Build auth section for sidebar in admin panel
func BuildAuthSidebarSection() []sidebarItem {
	return []sidebarItem{
		{
			Name:        "Permissions",
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
