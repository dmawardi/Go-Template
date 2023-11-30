package auth

import (
	"fmt"
	"net/http"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers"
)

// Middleware to check whether user is authenticated
func AuthenticateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Validate the token
		tokenData, err := ValidateAndParseToken(w, r)
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
		allowed := Authorize(tokenData.UserID, object, action)

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
func Authorize(userId, object, action string) bool {
	// Load Authorization policy from Database
	err := app.RBEnforcer.LoadPolicy()
	if err != nil {
		fmt.Printf("Failed to load RBAC Enforcer policy in Authorization middleware")
		return false
	}

	// Enforce policy for user's role using their ID
	ok, err := app.RBEnforcer.Enforce(userId, object, action)
	if err != nil {
		fmt.Print("Failed to enforce RBAC policy in Authorization middleware: ", err, "\nUser ID: ", userId, "\nObject: ", object, "\nAction: ", action, "\n")
		return false
	}
	fmt.Printf("User with ID %s is accessing %s to %s. Allowed? %v\n", userId, object, action, ok)

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
