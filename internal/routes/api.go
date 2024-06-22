package routes

import (
	"net/http"

	adminpanel "github.com/dmawardi/Go-Template/internal/admin-panel"
	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/controller/core"
	modulecontrollers "github.com/dmawardi/Go-Template/internal/controller/moduleControllers"
)

// Create new service repository
func BuildRouteState(a *config.AppConfig) {
	app = a
}

type Api interface {
	// Total route builder for API
	Routes(crudRouteSet CRUDRouteSet, adminRouteSet AdminRouteSet) http.Handler
}

// Api that contains all controllers for route creation
type api struct {
	Admin  adminpanel.AdminController
	User   core.UserController
	Policy core.AuthPolicyController
	Post   modulecontrollers.PostController
}

func NewApi(
	admin adminpanel.AdminController,
	user core.UserController,
	policy core.AuthPolicyController,
	post modulecontrollers.PostController) Api {
	return &api{Admin: admin, User: user, Policy: policy, Post: post}
}
