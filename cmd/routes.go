package main

import (
	"net/http"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/controller"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

type Api interface {
	routes() http.Handler
}

type api struct {
	user controller.UserController
}

func NewApi(user controller.UserController) Api {
	return &api{user}
}

func (a api) routes() http.Handler {
	// Create new router
	mux := chi.NewRouter()
	// Use built in Chi middleware
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)

	// Public routes
	mux.Group(func(mux chi.Router) {
		// @tag.name Public Routes
		// @tag.description Unprotected routes
		mux.Get("/", controller.GetJobs)
		// Login
		mux.Post("/api/users/login", a.user.Login)

		// Create new user
		mux.Post("/api/users", a.user.Create)

		// Private routes
		mux.Group(func(mux chi.Router) {
			mux.Use(auth.AuthenticateJWT)

			// @tag.name Private routes
			// @tag.description Protected routes
			// users
			mux.Get("/api/users", a.user.FindAll)
			mux.Get("/api/users/{id}", a.user.Find)
			mux.Put("/api/users/{id}", a.user.Update)
			mux.Delete("/api/users/{id}", a.user.Delete)

			// My profile
			mux.Get("/api/me", a.user.GetMyUserDetails)
			mux.Post("/api/me", controller.HealthCheck)
			mux.Put("/api/me", a.user.UpdateMyProfile)

		})

	})

	// Serve API Swagger docs
	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/docs/swagger.json"), //The url pointing to API definition
	))

	// Build fileserver using static directory
	fileServer := http.FileServer(http.Dir("./static"))
	// Handle all calls to /static/* by stripping prefix and sending to file server
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return mux
}
