package routes

import (
	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/controller"
	"github.com/go-chi/chi"
)

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

			// Email verification
			mux.Post("/api/users/send-verification-email", a.User.ResendVerificationEmail)

		})

	})
	return router
}
