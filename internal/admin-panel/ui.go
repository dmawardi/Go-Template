package adminpanel

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
)

// PAGE RENDER DATA
// Contains state for the page
type PageRenderData struct {
	// In HEAD
	PageTitle string
	// In BODY
	SectionTitle  string
	SectionDetail template.HTML
	SidebarList   AdminSideBar
	// Schema home used to return to the schema home page from delete
	SchemaHome string // eg. /admin/users/
	// Page type (Used for content selection)
	PageType PageType
	// Form
	FormData  FormData
	TableData TableData
	// Search
	SearchTerm             string
	RecordsPerPageSelector []int
	// Special section data for policies
	PolicySection PolicySection
	HeaderSection HeaderSection
}

// Variables for header
type HeaderSection struct {
	HomeUrl           template.URL
	ViewSiteUrl       template.URL
	ChangePasswordUrl template.URL
	LogOutUrl         template.URL
}

// Variables for policy section
type PolicySection struct {
	FocusedPolicies []PolicyEditDataRow
	PolicyResource  string
	Selectors       PolicyEditSelectors
}

// Page type (Used for dynamic selective rendering)
type PageType struct {
	HomePage    bool
	EditPage    bool
	ReadPage    bool
	CreatePage  bool
	DeletePage  bool
	SuccessPage bool
	// Used for policy section
	PolicyMode string // eg. "policy" or "inheritance"
}

// TEMPLATES
//
// Parses all the template files in the templates directory
func ParseAdminTemplates() (*template.Template, error) {
	// Parse the base template
	tmpl := template.New("/internal/admin-panel/templates/layout.tmpl")

	// Walk through all files in the templates directory
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// If the file is not a directory and has the .html extension
		if !info.IsDir() && filepath.Ext(path) == ".tmpl" {
			// Parse the file
			_, err = tmpl.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})

	// Return error if there is filepath walk issue
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// Function to render the Admin error page to the response
func serveAdminError(w http.ResponseWriter, sectionTitle string) {
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "Error - Admin",
		SectionTitle: sectionTitle,
		SidebarList:  sidebar,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   false,
			CreatePage: false,
			DeletePage: true,
		},
		FormData: FormData{},
	}

	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// Function to render the Admin success page to the response
func serveAdminSuccess(w http.ResponseWriter, pageTitle string, sectionTitle string) {
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    pageTitle,
		SectionTitle: sectionTitle,
		SidebarList:  sidebar,
		PageType: PageType{
			EditPage:    false,
			ReadPage:    false,
			CreatePage:  false,
			DeletePage:  false,
			SuccessPage: true,
		},
		FormData:      FormData{},
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// SIDEBAR
//
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
