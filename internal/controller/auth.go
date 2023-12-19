package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/service"
)

type AuthPolicyController interface {
	// Roles
	// Find a list of available roles
	FindAllRoles(w http.ResponseWriter, r *http.Request)
	// Assigns a role to a user
	AssignUserRole(w http.ResponseWriter, r *http.Request)
	// Policies
	FindAll(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type authPolicyController struct {
	service service.AuthPolicyService
}

func NewAuthPolicyController(service service.AuthPolicyService) AuthPolicyController {
	return &authPolicyController{service}
}

// API/POLICY
// Find a list of policies
// @Summary      Find a list of policies
// @Description  Returns list of policies
// @Tags         Policy
// @Accept       json
// @Produce      json
// @Success      200 {object} [][]string
// @Failure      400 {string} string "Can't find policies"
// @Router       /policy [get]
// @Security BearerToken
func (c authPolicyController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Find all
	policies, err := c.service.FindAll()
	if err != nil {
		http.Error(w, "Can't find policies", http.StatusBadRequest)
		return
	}
	// Return
	helpers.WriteAsJSON(w, policies)
}

// Find a list of roles
// @Summary      Find a list of roles
// @Description  Returns list of roles
// @Tags         Policy
// @Accept       json
// @Produce      json
// @Success      200 {object} []string
// @Failure      400 {string} string "Can't find roles"
// @Router       /policy/roles [get]
// @Security BearerToken
func (c authPolicyController) FindAllRoles(w http.ResponseWriter, r *http.Request) {
	// Find all roles
	roles, err := c.service.FindAllRoles()
	if err != nil {
		http.Error(w, "Can't find roles", http.StatusBadRequest)
		return
	}
	// Return posts
	helpers.WriteAsJSON(w, roles)
}

func (c authPolicyController) AssignUserRole(w http.ResponseWriter, r *http.Request) {
	// Grab request body as models.CasbinRule
	var pol models.CasbinRoleAssignment
	err := json.NewDecoder(r.Body).Decode(&pol)
	if err != nil {
		http.Error(w, "Invalid policy", http.StatusBadRequest)
		return
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&pol)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	success, err := c.service.AssignUserRole(pol.UserId, pol.Role)
	if err != nil {
		http.Error(w, "Can't assign user", http.StatusBadRequest)
		return
	}
	if !*success {
		http.Error(w, "Can't assign user", http.StatusBadRequest)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User assigned role successfully!"))
}

// Delete a policy
// @Summary      Delete a policy
// @Description  Delete a specific policy using the current policy
// @Tags         Policy
// @Accept       json
// @Produce      json
// @Success      200 {object} []string
// @Failure      400 {string} string "Can't delete policy"
// @Router       /policy [delete]
// @Security BearerToken
func (c authPolicyController) Delete(w http.ResponseWriter, r *http.Request) {
	// Grab request body as models.CasbinRule
	var pol models.CasbinRule
	err := json.NewDecoder(r.Body).Decode(&pol)
	if err != nil {
		http.Error(w, "Invalid policy", http.StatusBadRequest)
		return
	}

	err = c.service.Delete(pol)
	if err != nil {
		fmt.Printf("Error deleting policy: %v\n", err)
		http.Error(w, "Can't delete policy", http.StatusBadRequest)
		return
	}
	// Return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Policy deletion successful!"))
}

// Create a policy
// @Summary      Create a policy
// @Description  Create a specific policy using the current policy
// @Tags         Policy
// @Accept       json
// @Produce      json
// @Success      200 {object} string "Policy creation successful!"
// @Failure      400 {string} string "Can't create policy"
// @Router       /policy [post]
// @Security BearerToken
func (c authPolicyController) Create(w http.ResponseWriter, r *http.Request) {
	// Grab request body as models.CasbinRule
	var pol models.CasbinRule
	err := json.NewDecoder(r.Body).Decode(&pol)
	if err != nil {
		http.Error(w, "Invalid policy", http.StatusBadRequest)
		return
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&pol)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Ensure that the DTO is a policy and not a group policy
	if pol.PType == "g" {
		http.Error(w, "Can't assign a role. Try different route", http.StatusBadRequest)
		return
	}

	// Create the policy
	err = c.service.Create(pol)
	if err != nil {
		http.Error(w, "Can't create policy", http.StatusBadRequest)
		return
	}

	// Return success
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Policy creation successful!"))
}

// Update a policy
// @Summary      Update a policy
// @Description  Update a specific policy using the current policy
// @Tags         Policy
// @Accept       json
// @Produce      json
// @Success      200 {object} string "Policy update successful!"
// @Failure      400 {string} string "Can't update policy"
// @Router       /policy [put]
// @Security BearerToken
func (c authPolicyController) Update(w http.ResponseWriter, r *http.Request) {
	// Grab request body as models.CasbinRule
	var pol models.UpdateCasbinRule
	err := json.NewDecoder(r.Body).Decode(&pol)
	if err != nil {
		http.Error(w, "Invalid policy", http.StatusBadRequest)
		return
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&pol)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	err = c.service.Update(pol.OldPolicy, pol.NewPolicy)
	if err != nil {
		fmt.Printf("Error updating policy: %v\n", err)
		http.Error(w, "Can't update policy", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Policy update successful!"))
}
