package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/services"
)

// Middleware to check whether user is authenticated
func AuthenticateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Validate the token
		tokenData, err := auth.ValidateAndParseToken(w, r)
		// If error detected
		if err != nil {
			http.Error(w, "Error parsing authentication token", http.StatusForbidden)
			return
		}

		// Extract current URL being accessed
		object := helpers.ExtractBasePath(r)

		fmt.Println("current path: ", object)
		// Grab Http Method
		httpMethod := r.Method
		// Determine associated action based on HTTP method
		action := auth.ActionFromMethod(httpMethod)
		// Enforce RBAC policy and determine if user is authorized to perform action
		allowed := Authorize(tokenData.Email, object, action)

		// If not allowed
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
	fmt.Printf("\n%s is accessing %s to %s\n", subjectEmail, object, action)

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

	// Enforce policy for user's role
	ok, err := app.RBEnforcer.Enforce(foundUser.Role, object, action)
	if err != nil {
		log.Fatal("Failed to enforce RBAC policy in Authorization middleware")
		return false
	}

	// Return result of enforcement
	return ok
}
