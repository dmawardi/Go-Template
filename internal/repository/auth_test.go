package repository_test

import (
	"testing"

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
	// Setup
	err := testModule.auth.repo.Create(models.CasbinRule{V0: "admin", V1: "/api/v1/users", V2: "GET"})
	if err != nil {
		t.Errorf("Error creating policy: %v", err)
	}
	// Test function

	// Tear down
}
