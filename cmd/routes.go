package main

import (
	"net/http"

	"github.com/dmawardi/Go-Template/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes() http.Handler {
	// Create new router
	mux := chi.NewRouter()
	// Use built in Chi middleware
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)

	mux.Group(func(mux chi.Router) {
		// Public Routes
		mux.Get("/", handlers.GetJobs)
		// Login
		mux.Post("/api/login", handlers.LoginHandler)

		// Private routes
		mux.Group(func(mux chi.Router) {
			mux.Use(AuthenticateJWT)

			// User
			mux.Post("/api/user", handlers.CreateNewUser)
			mux.Patch("/api/user/{id}", handlers.UpdateUser)
			mux.Delete("/api/user/{id}", handlers.DeleteUser)

			mux.Get("/api/me", handlers.HealthCheck)

		})

	})

	// Build fileserver using static directory
	fileServer := http.FileServer(http.Dir("./static"))
	// Handle all calls to /static/* by stripping prefix and sending to file server
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return mux
}
