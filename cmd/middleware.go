package main

import (
	"net/http"
	"strings"

	"github.com/dmawardi/Go-Template/internal/auth"
)

// Middleware to check whether user is authorized
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
		_, err := auth.ValidateAndParseToken(tokenString)
		// If error detected
		if err != nil {
			http.Error(w, "Error parsing authentication token", http.StatusForbidden)
			return
		}

		// Else, allow through
		next.ServeHTTP(w, r)
	})
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
