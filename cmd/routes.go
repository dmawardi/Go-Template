package main

import (
	"net/http"

	"github.com/dmawardi/Go-Template/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

func routes() http.Handler {
	// Create new router
	mux := chi.NewRouter()
	// Use built in Chi middleware
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)

	// Public routes
	mux.Group(func(mux chi.Router) {
		// @tag.name Public Routes
		// @tag.description Unprotected routes
		mux.Get("/", handlers.GetJobs)
		// Login
		mux.Post("/api/users/login", handlers.LoginHandler)

		// Create new user
		mux.Post("/api/users", handlers.CreateNewUser)

		// Private routes
		mux.Group(func(mux chi.Router) {
			mux.Use(AuthenticateJWT)

			// @tag.name Private routes
			// @tag.description Protected routes
			// users
			mux.Get("/api/users", handlers.FindAllUsers)
			mux.Get("/api/users/{id}", handlers.FindUser)
			mux.Put("/api/users/{id}", handlers.UpdateUser)
			mux.Delete("/api/users/{id}", handlers.DeleteUser)

			// My profile
			mux.Get("/api/me", handlers.GetMyUserDetails)
			mux.Post("/api/me", handlers.HealthCheck)
			mux.Put("/api/me", handlers.UpdateMyProfile)

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
