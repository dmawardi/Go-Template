package repository

import (
	"errors"
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
	"gorm.io/gorm"
)

// CasbinPolicyRepository represents a repository for Casbin policies.
type AuthPolicyRepository interface {
	FindAll() ([][]string, error)
	FindAllRoles() ([]string, error)
	AssignUserRole(userId, roleToApply string) (*bool, error)
	// FindByUserID(userID string) ([]casbin.Policy, error)
	Create(policy models.CasbinRule) error
	Update(oldPolicy, newPolicy models.CasbinRule) error
	Delete(policy models.CasbinRule) error
}

// GormCasbinPolicyRepository is a GORM implementation of CasbinPolicyRepository.
type authPolicyRepository struct {
	db       *gorm.DB
	enforcer *casbin.Enforcer
}

// NewGormCasbinPolicyRepository creates a new instance of GormCasbinPolicyRepository.
func NewAuthPolicyRepository(db *gorm.DB) AuthPolicyRepository {
	return &authPolicyRepository{
		db:       db,
		enforcer: app.RBEnforcer,
	}
}

// FindAll returns all Casbin policies.
func (r *authPolicyRepository) FindAllRoles() ([]string, error) {
	// return all policies found in the databaseq
	roles := r.enforcer.GetAllRoles()
	fmt.Printf("Roles: %v\n", roles)
	return roles, nil
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
	_, err := r.enforcer.DeleteRolesForUser(userId)
	if err != nil {
		fmt.Printf("Error removing roles for user: %v\n", err)
		return nil, err
	}

	// Add the new role for the user.
	success, err := r.enforcer.AddRoleForUser(userId, roleToApply)
	if err != nil {
		fmt.Printf("Error assigning role to user: %v\n", err)
		return nil, err
	}

	return &success, nil
}

func (r *authPolicyRepository) FindAll() ([][]string, error) {
	// return all policies found in the database
	policies := r.enforcer.GetPolicy()
	return policies, nil
}

func (r *authPolicyRepository) Create(policy models.CasbinRule) error {
	// Add policy to enforcer
	newPolicy, err := r.enforcer.AddPolicy(policy.V0, policy.V1, policy.V2)
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
	removed, err = r.enforcer.RemovePolicy(policy.V0, policy.V1, policy.V2)
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
	removed, err := r.enforcer.RemovePolicy(oldPolicy.V0, oldPolicy.V1, oldPolicy.V2)
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
	addedPolicy, err := r.enforcer.AddPolicy(newPolicy.V0, newPolicy.V1, newPolicy.V2)
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
