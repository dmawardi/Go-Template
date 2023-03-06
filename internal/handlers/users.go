package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/services"
	"github.com/go-chi/chi"
)

// API/USERS
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
func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	// Init
	var user models.CreateUser
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Create user
	_, createErr := services.CreateUser(&user)
	if createErr != nil {
		http.Error(w, "User creation failed.", http.StatusBadRequest)
		return
	}

	// Set status to created
	w.WriteHeader(http.StatusCreated)
	// Send user success message in body
	w.Write([]byte("User creation successful!"))
}

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
func FindAllUsers(w http.ResponseWriter, r *http.Request) {
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
	foundUsers, err := services.FindAllUsers(limit, offset, orderBy)
	if err != nil {
		http.Error(w, "Can't find users", http.StatusBadRequest)
		fmt.Println("error in finding users: ", err)
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
func FindUser(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)
	fmt.Println("id parameter from request: ", idParameter)

	foundUser, err := services.FindUserById(idParameter)
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
// @Router       /users/{id} [put]
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

	fmt.Printf("JSON Received: %+v\n", user)

	// Update user
	updatedUser, createErr := services.UpdateUser(&idParameter, &user)
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
// @Router       /users/{id} [delete]
// @Security BearerToken
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)

	// Attampt to delete user using id
	err := services.DeleteUser(idParameter)

	// If error detected
	if err != nil {
		http.Error(w, "Failed user deletion", http.StatusBadRequest)
		return
	}
	// Else write success
	w.Write([]byte("Deletion successful!"))
	return
}
