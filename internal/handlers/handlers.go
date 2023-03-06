package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/services"
)

// Init state variable
var app *config.AppConfig

// Function called in main.go to connect app state to current file
func SetStateInHandlers(a *config.AppConfig) {
	app = a
}

// Login
// Handler to login with existing user
// @Summary      Login
// @Description  Log in to user account
// @Tags         Login
// @Accept       json
// @Produce      json
// @Param        user body models.Login true "Login JSON"
// @Success      200 {object} models.LoginResponse
// @Failure      401 {string} string "Invalid Credentials"
// @Failure      405 {string} string "Method not supported"
// @Router       /user/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Deny any request that is not a post
	if r.Method != "POST" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Init models for decoding
	var login models.Login
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Check if user exists in db
	foundUser, err := services.FindUserByEmail(login.Email)
	if err != nil {
		fmt.Println("Invalid credentials detected")
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	// If user is found
	// Compare stored (hashed) password with input password
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(login.Password))
	if err != nil {
		http.Error(w, "Incorrect username/password", http.StatusForbidden)
		return
	}

	// If match found (no errors)
	if err == nil {
		fmt.Println("User logging in: ", foundUser.Email)
		// Set login status to true
		tokenString, err := auth.GenerateJWT(int(foundUser.ID), foundUser.Email, foundUser.Role)
		if err != nil {
			fmt.Println("Failed to create JWT")
		}
		// Build login response
		var loginResponse = models.LoginResponse{Token: tokenString}
		// Send to user in body
		helpers.WriteAsJSON(w, loginResponse)
		return
	}
}

// API/ME
//
// Update the user's profile (using id from bearer token)
// @Summary      Update my profile
// @Description  Updates the currently logged in user
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user body models.UpdateUser true "Update User Json"
// @Param        id   path      int  true  "User ID"
// @Success      200 {object} models.PartialUser
// @Failure      400 {string} string "Failed user update"
// @Failure      403 {string} string "Authentication Token not detected"
// @Router       /user/{id} [put]
// @Security BearerToken
func UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	// grab id parameter
	var user models.UpdateUser
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}
	fmt.Printf("JSON Received: %+v\n", user)

	// Extract the user's id from their authentication token
	userId, err := auth.ExtractIdFromToken(w, r)
	if err != nil {
		http.Error(w, "Authentication Token not detected", http.StatusForbidden)
	}

	// Update user
	updatedUser, createErr := services.UpdateUser(userId, &user)
	if createErr != nil {
		http.Error(w, fmt.Sprintf("Failed user update: %s", createErr), http.StatusBadRequest)
		return
	}
	// Write updated user to output
	err = helpers.WriteAsJSON(w, updatedUser)
	if err != nil {
		fmt.Println("Error writing to JSON", err)
		return
	}
}

// Detail to display a user's profile details
// @Summary      Get my user profile details
// @Description  Return my user details
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200 {object} models.CreatedUser
// @Failure      400 {string} string "Can't find user details"
// @Router       /me [get]
// @Security BearerToken
func GetMyUserDetails(w http.ResponseWriter, r *http.Request) {
	// Grab ID from cookie
	// Validate the token
	tokenData, err := auth.ValidateAndParseToken(w, r)
	// If error detected
	if err != nil {
		http.Error(w, "Error parsing authentication token", http.StatusForbidden)
		return
	}

	// Convert to int
	idParameter, err := strconv.Atoi(tokenData.UserID)
	// If error detected
	if err != nil {
		fmt.Println("error parsing token to string: ", err)
		http.Error(w, "Error parsing authentication token", http.StatusForbidden)
		return
	}

	// Find user by id from cookie
	foundUser, err := services.FindUserById(idParameter)
	if err != nil {
		http.Error(w, "Can't find user details", http.StatusBadRequest)
		fmt.Println("error in finding user: ", err)
		return
	}

	// Write found user data to Response
	err = helpers.WriteAsJSON(w, foundUser)
	if err != nil {
		http.Error(w, "Can't find user details", http.StatusBadRequest)
		fmt.Println("error in finding user: ", err)
		return
	}
}

// Sample handler for JSON data: Jobs
func GetJobs(w http.ResponseWriter, r *http.Request) {
	var jobs []models.Job

	jobs = append(jobs, models.Job{ID: 1, Name: "Accounting"})
	jobs = append(jobs, models.Job{ID: 2, Name: "Programming"})

	// Set header
	w.Header().Set("Content-Type", "application/json")

	// Build new JSON encoder to write to, then write jobs data
	json.NewEncoder(w).Encode(jobs)
}

// Login URL check
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome!"))
}
