package routes

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/controller"
	"github.com/go-chi/chi"
)

// Adds a basic fully authorized CRUD route set to a Chi mux router
func AddBasicCrudApiRoutes(router *chi.Mux, urlExtension string, controller controller.BasicController) *chi.Mux {
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
