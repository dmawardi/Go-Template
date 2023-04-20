package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/service"
	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
)

type UserController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	// API/ME
	GetMyUserDetails(w http.ResponseWriter, r *http.Request)
	UpdateMyProfile(w http.ResponseWriter, r *http.Request)
	// Login
	Login(w http.ResponseWriter, r *http.Request)
}

type userController struct {
	service service.UserService
}

func NewUserController(service service.UserService) UserController {
	return &userController{service}
}

// API/USERS
// Find a list of users
// @Summary      Find a list of users
// @Description  Accepts limit, offset, and order params and returns list of users
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        limit   path      int  true  "limit"
// @Param        offset   path      int  true  "offset"
// @Param        order   path      int  true  "order by"
// @Success      200 {object} []models.CreatedUser
// @Failure      400 {string} string "Can't find users"
// @Failure      400 {string} string "Must include limit parameter with a max value of 50"
// @Router       /users/{id} [get]
// @Security BearerToken
func (c userController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Grab URL query parameters
	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")
	orderBy := r.URL.Query().Get("order")

	// Convert to int
	limit, _ := strconv.Atoi(limitParam)
	offset, _ := strconv.Atoi(offsetParam)

	// Check that limit is present as requirement
	if (limit == 0) || (limit >= 50) {
		http.Error(w, "Must include limit parameter with a max value of 50", http.StatusBadRequest)
		return
	}

	// Query database for all users using query params
	foundUsers, err := c.service.FindAll(limit, offset, orderBy)
	if err != nil {
		http.Error(w, "Can't find users", http.StatusBadRequest)
		return
	}
	err = helpers.WriteAsJSON(w, foundUsers)
	if err != nil {
		http.Error(w, "Can't find users", http.StatusBadRequest)
		fmt.Println("error writing users to response: ", err)
		return
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
// @Router       /users/{id} [get]
// @Security BearerToken
func (c userController) Find(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	fmt.Println("id parameter from request: ", stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	foundUser, err := c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find user with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
	err = helpers.WriteAsJSON(w, foundUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find user with ID: %v\n", idParameter), http.StatusBadRequest)
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
// @Router       /users [post]
func (c userController) Create(w http.ResponseWriter, r *http.Request) {
	// Init
	var user models.CreateUser
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&user)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Create user
	_, createErr := c.service.Create(&user)
	if createErr != nil {
		http.Error(w, "User creation failed.", http.StatusBadRequest)
		return
	}

	// Set status to created
	w.WriteHeader(http.StatusCreated)
	// Send user success message in body
	w.Write([]byte("User creation successful!"))
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
// @Router       /users/{id} [put]
// @Security BearerToken
func (c userController) Update(w http.ResponseWriter, r *http.Request) {
	// grab id parameter
	var user models.UpdateUser
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&user)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)

	// Update user
	updatedUser, createErr := c.service.Update(idParameter, &user)
	if createErr != nil {
		http.Error(w, fmt.Sprintf("Failed user update: %s", createErr), http.StatusBadRequest)
		return
	}
	// Write user to output
	err = helpers.WriteAsJSON(w, updatedUser)
	if err != nil {
		fmt.Printf("Error encountered when writing to JSON. Err: %s", err)
	}
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
// @Router       /users/{id} [delete]
// @Security BearerToken
func (c userController) Delete(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)

	// Attampt to delete user using id
	err := c.service.Delete(idParameter)

	// If error detected
	if err != nil {
		http.Error(w, "Failed user deletion", http.StatusBadRequest)
		return
	}
	// Else write success
	w.Write([]byte("Deletion successful!"))
	return
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
// @Failure      400 {string} string "Bad request"
// @Router       /user/{id} [put]
// @Security BearerToken
func (c userController) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	// grab id parameter
	var user models.UpdateUser
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("Decoding error: ", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&user)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Extract the user's id from their authentication token
	userId, err := auth.ExtractIdFromToken(w, r)
	if err != nil {
		http.Error(w, "Authentication Token not detected", http.StatusForbidden)
	}

	// Update user
	updatedUser, createErr := c.service.Update(*userId, &user)
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
func (c userController) GetMyUserDetails(w http.ResponseWriter, r *http.Request) {
	// Grab ID from cookie
	// Validate the token
	tokenData, err := auth.ValidateAndParseToken(w, r)
	// If error detected
	if err != nil {
		http.Error(w, "Error parsing authentication token:1", http.StatusForbidden)
		return
	}

	// Convert to int
	idParameter, err := strconv.Atoi(tokenData.UserID)
	// If error detected
	if err != nil {
		http.Error(w, "Error parsing authentication token:2", http.StatusForbidden)
		return
	}

	// Find user by id from cookie
	foundUser, err := c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, "Can't find user details", http.StatusBadRequest)
		return
	}

	// Write found user data to Response
	err = helpers.WriteAsJSON(w, foundUser)
	if err != nil {
		http.Error(w, "Can't find user details", http.StatusBadRequest)
		return
	}
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
func (c userController) Login(w http.ResponseWriter, r *http.Request) {
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

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&login)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through
	// Check if user exists in db
	foundUser, err := c.service.FindByEmail(login.Email)
	if err != nil {
		fmt.Println("Invalid credentials detected")
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	// If user is found
	// Compare stored (hashed) password with input password
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(login.Password))
	if err != nil {
		http.Error(w, "Incorrect username/password", http.StatusUnauthorized)
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
