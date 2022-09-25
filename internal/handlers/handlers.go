package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
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
	if r.Method != "POST" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Init
	var login models.Login
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Check if user exists in db
	foundUser, err := services.FindUserByEmail(app.Ctx, app.DbClient, login.Email)
	fmt.Println("founduser: ", foundUser)
	// If user found
	if err == nil {
		fmt.Println("User logging in: ", foundUser)

		// Compare stored (hashed) password with input password
		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(login.Password))

		// If match found (no errors)
		if err == nil {
			// Set login status to true
			tokenString, err := auth.GenerateJWT(foundUser.ID, foundUser.Email, foundUser.Role)
			// helpers.WriteAsJSON(w, )
			if err != nil {
				fmt.Println("Failed to create JWT")
			}
			// Build login response
			var loginResponse = models.LoginResponse{Token: tokenString}
			// Send to user in body
			helpers.WriteAsJSON(w, loginResponse)
			return
			// else if user password doesn't match
		} else {
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}
	}
}

// Login URL check
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome!"))
}

// Users

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
		http.Error(w, "Error parsing authentication token", http.StatusForbidden)
		fmt.Println("error parsing token to string: ", err)
		return
	}

	// Find user by id from cookie
	foundUser, err := services.FindUserById(app.Ctx, app.DbClient, idParameter)
	if err != nil {
		http.Error(w, "Can't find user details", http.StatusBadRequest)
		fmt.Println("error in finding user: ", err)
		return
	}
	err = helpers.WriteAsJSON(w, foundUser)
	if err != nil {
		http.Error(w, "Can't find user details", http.StatusBadRequest)
		fmt.Println("error in finding user: ", err)
		return
	}
}

// Create a new user
// @Summary      Create User
// @Description  Creates a new user
// @Tags         User
// @Accept       json
// @Produce      plain
// @Param        user body models.CreateUser true "NewUserJson"
// @Success      201 {string} string "User creation successful!"
// @Failure      400 {string} string "User creation failed."
// @Router       /user [post]
func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	// Init
	var user models.CreateUser
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}
	fmt.Printf("Create User Received: %+v\n", user)

	// Create user
	createdUser, createErr := services.CreateUser(app.Ctx, app.DbClient, &user)
	if createErr != nil {
		http.Error(w, "User creation failed.", http.StatusBadRequest)
		return
	}
	fmt.Println("created user: ", createdUser)
	// Set status to created
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User creation successful!"))
	if err != nil {
		fmt.Println(err)
	}
}

// Find a created user
// @Summary      Find User
// @Description  Find a user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200 {object} models.CreatedUser
// @Failure      400 {string} string "Can't find user"
// @Router       /user/{id} [get]
// @Security BearerToken
func FindUser(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)
	fmt.Println("id parameter from request: ", idParameter)

	foundUser, err := services.FindUserById(app.Ctx, app.DbClient, idParameter)
	if err != nil {
		http.Error(w, "Can't find user", http.StatusBadRequest)
		fmt.Println("error in finding user: ", err)
		return
	}
	err = helpers.WriteAsJSON(w, foundUser)
	if err != nil {
		http.Error(w, "Can't find user", http.StatusBadRequest)
		fmt.Println("error in finding user: ", err)
		return
	}
}

// Update a user (using URL parameter id)
// @Summary      Update User
// @Description  Updates an existing user
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user body models.UpdateUser true "Update User Json"
// @Param        id   path      int  true  "User ID"
// @Success      200 {object} models.UpdatedUser
// @Failure      400 {string} string "Failed user update"
// @Failure      403 {string} string "Authentication Token not detected"
// @Router       /user/{id} [put]
// @Security BearerToken
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// grab id parameter
	var user models.UpdateUser
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)
	// Set within request body data object
	user.Id = idParameter

	fmt.Printf("JSON Received: %+v\n", user)

	// Update user
	updatedUser, createErr := services.UpdateUser(app.Ctx, app.DbClient, &user)
	if createErr != nil {
		http.Error(w, fmt.Sprintf("Failed user update: %s", createErr), http.StatusBadRequest)
		return
	}
	// Write user to output
	err = helpers.WriteAsJSON(w, updatedUser)
	fmt.Println(err)
}

// Delete user (using URL parameter id)
// @Summary      Delete User
// @Description  Deletes an existing user
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200 {string} string "Deletion successful!"
// @Failure      400 {string} string "Failed user deletion"
// @Router       /user/{id} [delete]
// @Security BearerToken
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)

	// Attampt to delete user using id
	err := services.DeleteUser(app.Ctx, app.DbClient, idParameter)

	// If error detected
	if err != nil {
		http.Error(w, "Failed user deletion", http.StatusBadRequest)
		return
	}
	// Else write success
	w.Write([]byte("Deletion successful!"))
	return
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
