package adminpanel

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dmawardi/Go-Template/internal/config"
)

// init state for db access
var app *config.AppConfig

// Build schema item list for sidebar
var sidebarList = []string{"Users", "Groups"}

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
