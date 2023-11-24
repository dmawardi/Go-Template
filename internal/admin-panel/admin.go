package adminpanel

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/go-chi/chi"
)

// init state for db access
var app *config.AppConfig

// Build admin api controller
var adminApi = NewAdminApiController(NewAdminBaseController(), NewUserAdminController())

// Build schema item list for sidebar
var sidebarList = []string{"Users", "Groups"}

// Function to add new routes to an existing Chi mux router
func AddAdminRoutes(router *chi.Mux) *chi.Mux {
	// Admin routes
	// router.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("This is the admin login page"))
	// })
	// Public routes
	router.Group(func(mux chi.Router) {
		// @tag.name Public Routes
		// @tag.description Unprotected routes
		mux.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("This is the admin login page"))
		})
		// admin users
		// Read (all users)
		mux.Get("/admin/users", adminApi.users.FindAll)
		// Create (GET form / POST form)
		mux.Get("/admin/users/create", adminApi.users.Create)
		mux.Post("/admin/users/create", adminApi.users.Create)
		// Delete
		mux.Post("/admin/users/delete", adminApi.users.Delete)
		// Edit/Update (GET data in form / POST form)
		mux.Get("/admin/users/{id}", adminApi.users.Edit)
		mux.Post("/admin/users/{id}", adminApi.users.Edit)

		// Private routes
		mux.Group(func(mux chi.Router) {
			mux.Use(auth.AuthenticateJWT)

			// @tag.name Private routes
			// @tag.description Protected routes
			// admin home
			mux.Get("/admin/home", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("This is the admin main home page"))
			})

		})

	})

	return router
}

// Admin home controller
type AdminBaseController interface {
	Home(w http.ResponseWriter, r *http.Request)
}

type adminBaseController struct {
}

// Constructor
func NewAdminBaseController() AdminBaseController {
	return &adminBaseController{}
}

func (c adminBaseController) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the admin main home page"))
}

// Admin API controller
type AdminApiController struct {
	base  AdminBaseController
	users AdminUserController
}

// Decleare new admin api controller
func NewAdminApiController(base AdminBaseController, users AdminUserController) AdminApiController {
	return AdminApiController{base, users}
}

// Parses all the template files in the templates directory
func parseAdminTemplates() (*template.Template, error) {
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

// Function called in main.go to connect app state to current file
func SetStateInAdminPanel(a *config.AppConfig) {
	app = a
}
