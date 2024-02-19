package repository_test

import (
	"fmt"
	"testing"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
)

// Policies
func TestAuthPolicyRepository_Create(t *testing.T) {
	policyToCreate := models.CasbinRule{V0: "admin", V1: "/api/v1/users", V2: "GET"}
	// Test function
	err := testModule.auth.repo.Create(policyToCreate)
	if err != nil {
		t.Errorf("Error creating policy: %v", err)
	}
	// Tear down
	err = testModule.auth.repo.Delete(policyToCreate)
	if err != nil {
		t.Errorf("Error deleting policy: %v", err)
	}
}

func TestAuthPolicyRepository_Delete(t *testing.T) {
	policyToCreate := models.CasbinRule{V0: "admin", V1: "/api/v1/users", V2: "GET"}
	err := testModule.auth.repo.Create(policyToCreate)
	if err != nil {
		t.Errorf("Error creating policy: %v", err)
	}
	// Test function
	err = testModule.auth.repo.Delete(policyToCreate)
	if err != nil {
		t.Fatalf("Error deleting policy: %v", err)
	}
}

func TestAuthPolicyRepository_FindAll(t *testing.T) {
	policy1 := models.CasbinRule{V0: "admin", V1: "/api/v1/users", V2: "create"}
	policy2 := models.CasbinRule{V0: "admin", V1: "/api/v1/users", V2: "update"}
	// Setup
	err := testModule.auth.repo.Create(policy1)
	if err != nil {
		t.Errorf("Error creating policy: %v", err)
	}
	err = testModule.auth.repo.Create(policy2)
	if err != nil {
		t.Errorf("Error creating policy: %v", err)
	}
	// Test function
	policies, err := testModule.auth.repo.FindAll()
	if err != nil {
		t.Errorf("Error finding policies: %v", err)
	}
	if len(policies) != 2 {
		t.Errorf("Expected 2 policies, found %v", len(policies))
	}

	// Iterate through policies
	for _, policy := range policies {
		// Check if policy resource matches policy1 or policy2
		if policy[1] == policy1.V0 {
			// Iterate through policy checking details against policy1
			checkArrayStringPolicyAgainstCasbinRule(policy, policy1, t)
		} else if policy[1] == policy2.V0 {
			checkArrayStringPolicyAgainstCasbinRule(policy, policy2, t)
		}

	}
	// Cleanup
	err = testModule.auth.repo.Delete(policy1)
	if err != nil {
		t.Errorf("Error deleting policy: %v", err)
	}
	err = testModule.auth.repo.Delete(policy2)
	if err != nil {
		t.Errorf("Error deleting policy: %v", err)
	}
}

func TestAuthPolicyRepository_Update(t *testing.T) {
	oldPolicy := models.CasbinRule{V0: "admin", V1: "/api/v1/users", V2: "GET"}
	newPolicy := models.CasbinRule{V0: "admin", V1: "/api/v1/users", V2: "POST"}
	// Setup
	err := testModule.auth.repo.Create(oldPolicy)
	if err != nil {
		t.Errorf("Error creating policy: %v", err)
	}
	// Test function
	err = testModule.auth.repo.Update(oldPolicy, newPolicy)
	if err != nil {
		t.Errorf("Error updating policy: %v", err)
	}

	// Cleanup
	err = testModule.auth.repo.Delete(newPolicy)
	if err != nil {
		t.Errorf("Error deleting policy: %v", err)
	}
}

// Roles
func TestAuthRoleRepository_AssignUserRole(t *testing.T) {
	createdUser, err := testModule.users.repo.Create(&db.User{Email: "ratbag@gmail.com", Password: "password"})
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
	// Test function
	success, err := testModule.auth.repo.AssignUserRole(fmt.Sprint(createdUser.ID), "admin")
	if err != nil {
		t.Errorf("Error assigning role to user: %v", err)
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
	}
	// Cleanup
	err = testModule.users.repo.Delete(int(createdUser.ID))
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}
	success, err = testModule.auth.repo.DeleteRolesForUser(fmt.Sprint(createdUser.ID))
	if err != nil {
		t.Errorf("Error deleting roles for user: %v", err)
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
	}
}
func TestAuthRoleRepository_FindAllRoles(t *testing.T) {

}

func TestAuthRoleRepository_RoleByUserId(t *testing.T) {
}

func TestAuthRoleRepository_Delete(t *testing.T) {
}

func checkArrayStringPolicyAgainstCasbinRule(policy []string, casbinRule models.CasbinRule, t *testing.T) bool {
	if policy[0] != casbinRule.V0 {
		t.Errorf("Expected %v, found %v", casbinRule.V0, policy[0])
	}
	if policy[1] != casbinRule.V1 {
		t.Errorf("Expected %v, found %v", casbinRule.V1, policy[1])
		return false
	}
	if policy[2] != casbinRule.V2 {
		t.Errorf("Expected %v, found %v", casbinRule.V2, policy[2])
		return false
	}
	return true
}
