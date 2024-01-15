package routes

import (
	"fmt"
	"net/http"

	adminpanel "github.com/dmawardi/Go-Template/internal/admin-panel"
	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/controller"
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
	User   controller.UserController
	Post   controller.PostController
	Policy controller.AuthPolicyController
}

func NewApi(admin adminpanel.AdminController, user controller.UserController, policy controller.AuthPolicyController, post controller.PostController) Api {
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

	// // Add basic admin routes
	mux = a.AddBasicAdminRoutes(mux, a.Admin.Base)
	// mux = a.AddAdminCrudRoutes(mux, true, "users", a.Admin.User)
	mux = a.AddAdminRouteSet(mux, false, "users", a.Admin.User)
	mux = a.AddAdminRouteSet(mux, false, "posts", a.Admin.Post)
	mux = a.AddAdminPolicySet(mux, false, "policy", a.Admin.Auth)

	// Serve API Swagger docs at built URL from config state
	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s/static/docs/swagger.json", app.BaseURL)), //The url pointing to API definition
	))

	// Build fileserver using static directory
	fileServer := http.FileServer(http.Dir("./static"))
	// Handle all calls to /static/* by stripping prefix and sending to file server
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return mux
}

// Adds User routes to a Chi mux router (includes login, forgot password, etc)
func (a api) AddUserApiRoutes(router *chi.Mux) *chi.Mux {
	// Public routes
	router.Group(func(mux chi.Router) {
		// @tag.name Public Routes
		// @tag.description Unprotected routes
		mux.Get("/", controller.GetJobs)
		// Login
		mux.Post("/api/users/login", a.User.Login)
		// Forgot password
		mux.Post("/api/users/forgot-password", a.User.ResetPassword)
		// Verify Email
		mux.Get("/api/users/verify-email/{token}", a.User.EmailVerification)
		mux.Post("/api/users/send-verification-email", a.User.ResendVerificationEmail)

		// Create new user
		mux.Post("/api/users", a.User.Create)

		// Private routes
		mux.Group(func(mux chi.Router) {
			mux.Use(auth.AuthenticateJWT)

			// @tag.name Private routes
			// @tag.description Protected routes
			// users
			mux.Get("/api/users", a.User.FindAll)
			mux.Get("/api/users/{id}", a.User.Find)
			mux.Put("/api/users/{id}", a.User.Update)
			mux.Delete("/api/users/{id}", a.User.Delete)

			// My profile
			mux.Get("/api/me", a.User.GetMyUserDetails)
			mux.Post("/api/me", controller.HealthCheck)
			mux.Put("/api/me", a.User.UpdateMyProfile)
		})

	})
	return router
}

// Adds Authorization routes to a Chi mux router
func (a api) AddAuthRBACApiRoutes(router *chi.Mux) *chi.Mux {
	// Public routes
	router.Group(func(mux chi.Router) {
		// @tag.name Public Routes
		// @tag.description Unprotected routes

		// Private routes
		mux.Group(func(mux chi.Router) {
			mux.Use(auth.AuthenticateJWT)

			// @tag.name Private routes
			// @tag.description Protected routes
			// Auth
			mux.Get("/api/auth", a.Policy.FindAll)
			mux.Get("/api/auth/roles", a.Policy.FindAllRoles)
			mux.Put("/api/auth/roles", a.Policy.AssignUserRole)

			mux.Post("/api/auth", a.Policy.Create)
			mux.Put("/api/auth", a.Policy.Update)
			mux.Delete("/api/auth", a.Policy.Delete)
		})

	})
	return router
}

// Adds a basic fully authorized CRUD route set to a Chi mux router
func (a api) AddBasicCrudApiRoutes(router *chi.Mux, urlExtension string, controller controller.BasicController) *chi.Mux {
	// Public routes
	router.Group(func(mux chi.Router) {
		// Private routes
		mux.Use(auth.AuthenticateJWT)
		// @tag.name Private routes
		// @tag.description Protected routes
		// Route set
		mux.Get(fmt.Sprintf("/api/%s", urlExtension), controller.FindAll)
		mux.Get(fmt.Sprintf("/api/%s/{id}", urlExtension), controller.Find)
		mux.Put(fmt.Sprintf("/api/%s/{id}", urlExtension), controller.Update)
		mux.Post(fmt.Sprintf("/api/%s", urlExtension), controller.Create)
		mux.Delete(fmt.Sprintf("/api/%s/{id}", urlExtension), controller.Delete)
	})

	return router
}

// Adds a basic Admin CRUD route set to a Chi mux router
func (a api) AddAdminRouteSet(router *chi.Mux, protected bool, urlExtension string, controller controller.BasicAdminController) *chi.Mux {
	// Reassign for consistency
	r := router
	r.Group(func(mux chi.Router) {
		// Set to use JWT authentication if protected
		if protected {
			mux.Use(auth.AuthenticateJWT)
		}
		// Read (all users)
		mux.Get(fmt.Sprintf("/admin/%s", urlExtension), controller.FindAll)
		// Create (GET form / POST form)
		mux.Get(fmt.Sprintf("/admin/%s/create", urlExtension), controller.Create)
		mux.Post(fmt.Sprintf("/admin/%s/create", urlExtension), controller.Create)
		mux.Get(fmt.Sprintf("/admin/%s/create/success", urlExtension), controller.CreateSuccess)
		// Delete
		mux.Get(fmt.Sprintf("/admin/%s/delete/{id}", urlExtension), controller.Delete)
		mux.Post(fmt.Sprintf("/admin/%s/delete/{id}", urlExtension), controller.Delete)
		mux.Get(fmt.Sprintf("/admin/%s/delete/success", urlExtension), controller.DeleteSuccess)
		// Bulk delete (from table)
		mux.Delete(fmt.Sprintf("/admin/%s/bulk-delete", urlExtension), controller.BulkDelete)

		// Edit/Update (GET data in form / POST form)
		mux.Get(fmt.Sprintf("/admin/%s/{id}", urlExtension), controller.Edit)
		mux.Post(fmt.Sprintf("/admin/%s/{id}", urlExtension), controller.Edit)
		mux.Get(fmt.Sprintf("/admin/%s/edit/success", urlExtension), controller.EditSuccess)
	})
	return router
}

func (a api) AddAdminPolicySet(router *chi.Mux, protected bool, urlExtension string, controller adminpanel.AdminAuthPolicyController) *chi.Mux {
	// Reassign for consistency
	r := router
	r.Group(func(mux chi.Router) {
		// Set to use JWT authentication if protected
		if protected {
			mux.Use(auth.AuthenticateJWT)
		}
		// Read (all users)
		mux.Get(fmt.Sprintf("/admin/%s", urlExtension), controller.FindAll)
		mux.Get(fmt.Sprintf("/admin/%s/roles", urlExtension), controller.FindAllRoles)
		mux.Get(fmt.Sprintf("/admin/%s/inheritance", urlExtension), controller.FindAllRoleInheritance)
		// Create Policy (GET form / POST form)
		mux.Get(fmt.Sprintf("/admin/%s/create", urlExtension), controller.Create)
		mux.Post(fmt.Sprintf("/admin/%s/create", urlExtension), controller.Create)
		mux.Get(fmt.Sprintf("/admin/%s/create/success", urlExtension), controller.CreateSuccess)
		// Create Role
		mux.Get(fmt.Sprintf("/admin/%s/create-role", urlExtension), controller.CreateRole)
		mux.Post(fmt.Sprintf("/admin/%s/create-role", urlExtension), controller.CreateRole)
		mux.Get(fmt.Sprintf("/admin/%s/create-role/success", urlExtension), controller.CreateRoleSuccess)
		// Create inheritance
		mux.Get(fmt.Sprintf("/admin/%s/create-inheritance", urlExtension), controller.CreateInheritance)
		mux.Post(fmt.Sprintf("/admin/%s/create-inheritance", urlExtension), controller.CreateInheritance)
		mux.Get(fmt.Sprintf("/admin/%s/create-inheritance/success", urlExtension), controller.CreateInheritanceSuccess)
		// Delete Inheritance
		mux.Get(fmt.Sprintf("/admin/%s/delete-inheritance/{inherit-slug}", urlExtension), controller.DeleteInheritance)
		mux.Post(fmt.Sprintf("/admin/%s/delete-inheritance/{inherit-slug}", urlExtension), controller.DeleteInheritance)
		mux.Get(fmt.Sprintf("/admin/%s/delete-inheritance/success", urlExtension), controller.DeleteInheritanceSuccess)

		// Edit/Update (GET data in form / POST form)
		mux.Get(fmt.Sprintf("/admin/%s/{id}", urlExtension), controller.Edit)
		mux.Post(fmt.Sprintf("/admin/%s/{id}", urlExtension), controller.Edit)
		mux.Delete(fmt.Sprintf("/admin/%s/{id}", urlExtension), controller.Edit)
	})
	return router
}

// Function to add new routes to an existing Chi mux router
func (a api) AddBasicAdminRoutes(router *chi.Mux, controller adminpanel.AdminBaseController) *chi.Mux {
	// Public routes
	router.Group(func(mux chi.Router) {
		// @tag.name Public Routes
		// @tag.description Unprotected routes
		mux.Get("/admin", controller.Login)
		mux.Get("/admin/login", controller.Login)
		mux.Post("/admin/login", controller.Login)

		// admin logout
		mux.Get("/admin/logout", controller.Logout)

		// Private routes
		mux.Group(func(mux chi.Router) {
			mux.Use(auth.AuthenticateJWT)

			// @tag.name Private routes
			// @tag.description Protected routes
			// admin home
			mux.Get("/admin/home", controller.Home)
			// Change password
			mux.Get("/admin/change-password", controller.ChangePassword)

		})

	})

	return router
}
