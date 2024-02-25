package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"gorm.io/gorm"
)

// CasbinPolicyRepository represents a repository for Casbin policies.
type AuthPolicyRepository interface {
	// Roles
	FindAllRoles() ([]string, error)
	FindRoleByUserId(userId string) (string, error)
	CreateRole(userId, roleToApply string) (*bool, error)
	AssignUserRole(userId, roleToApply string) (*bool, error)
	DeleteRolesForUser(userID string) (*bool, error)

	// Role Inheritance
	FindAllRoleInheritance() ([]models.GRecord, error)
	// Should not contain role: prefix when passed. All handled in repository
	CreateInheritance(inherit models.GRecord) error
	// Should not contain role: prefix when passed. All handled in repository
	DeleteInheritance(inherit models.GRecord) error

	// Policies
	FindAll() ([][]string, error)
	Create(policy models.CasbinRule) error
	Update(oldPolicy, newPolicy models.CasbinRule) error
	Delete(policy models.CasbinRule) error
}

// GormCasbinPolicyRepository is a GORM implementation of CasbinPolicyRepository.
type authPolicyRepository struct {
	db   *gorm.DB
	auth config.AuthEnforcer
}

// NewGormCasbinPolicyRepository creates a new instance of GormCasbinPolicyRepository.
func NewAuthPolicyRepository(db *gorm.DB) AuthPolicyRepository {
	return &authPolicyRepository{
		db:   db,
		auth: app.Auth,
	}
}

// Role inheritance
// Returns all role inheritance records
func (r *authPolicyRepository) FindAllRoleInheritance() ([]models.GRecord, error) {
	// return all policies found in the database
	rolesAndAssignments := r.auth.Enforcer.GetNamedGroupingPolicy("g")
	fmt.Printf("Roles: %v\n", rolesAndAssignments)
	// Filter out the roles that are not user assigned
	// Filtered slice
	// var roles []string

	var roleInheritancePolicies []models.GRecord

	for _, policy := range rolesAndAssignments {
		// Assuming policy[0] contains the role/subject and policy[1] contains the inherited role
		// Adjust the indexing based on your actual policy structure
		if strings.HasPrefix(policy[0], "role:") && strings.HasPrefix(policy[1], "role:") {
			roleInheritancePolicies = append(roleInheritancePolicies, models.GRecord{Role: policy[0], InheritsFrom: policy[1]})
		}
	}

	return roleInheritancePolicies, nil
}
func (r *authPolicyRepository) CreateInheritance(inherit models.GRecord) error {
	// Apply naming convention to new role record
	addRolePrefix(&inherit)

	roles := r.auth.Enforcer.GetAllRoles()
	roleFound1 := helpers.ArrayContainsString(roles, inherit.Role)
	roleFound2 := helpers.ArrayContainsString(roles, inherit.InheritsFrom)
	if !roleFound1 || !roleFound2 {
		return errors.New("inheritance roles not found")
	}

	// Else, proceed to add the policy
	// Add policy to enforcer using add role for user, but it will be role to role
	newPolicy, err := r.auth.Enforcer.AddNamedGroupingPolicy("g", inherit.Role, inherit.InheritsFrom)
	if err != nil {
		return err
	}
	// If not new, return error
	if !newPolicy {
		return errors.New("policy already exists")
	}
	// else, return success
	return nil
}
func (r *authPolicyRepository) DeleteInheritance(inherit models.GRecord) error {
	// Apply naming convention to new role record
	addRolePrefix(&inherit)
	// Remove policy from enforcer
	removed, err := r.auth.Enforcer.RemoveNamedGroupingPolicy("g", inherit.Role, inherit.InheritsFrom)
	if err != nil {
		return err
	}

	// If not removed, return error
	if !removed {
		return errors.New("policy does not exist")
	}

	// else, return success
	return nil
}

// Roles
func (r *authPolicyRepository) FindAllRoles() ([]string, error) {
	// return all policies found in the databaseq
	roles := r.auth.Enforcer.GetAllRoles()
	return roles, nil
}
func (r *authPolicyRepository) FindRoleByUserId(userId string) (string, error) {
	// return all policies found in the databaseq
	roles, err := r.auth.Enforcer.GetRolesForUser(userId)
	if err != nil {
		return "", err
	}
	// If no roles found, return error
	if len(roles) == 0 {
		return "", errors.New("no roles found for user")
	}
	// Return first found role (should be only role)
	return roles[0], nil
}
func (r *authPolicyRepository) AssignUserRole(userId, roleToApply string) (*bool, error) {
	// Apply naming convention role to apply
	roleToApply = "role:" + roleToApply
	// Check if user exists
	user := db.User{}
	result := r.db.Where("id = ?", userId).First(&user)
	if result.Error != nil {
		fmt.Printf("Error finding user with id: %v\n", userId)
		return nil, result.Error
	}

	// If user exists, proceed to check if role exists
	roles, err := r.FindAllRoles()
	if err != nil {
		fmt.Printf("Error finding roles: %v\n", err)
		return nil, err
	}

	roleFound := helpers.ArrayContainsString(roles, roleToApply)
	if !roleFound {
		fmt.Printf("Role not found: %v\nCurrent roles: %v\n", roleToApply, roles)
		return nil, errors.New("role not found")
	}

	// First, remove the existing roles for the user (if found)
	_, err = r.auth.Enforcer.DeleteRolesForUser(userId)
	if err != nil {
		fmt.Printf("Error removing roles for user: %v\n", err)
		return nil, err
	}

	// Add the new role for the user.
	success, err := r.auth.Enforcer.AddRoleForUser(userId, roleToApply)
	if err != nil {
		fmt.Printf("Error assigning role to user: %v\n", err)
		return nil, err
	}

	return &success, nil
}
func (r *authPolicyRepository) CreateRole(userId, roleToApply string) (*bool, error) {
	// Check if user exists
	user := db.User{}
	result := r.db.Where("id = ?", userId).First(&user)
	if result.Error != nil {
		fmt.Printf("Error finding user with id: %v\n", userId)
		return nil, result.Error
	}

	// Apply naming convention to new role record
	roleToApply = "role:" + roleToApply

	// Create the new role with the user as the first member
	success, err := r.auth.Enforcer.AddRoleForUser(userId, roleToApply)
	if err != nil {
		fmt.Printf("Error assigning role to user: %v\n", err)
		return nil, err
	}

	return &success, nil
}
func (r *authPolicyRepository) DeleteRolesForUser(userID string) (*bool, error) {
	// Set default result
	result := false
	// Remove all roles for user
	_, err := r.auth.Enforcer.DeleteRolesForUser(userID)
	if err != nil {
		fmt.Printf("Error removing roles for user: %v\n", err)
		result = false
		return &result, err
	}
	// Determine as success
	result = true

	return &result, nil
}

// Policies
func (r *authPolicyRepository) FindAll() ([][]string, error) {
	// return all policies found in the database
	policies := r.auth.Enforcer.GetPolicy()
	return policies, nil
}
func (r *authPolicyRepository) Create(policy models.CasbinRule) error {
	// Add policy to enforcer
	newPolicy, err := r.auth.Enforcer.AddPolicy(policy.V0, policy.V1, policy.V2)
	if err != nil {
		return err
	}
	// If not new, return error
	if !newPolicy {
		return errors.New("policy already exists")
	}
	// else, return success
	return nil
}
func (r *authPolicyRepository) Delete(policy models.CasbinRule) error {
	var removed bool
	var err error

	// Remove policy from enforcer
	removed, err = r.auth.Enforcer.RemovePolicy(policy.V0, policy.V1, policy.V2)
	if err != nil {
		return err
	}

	// If not removed, return error
	if !removed {
		return errors.New("policy does not exist")
	}

	// else, return success
	return nil
}
func (r *authPolicyRepository) Update(oldPolicy, newPolicy models.CasbinRule) error {
	// Remove old policy from enforcer
	removed, err := r.auth.Enforcer.RemovePolicy(oldPolicy.V0, oldPolicy.V1, oldPolicy.V2)
	if err != nil {
		fmt.Printf("Error removing old policy: %v\n", err)
		return err
	}
	// If not removed, return error
	if !removed {
		fmt.Printf("Policy to update doesn't exist: %v\n", oldPolicy)
		return errors.New("policy to update does not exist")
	}
	// Add new policy to enforcer
	addedPolicy, err := r.auth.Enforcer.AddPolicy(newPolicy.V0, newPolicy.V1, newPolicy.V2)
	if err != nil {
		return err
	}
	// If not new, return error
	if !addedPolicy {
		return errors.New("policy already exists")
	}
	// else, return success
	return nil
}

// Helper functions
// Used to implement role naming convention
func addRolePrefix(inherit *models.GRecord) {
	// Apply naming convention to new role record
	inherit.Role = "role:" + inherit.Role
	inherit.InheritsFrom = "role:" + inherit.InheritsFrom
}
