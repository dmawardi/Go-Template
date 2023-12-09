package adminpanel

import (
	"github.com/dmawardi/Go-Template/internal/config"
)

// init state for db access
var app *config.AppConfig

// Function called in main.go to connect app state to current file
func SetStateInAdminPanel(a *config.AppConfig) {
	app = a
}

// Build item list for sidebar (Add for every module)
var sidebarList = []sidebarItem{
	{Name: "Users", AddLink: "/admin/users/create", FindAllLink: "/admin/users"},
	{Name: "Groups", AddLink: "/admin/groups/create", FindAllLink: "/admin/groups"},
	{Name: "Posts", AddLink: "/admin/posts/create", FindAllLink: "/admin/posts"},
}

// Default ReDisplayed on find all pages
var recordsPerPage = []int{10, 25, 50, 100}

// Admin controller (used in API)
type AdminController struct {
	Base AdminBaseController
	User AdminUserController
}

// Constructor
func NewAdminController(base AdminBaseController, users AdminUserController) AdminController {
	return AdminController{base, users}
}

// Used for rendering admin sidebar
type sidebarItem struct {
	Name        string
	FindAllLink string
	AddLink     string
}

// Interface for all schemas (used for Admin panel) (Add for every schema)
type AdminPanelSchema interface {
	// Returns ID of record
	GetID() string
	// Returns value of schema field
	ObtainValue(keyValue string) interface{}
}
