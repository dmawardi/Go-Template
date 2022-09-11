package auth

import (
	"net/http"

	"github.com/dmawardi/Go-Template/internal/config"
)

var app *config.AppConfig

// Function called in main.go to connect app state to current file
func SetStateInAuth(a *config.AppConfig) {
	app = a
}

// Uses session to update login status
func SetLoginStatus(w http.ResponseWriter, r *http.Request, status bool) error {
	// Get session or create new if doesn't already exist
	session, err := app.Session.Get(r, "session.id")
	if err != nil {
		return err
	}
	// Set authenticated value
	session.Values["authenticated"] = status
	// Saves session during the current request
	session.Save(r, w)
	// return nil / success
	return nil
}

func IsAuthenticated(r *http.Request) bool {
	// Grab session from app
	session, _ := app.Session.Get(r, "session.id")
	// Grab authenticated value from session
	authenticated := session.Values["authenticated"]

	// If value is present and not false
	if authenticated != nil && authenticated != false {
		return true
	}
	// else
	return false
}
