package service_test

import (
	"fmt"
	"testing"

	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
)

func TestAuthPolicyService_Create(t *testing.T) {
	policyToCreate := models.PolicyRule{
		Resource: "/testResource",
		Action:   "create",
		Role:     "admin",
	}

	// Test function
	err := testModule.auth.serv.Create(policyToCreate)
	if err != nil {
		t.Errorf("Error creating policy: %v", err)
	}

	// Test if policy is created
	policies, err := testModule.auth.serv.FindByResource(policyToCreate.Resource)
	if err != nil {
		t.Errorf("Error finding policy: %v", err)
	}

	if len(policies) != 1 {
		t.Errorf("Expected 1 policy, got %d", len(policies))
	}

	checkPolicyMatch(t, policyToCreate, policies[0])

	// Clean up
	err = testModule.auth.serv.Delete(policyToCreate)
	if err != nil {
		t.Errorf("Error deleting policy: %v", err)
	}
}

func TestAuthPolicyService_Delete(t *testing.T) {
	// Create a policy
	policyToCreate := models.PolicyRule{
		Resource: "/testResource",
		Action:   "create",
		Role:     "admin",
	}
	err := testModule.auth.serv.Create(policyToCreate)
	if err != nil {
		t.Fatalf("Error creating policy: %v", err)
	}

	// Test function
	err = testModule.auth.serv.Delete(policyToCreate)
	if err != nil {
		t.Errorf("Error deleting policy: %v", err)
	}

	// Test if policy is deleted
	policies, err := testModule.auth.serv.FindByResource(policyToCreate.Resource)
	if err != nil {
		t.Errorf("Error finding policy: %v", err)
	}

	if len(policies) != 0 {
		t.Errorf("Expected 0 policy, got %d", len(policies))
	}
}

func TestAuthPolicyService_FindByResource(t *testing.T) {
	// Create a policy
	policyToCreate := models.PolicyRule{
		Resource: "/testResource",
		Action:   "create",
		Role:     "admin",
	}
	err := testModule.auth.serv.Create(policyToCreate)
	if err != nil {
		t.Fatalf("Error creating policy: %v", err)
	}
	// Test function
	policies, err := testModule.auth.serv.FindByResource(policyToCreate.Resource)
	if err != nil {
		t.Errorf("Error finding policy: %v", err)
	}

	if len(policies) != 1 {
		t.Errorf("Expected 1 policy, got %d", len(policies))
	}

	checkPolicyMatch(t, policyToCreate, policies[0])

	// Clean up
	err = testModule.auth.serv.Delete(policyToCreate)
	if err != nil {
		t.Errorf("Error deleting policy: %v", err)
	}
}

func TestAuthPolicyService_FindAll(t *testing.T) {
	// Create a policy
	policy1 := models.PolicyRule{
		Resource: "/testResource",
		Action:   "create",
		Role:     "admin",
	}
	sameResourceDifferentActionPolicy := models.PolicyRule{
		Resource: "/testResource",
		Action:   "update",
		Role:     "admin",
	}
	sameResourceDifferentRolePolicy := models.PolicyRule{
		Resource: "/testResource",
		Action:   "create",
		Role:     "user",
	}
	policy2 := models.PolicyRule{
		Resource: "/nextResource",
		Action:   "create",
		Role:     "admin",
	}
	policiesToCreate := []models.PolicyRule{policy1, sameResourceDifferentActionPolicy, sameResourceDifferentRolePolicy, policy2}
	for _, policy := range policiesToCreate {
		err := testModule.auth.serv.Create(policy)
		if err != nil {
			t.Fatalf("Error creating policy: %v", err)
		}
	}

	// Test function
	policies, err := testModule.auth.serv.FindAll("")
	if err != nil {
		t.Errorf("Error finding policy: %v", err)
	}

	// Test if all policies are found (one for each role)
	if len(policies) != 3 {
		t.Errorf("Expected 1 policy, got %d", len(policies))
	}

	// Check details
	for _, policy := range policies {
		// If it's nextResource, check if it matches policy2
		if policy.Resource == policy2.Resource {
			checkPolicyMatch(t, policy2, policy)

		} else if policy.Resource == policy1.Resource {
			// Else if it's testResource, check what role is to determine which match to check against
			if policy.Role == policy1.Role {
				checkPolicyMatch(t, policy1, policy)
			} else if policy.Role == sameResourceDifferentActionPolicy.Role {
				checkPolicyMatch(t, sameResourceDifferentActionPolicy, policy)
			} else if policy.Role == sameResourceDifferentRolePolicy.Role {
				checkPolicyMatch(t, sameResourceDifferentRolePolicy, policy)
			}

		}

	}
	// Clean up
	for _, policy := range policiesToCreate {
		err = testModule.auth.serv.Delete(policy)
		if err != nil {
			t.Errorf("Error deleting policy: %v", err)
		}
	}
}

func TestAuthPolicyService_Update(t *testing.T) {
	// Create a policy
	policyToCreate := models.PolicyRule{
		Resource: "/testResource",
		Action:   "create",
		Role:     "admin",
	}
	err := testModule.auth.serv.Create(policyToCreate)
	if err != nil {
		t.Fatalf("Error creating policy: %v", err)
	}

	// Test function
	policyToUpdate := models.PolicyRule{
		Resource: "/testResource",
		Action:   "create",
		Role:     "user",
	}
	err = testModule.auth.serv.Update(policyToCreate, policyToUpdate)
	if err != nil {
		t.Errorf("Error updating policy: %v", err)
	}

	// Test if policy is updated
	policies, err := testModule.auth.serv.FindByResource(policyToCreate.Resource)
	if err != nil {
		t.Errorf("Error finding policy: %v", err)
	}

	if len(policies) != 1 {
		t.Errorf("Expected 1 policy, got %d", len(policies))
	}

	checkPolicyMatch(t, policyToUpdate, policies[0])

	// Clean up
	err = testModule.auth.serv.Delete(policyToUpdate)
	if err != nil {
		t.Errorf("Error deleting policy: %v", err)
	}
}

// Roles

func TestAuthPolicyService_FindAllRoles(t *testing.T) {
	// Create a user with a role
	createdUser1, err := testModule.users.serv.Create(&models.CreateUser{
		Email:    "wallyhjango@gmial.com",
		Password: "password",
		Role:     "admin",
	})
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
	createdUser2, err := testModule.users.serv.Create(&models.CreateUser{
		Email:    "whereswally@gmial.com",
		Password: "password",
		Role:     "user",
	})
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Test function
	roles, err := testModule.auth.serv.FindAllRoles()
	if err != nil {
		t.Errorf("Error finding roles: %v", err)
	}

	// Test if all roles are found
	if len(roles) != 2 {
		t.Errorf("Expected 2 roles, got %d", len(roles))
	}

	found := helpers.ArrayContainsString(roles, createdUser1.Role)
	if !found {
		t.Errorf("Expected role %s, got %v", createdUser1.Role, roles)
	}
	found = helpers.ArrayContainsString(roles, createdUser2.Role)
	if !found {
		t.Errorf("Expected role %s, got %v", createdUser2.Role, roles)
	}

	// Clean up
	err = testModule.users.serv.Delete(int(createdUser1.ID))
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}
	err = testModule.users.serv.Delete(int(createdUser2.ID))
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}

	// Check roles are deleted
	roles, err = testModule.auth.serv.FindAllRoles()
	if err != nil {
		t.Errorf("Error finding roles: %v", err)
	}

	if len(roles) != 0 {
		t.Errorf("Expected 0 roles, got %d", len(roles))
	}
}

func TestAuthPolicyService_FindRoleByUserId(t *testing.T) {
	// Create a user with a role
	createdUser, err := testModule.users.serv.Create(&models.CreateUser{
		Email:    "banjo@gmial.com",
		Password: "password",
		Role:     "admin",
	})
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Test function
	role, err := testModule.auth.serv.FindRoleByUserId(int(createdUser.ID))
	if err != nil {
		t.Errorf("Error finding role: %v", err)
	}

	// Test if role is found
	if role != createdUser.Role {
		t.Errorf("Expected role %s, got %s", createdUser.Role, role)
	}

	// Clean up
	err = testModule.users.serv.Delete(int(createdUser.ID))
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}
}

func TestAuthPolicyService_AssignUserRole(t *testing.T) {
	// Create a user
	createdUser, err := testModule.users.serv.Create(&models.CreateUser{
		Email:    "willybongo@gmial.com",
		Password: "password",
		Role:     "admin",
	})
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Test function
	roleToApply := "user"
	success, err := testModule.auth.serv.AssignUserRole(fmt.Sprint(createdUser.ID), roleToApply)
	if err != nil {
		t.Errorf("Error assigning role: %v", err)
	}
	if !*success {
		t.Errorf("Expected success, got %v", success)
	}

	// Test if role is assigned
	role, err := testModule.auth.serv.FindRoleByUserId(int(createdUser.ID))
	if err != nil {
		t.Errorf("Error finding role: %v", err)
	}

	// Test if role is found as required
	if role != roleToApply {
		t.Errorf("Expected role %s, got %s", roleToApply, role)
	}

	// Clean up
	err = testModule.users.serv.Delete(int(createdUser.ID))
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}
}

func checkPolicyMatch(t *testing.T, created models.PolicyRule, found models.PolicyRuleCombinedActions) {
	if found.Resource != created.Resource {
		t.Errorf("Expected resource %s, got %s", created.Resource, found.Resource)
	}
	if found.Role != created.Role {
		t.Errorf("Expected role %s, got %s", created.Role, found.Role)
	}
	containsPolicy := helpers.ArrayContainsString(found.Action, created.Action)
	if !containsPolicy {
		t.Errorf("Expected action %s, got %s", created.Action, found.Action)
	}
}
