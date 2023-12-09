package adminpanel

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dmawardi/Go-Template/internal/config"
)

// init state for db access
var app *config.AppConfig

// Build item list for sidebar (Add for every module)
var sidebarList = []sidebarItem{
	{Name: "Users", AddLink: "/admin/users/create", FindAllLink: "/admin/users"},
	{Name: "Groups", AddLink: "/admin/groups/create", FindAllLink: "/admin/groups"},
	{Name: "Posts", AddLink: "/admin/posts/create", FindAllLink: "/admin/posts"},
}

// Displayed on find all pages
var recordsPerPage = []int{10, 25, 50, 100}

// Function called in main.go to connect app state to current file
func SetStateInAdminPanel(a *config.AppConfig) {
	app = a
}

// Admin base controller (non-schema related routes)
type AdminBaseController interface {
	Home(w http.ResponseWriter, r *http.Request)
}
type adminBaseController struct {
}

// Constructor
func NewAdminBaseController() AdminBaseController {
	return &adminBaseController{}
}

// Admin controller (used in API)
type AdminController struct {
	Base AdminBaseController
	User AdminUserController
}

// Constructor
func NewAdminController(base AdminBaseController, users AdminUserController) AdminController {
	return AdminController{base, users}
}

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

// RECEIVER FUNCTIONS
func (c adminBaseController) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the admin main home page"))
}

// Function to render the Admin error page to the response
func serveAdminError(w http.ResponseWriter, sectionTitle string) {
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "Error - Admin",
		SectionTitle: sectionTitle,
		SidebarList:  sidebarList,
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

// Used for rendering admin sidebar
type sidebarItem struct {
	Name        string
	FindAllLink string
	AddLink     string
}

// Interface for all schemas (used for Admin panel) (Add for every schema)
type AdminPanelSchema interface {
	GetID() string
	ObtainValue(keyValue string) interface{}
}

// Map used to group form selectors for a schema (eg. FormSelector["field_name"])
type FormSelectors map[string]func() []FormFieldSelector
