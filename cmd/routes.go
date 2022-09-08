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

	// Routes
	mux.Get("/", handlers.GetJobs)

	// Build fileserver using static directory
	fileServer := http.FileServer(http.Dir("./static"))
	// Handle all calls to /static/* by stripping prefix and sending to file server
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return mux
}
