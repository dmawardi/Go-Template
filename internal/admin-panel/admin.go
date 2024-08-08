package adminpanel

import (
	"html/template"

	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/models"
	coreservices "github.com/dmawardi/Go-Template/internal/service/core"
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
	// Action service for recording admin actions
	Action coreservices.ActionService
	// Basic modules
	Base AdminCoreController
	User AdminUserController
	Auth AdminAuthPolicyController
	// Additional modules contained in module map
	ModuleMap models.ModuleMap
}

// Constructor
func NewAdminPanelController(action coreservices.ActionService, 
							base AdminCoreController, 
							users AdminUserController, 
							authPolicies AdminAuthPolicyController, 
							moduleMap models.ModuleMap) AdminPanelController {
	return AdminPanelController{action, base, users, authPolicies, moduleMap}
}



