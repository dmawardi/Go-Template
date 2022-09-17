package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/services"
)

// Middleware to check whether user is authenticated
func AuthenticateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Grab request header
		header := r.Header
		// Extract token string from Authorization header by removing prefix "Bearer "
		_, tokenString, _ := strings.Cut(header.Get("Authorization"), " ")

		if tokenString == "" {
			http.Error(w, "Authentication Token not detected", http.StatusForbidden)
			return
		}
		// Validate the token
		tokenData, err := auth.ValidateAndParseToken(tokenString)
		// If error detected
		if err != nil {
			http.Error(w, "Error parsing authentication token", http.StatusForbidden)
			return
		}

		// Extract current URL being accessed
		fmt.Println("URL being accessed", r.URL.Path)
		object := r.URL.Path
		// Grab Http Method
		httpMethod := r.Method
		// Determine associated action based on HTTP method
		action := auth.ActionFromMethod(httpMethod)
		// Enforce RBAC policy and determine if user is authorized to perform action
		allowed := Authorize(tokenData.Email, object, action)

		if !allowed {
			http.Error(w, "Not authorized to perform that action", http.StatusForbidden)
			return
		}

		// Else, allow through
		next.ServeHTTP(w, r)
	})
}

// Middleware to check whether user is authorized
func Authorize(subjectEmail, object, action string) bool {
	fmt.Printf("\n%s is accessing %s to %s", subjectEmail, object, action)

	// Extract user ID from JWT and check if user exists in database.
	foundUser, err := services.FindUserByEmail(app.Ctx, app.DbClient, subjectEmail)
	if err != nil {
		fmt.Println("No user has been found in db with that id")
	}

	// Load Authorization policy from Database
	err = app.RBEnforcer.LoadPolicy()
	if err != nil {
		log.Fatal("Failed to load RBAC Enforcer policy in Authorization middleware")
		return false
	}

	fmt.Printf("\nEnforcing: %s accessing %s to %s\n", foundUser.Role, object, action)

	// Casbin enforces policy for user's role
	ok, err := app.RBEnforcer.Enforce(foundUser.Role, object, action)
	if err != nil {
		log.Fatal("Failed to enforce RBAC policy in Authorization middleware")
		return false
	}

	// Return result of enforcement
	return ok

}

// func Middleware(a auth.Authorizer, mux chi.Router) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			username, _, ok := r.BasicAuth()
// 			// This is where the password would normally be verified

// 			// asset := mux.Vars(r)["asset"]
// 			action := auth.ActionFromMethod(r.Method)
// 			if !ok || !a.HasPermission(username, action, asset) {
// 				log.Printf("User '%s' is not allowed to '%s' resource '%s'", username, action, asset)
// 				w.WriteHeader(http.StatusForbidden)
// 				return
// 			}

// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }
