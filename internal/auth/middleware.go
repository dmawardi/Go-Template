package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers"
)

// Middleware to check whether user is authenticated
func AuthenticateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Validate the token
		tokenData, err := ValidateAndParseToken(w, r)
		fmt.Println("tokendata received: ", tokenData)
		// If error detected
		if err != nil {
			http.Error(w, "Error parsing authentication token", http.StatusForbidden)
			return
		}

		// Extract current URL being accessed
		object := helpers.ExtractBasePath(r)

		// Grab Http Method
		httpMethod := r.Method
		// Determine associated action based on HTTP method
		action := ActionFromMethod(httpMethod)
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

	// Extract user ID from JWT and check if user exists in database.
	foundUser, err := FindByEmail(subjectEmail)
	if err != nil {
		fmt.Println("No user has been found in db with that id")
		return false
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
	fmt.Printf("%s is accessing %s to %s. Allowed? %v\n", subjectEmail, object, action, ok)

	// Return result of enforcement
	return ok
}

// Find user in database by email (for authentication)
func FindByEmail(email string) (*db.User, error) {
	// Create an empty ref object of type user
	user := db.User{}
	// Check if user exists in db
	result := app.DbClient.Where("email = ?", email).First(&user)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in finding user in authentication: ", result.Error)
		return nil, result.Error
	}
	// else
	return &user, nil
}
