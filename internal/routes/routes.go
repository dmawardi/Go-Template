package routes

import (
	"fmt"
	"net/http"

	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/controller"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

var app *config.AppConfig

type CRUDRouteSet struct {
	Controller controller.BasicController
	Name       string
}

type AdminRouteSet struct {
	Controller controller.BasicAdminController
	Name       string
}

// Overall Routes builder for server
func (a api) Routes(crudRouteSet CRUDRouteSet, adminRouteSet AdminRouteSet) http.Handler {
	// Create new router
	mux := chi.NewRouter()
	// Use built in Chi middleware
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)
	mux.Use(corsMiddleware)

	// Add user and group API routes
	mux = a.AddUserApiRoutes(mux)
	mux = a.AddAuthRBACApiRoutes(mux)

	// Add basic admin panel routes (home, login, etc)
	mux = a.AddBasicAdminRoutes(mux, a.Admin.Base)
	// Add admin user routes
	mux = a.AddAdminRouteSet(mux, false, "users", a.Admin.User)
	// Add admin policy routes
	mux = a.AddAdminPolicySet(mux, true, "policy", a.Admin.Auth)

	// Other schemas
	mux = a.AddBasicCrudApiRoutes(mux, "posts", a.Post)
	// Add admin panel schema route sets
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
