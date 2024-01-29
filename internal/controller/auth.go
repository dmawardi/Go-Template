package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/service"
	"github.com/go-chi/chi"
)

type AuthPolicyController interface {
	// Roles
	FindAllRoles(w http.ResponseWriter, r *http.Request)
	AssignUserRole(w http.ResponseWriter, r *http.Request)
	// Policies
	FindAll(w http.ResponseWriter, r *http.Request)
	FindByResource(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	// Inheritance
	FindAllRoleInheritance(w http.ResponseWriter, r *http.Request)
	CreateInheritance(w http.ResponseWriter, r *http.Request)
	DeleteInheritance(w http.ResponseWriter, r *http.Request)
}

type authPolicyController struct {
	service service.AuthPolicyService
}

func NewAuthPolicyController(service service.AuthPolicyService) AuthPolicyController {
	return &authPolicyController{service}
}

// API/POLICY

// POLICIES
//

// Find a list of policies
// @Summary      Find a list of policies
// @Description  Returns list of policies
// @Tags         Policy
// @Accept       json
// @Produce      json
// @Param        order   query      int  false  "order by"
// @Success      200 {object} [][]string
// @Failure      400 {string} string "Can't find policies"
// @Router       /policy [get]
// @Security BearerToken
func (c authPolicyController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Grab search query
	searchQuery := r.URL.Query().Get("searchQuery")
	// Find all
	policies, err := c.service.FindAll(searchQuery)
	if err != nil {
		http.Error(w, "Can't find policies", http.StatusBadRequest)
		return
	}
	// Return
	helpers.WriteAsJSON(w, policies)
}

// @Summary      Find policies asssociated with a resource
// @Description  Returns list of policies associated with a resource
// @Tags         Policy
// @Accept       json
// @Produce      json
// @Param        policy-slug   path      string  true  "Policy Slug"
// @Success      200 {object} [][]string
// @Failure      400 {string} string "Can't find policies"
// @Router       /policy [get]
// @Security BearerToken
func (c authPolicyController) FindByResource(w http.ResponseWriter, r *http.Request) {
	// Grab search query
	policyResource := chi.URLParam(r, "policySlug")

	// Unslugify
	policyResource = helpers.UnslugifyResourceName(policyResource)

	// Find all
	policies, err := c.service.FindByResource(policyResource)
	if err != nil || len(policies) == 0 {
		http.Error(w, "Can't find policies", http.StatusBadRequest)
		return
	}
	// Return
	helpers.WriteAsJSON(w, policies)
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
	var pol models.PolicyRule
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
	var pol models.PolicyRule
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

// ROLES
//

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

// Assigns a role to a user or creates a new role for a user
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

// INHERITANCE
//

func (c authPolicyController) FindAllRoleInheritance(w http.ResponseWriter, r *http.Request) {
	// Find all roles
	roles, err := c.service.FindAllRoleInheritance()
	if err != nil {
		http.Error(w, "Can't find roles", http.StatusBadRequest)
		return
	}
	// Return posts
	helpers.WriteAsJSON(w, roles)
}

func (c authPolicyController) CreateInheritance(w http.ResponseWriter, r *http.Request) {
	// Grab request body as models.CasbinRule
	var pol models.G2Record
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

	err = c.service.CreateInheritance(pol)
	if err != nil {
		http.Error(w, "Can't create inheritance", http.StatusBadRequest)
		return
	}

	// Return success
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Inheritance creation successful!"))
}

func (c authPolicyController) DeleteInheritance(w http.ResponseWriter, r *http.Request) {
	// Grab request body as models.CasbinRule
	var pol models.G2Record
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

	err = c.service.DeleteInheritance(pol)
	if err != nil {
		http.Error(w, "Can't delete inheritance", http.StatusBadRequest)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Inheritance deletion successful!"))
}
