package repository

import (
	"errors"
	"fmt"

	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
	"gorm.io/gorm"
)

// CasbinPolicyRepository represents a repository for Casbin policies.
type AuthPolicyRepository interface {
	// Roles
	FindAllRoles() ([]string, error)
	FindRoleByUserId(userId string) (string, error)
	AssignUserRole(userId, roleToApply string) (*bool, error)
	DeleteAllUserRoles(userID string) (*bool, error)
	// Role Inheritance
	FindAllRoleInheritance() ([][]string, error)
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

// Roles
// Returns all role inheritance records
func (r *authPolicyRepository) FindAllRoleInheritance() ([][]string, error) {
	// return all policies found in the database
	g2Records := r.auth.Enforcer.GetNamedGroupingPolicy("g2")

	return g2Records, nil
}

// FindAll returns all Casbin policies.
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
	// Check if user exists
	user := db.User{}
	result := r.db.Where("id = ?", userId).First(&user)
	if result.Error != nil {
		fmt.Printf("Error finding user with id: %v\n", userId)
		return nil, result.Error
	}

	// If user exists, proceed
	// First, remove the existing roles for the user (if found)
	_, err := r.auth.Enforcer.DeleteRolesForUser(userId)
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
func (r *authPolicyRepository) DeleteAllUserRoles(userID string) (*bool, error) {
	// Set default result
	result := false
	// Remove all roles for user
	_, err := r.auth.Enforcer.DeleteRolesForUser(userID)
	if err != nil {
		fmt.Printf("Error removing roles for user: %v\n", err)
		result = false
		return &result, err
	}
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
	fmt.Printf("Old policy: %v\n", oldPolicy)
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
