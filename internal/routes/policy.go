package routes

import (
	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/go-chi/chi"
)

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
			// AUTH
			//
			// Policies
			mux.Get("/api/auth", a.Policy.FindAll)
			mux.Get("/api/auth/{policy-slug}", a.Policy.FindByResource)
			mux.Post("/api/auth", a.Policy.Create)
			mux.Put("/api/auth", a.Policy.Update)
			mux.Delete("/api/auth", a.Policy.Delete)
			// Roles
			mux.Get("/api/auth/roles", a.Policy.FindAllRoles)
			mux.Put("/api/auth/roles", a.Policy.AssignUserRole)
			mux.Post("/api/auth/roles", a.Policy.CreateRole)
			// Inheritance
			mux.Get("/api/auth/inheritance", a.Policy.FindAllRoleInheritance)
			mux.Post("/api/auth/inheritance", a.Policy.CreateInheritance)
			mux.Delete("/api/auth/inheritance", a.Policy.DeleteInheritance)
		})

	})
	return router
}
