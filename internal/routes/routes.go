package routes

import (
	"fmt"
	"net/http"

	adminpanel "github.com/dmawardi/Go-Template/internal/admin-panel"
	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/controller"
	"github.com/dmawardi/Go-Template/internal/controller/core"
	modulecontrollers "github.com/dmawardi/Go-Template/internal/controller/moduleControllers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

var app *config.AppConfig

// Create new service repository
func BuildRouteState(a *config.AppConfig) {
	app = a
}

type Api interface {
	// Total route builder for API
	Routes() http.Handler
	// Api routes
	AddUserApiRoutes(router *chi.Mux) *chi.Mux
	AddBasicCrudApiRoutes(router *chi.Mux, urlExtension string, controller controller.BasicController) *chi.Mux
	// Admin routes
	AddBasicAdminRoutes(router *chi.Mux, controller adminpanel.AdminBaseController) *chi.Mux
	AddAdminRouteSet(router *chi.Mux, protected bool, urlExtension string, controller controller.BasicAdminController) *chi.Mux
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

// Overall Routes builder for server
func (a api) Routes() http.Handler {
	// Create new router
	mux := chi.NewRouter()
	// Use built in Chi middleware
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)
	mux.Use(corsMiddleware)

	// Add special user and group routes
	mux = a.AddUserApiRoutes(mux)
	mux = a.AddAuthRBACApiRoutes(mux)
	// Other schemas
	mux = a.AddBasicCrudApiRoutes(mux, "posts", a.Post)

	// Add basic admin routes (home, login, etc)
	mux = a.AddBasicAdminRoutes(mux, a.Admin.Base)
	// Add admin policy routes
	mux = a.AddAdminPolicySet(mux, true, "policy", a.Admin.Auth)

	// Add admin panel schema route sets
	mux = a.AddAdminRouteSet(mux, false, "users", a.Admin.User)
	mux = a.AddAdminRouteSet(mux, false, "posts", a.Admin.Post)

	// Serve API Swagger docs at built URL from config state
	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s/static/docs/swagger.json", app.BaseURL)), //The url pointing to API definition
	))
	fmt.Printf("Serving Swagger docs at http://%s/swagger/index.html\n", app.BaseURL)

	// Build fileserver using static directory
	fileServer := http.FileServer(http.Dir("./static"))
	// Handle all calls to /static/* by stripping prefix and sending to file server
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return mux
}
