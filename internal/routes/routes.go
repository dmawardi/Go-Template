package routes

import (
	"net/http"

	adminpanel "github.com/dmawardi/Go-Template/internal/admin-panel"
	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/controller"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

type Api interface {
	Routes() http.Handler
	AddAdminRoutes(router *chi.Mux) *chi.Mux
}

type api struct {
	Admin adminpanel.AdminController
	User  controller.UserController
	Post  controller.PostController
}

func NewApi(admin adminpanel.AdminController, user controller.UserController, post controller.PostController) Api {
	return &api{Admin: admin, User: user, Post: post}
}

func (a api) Routes() http.Handler {
	// Create new router
	mux := chi.NewRouter()
	// Use built in Chi middleware
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)
	mux.Use(corsMiddleware)

	// Public routes
	mux.Group(func(mux chi.Router) {
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

			// Posts
			mux.Get("/api/posts", a.Post.FindAll)
			mux.Get("/api/posts/{id}", a.Post.Find)
			mux.Put("/api/posts/{id}", a.Post.Update)
			mux.Post("/api/posts", a.Post.Create)
			mux.Delete("/api/posts/{id}", a.Post.Delete)
		})

	})

	// Serve API Swagger docs
	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/static/docs/swagger.json"), //The url pointing to API definition
	))

	// Build fileserver using static directory
	fileServer := http.FileServer(http.Dir("./static"))
	// Handle all calls to /static/* by stripping prefix and sending to file server
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// // Add admin routes
	mux = a.AddAdminRoutes(mux)

	return mux
}

// Function to add new routes to an existing Chi mux router
func (a api) AddAdminRoutes(router *chi.Mux) *chi.Mux {
	// Admin routes
	// router.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("This is the admin login page"))
	// })
	// Public routes
	router.Group(func(mux chi.Router) {
		// @tag.name Public Routes
		// @tag.description Unprotected routes
		mux.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("This is the admin login page sailor"))
		})
		// admin users
		// Read (all users)
		mux.Get("/admin/users", a.Admin.User.FindAll)
		// Create (GET form / POST form)
		mux.Get("/admin/users/create", a.Admin.User.Create)
		mux.Post("/admin/users/create", a.Admin.User.Create)
		mux.Get("/admin/users/create/success", a.Admin.User.CreateSuccess)
		// Delete
		mux.Get("/admin/users/delete/{id}", a.Admin.User.Delete)
		mux.Post("/admin/users/delete/{id}", a.Admin.User.Delete)
		// Edit/Update (GET data in form / POST form)
		mux.Get("/admin/users/{id}", a.Admin.User.Edit)
		mux.Post("/admin/users/{id}", a.Admin.User.Edit)
		mux.Get("/admin/users/edit/success", a.Admin.User.EditSuccess)

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
