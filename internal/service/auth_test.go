package service_test

import (
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

// func TestAuthPolicyService_FindAll(t *testing.T) {
// 	// Create a policy
// 	policyToCreate := models.PolicyRule{
// 		Resource: "/testResource",
// 		Action:   "create",
// 		Role:     "admin",
// 	}
// 	err := testModule.auth.serv.Create(policyToCreate)
// 	if err != nil {
// 		t.Fatalf("Error creating policy: %v", err)
// 	}

// 	// Test function
// 	policies, err := testModule.auth.serv.FindAll("")
// 	if err != nil {
// 		t.Errorf("Error finding policy: %v", err)
// 	}

// 	if len(policies) != 1 {
// 		t.Errorf("Expected 1 policy, got %d", len(policies))
// 	}

// 	checkPolicyMatch(t, policyToCreate, policies[0])

// 	// Clean up
// 	err = testModule.auth.serv.Delete(policyToCreate)
// 	if err != nil {
// 		t.Errorf("Error deleting policy: %v", err)
// 	}
// }

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
