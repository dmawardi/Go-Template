package service_test

import (
	"fmt"
	"testing"

	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
)

// Policies
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
	// Test function
	roles, err := testModule.auth.serv.FindAllRoles()
	if err != nil {
		t.Errorf("Error finding roles: %v", err)
	}

	// Test if all roles are found (default roles are: admin, moderator, user)
	if len(roles) != 3 {
		t.Errorf("Expected 3 roles, got %d", len(roles))
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

	var tests = []struct {
		roleToApply     string
		expectedSuccess bool
	}{
		{"admin", true},
		{"moderator", true},
		{"user", true},
		{"nonexistentrole", false},
	}

	for _, v := range tests {
		// Test function
		success, err := testModule.auth.serv.AssignUserRole(fmt.Sprint(createdUser.ID), v.roleToApply)

		// If success is expected
		if v.expectedSuccess {
			// Check that there is no error
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
			if role != v.roleToApply {
				t.Errorf("Expected role %s, got %s", v.roleToApply, role)
			}

			// Else if failure expected
		} else {
			// Check that there is an error
			if err == nil {
				t.Errorf("Expected error, got nil")
			}
			if success != nil {
				t.Errorf("Expected success to be nil, got %v", success)
			}
		}
	}

	// Clean up
	err = testModule.users.serv.Delete(int(createdUser.ID))
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}
}

func TestAuthPolicyService_CreateRole(t *testing.T) {
	// Create a user
	createdUser, err := testModule.users.serv.Create(&models.CreateUser{
		Email:    "KurtBangle@gmial.com",
		Password: "password",
		Role:     "admin",
	})
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	var tests = []struct {
		roleToApply     string
		expectedSuccess bool
	}{
		{"squire", true},
		{"knight", true},
		{"zebra", true},
		// Existing role should fail
		{"admin", false},
	}

	for _, v := range tests {
		// Test function
		success, err := testModule.auth.serv.CreateRole(fmt.Sprint(createdUser.ID), v.roleToApply)

		// Check that there is no error
		if err != nil {
			t.Errorf("Error assigning role: %v", err)
		}
		if !*success {
			t.Errorf("Expected success, got %v", success)
		}

		roles, err := testModule.auth.serv.FindAllRoles()
		if err != nil {
			t.Errorf("Error finding roles: %v", err)
		}
		if !helpers.ArrayContainsString(roles, v.roleToApply) {
			t.Errorf("Expected role %s to be created, got %v", v.roleToApply, roles)
		}

		// Test if role is assigned
		role, err := testModule.auth.serv.FindRoleByUserId(int(createdUser.ID))
		if err != nil {
			t.Errorf("Error finding role: %v", err)
		}

		// Test if role is found as required
		if role != v.roleToApply {
			t.Errorf("Expected role %s, got %s", v.roleToApply, role)
		}
	}

	// Clean up
	err = testModule.users.serv.Delete(int(createdUser.ID))
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}
}

// Role inheritance

// func TestAuthPolicyService_FindAllRoleInheritance(t *testing.T) {
// 	// Create a user
// 	createdUser1, err := testModule.users.serv.Create(&models.CreateUser{
// 		Email:    "raylongo@gmial.com",
// 		Password: "password",
// 		Role:     "alien",
// 	})
// 	if err != nil {
// 		t.Errorf("Error creating user: %v", err)
// 	}
// 	createdUser2, err := testModule.users.serv.Create(&models.CreateUser{
// 		Email:    "mongopak@gmial.com",
// 		Password: "password",
// 		Role:     "eevee",
// 	})
// 	if err != nil {
// 		t.Errorf("Error creating user: %v", err)
// 	}
// 	createdUser3, err := testModule.users.serv.Create(&models.CreateUser{
// 		Email:    "sideshowbob@gmial.com",
// 		Password: "password",
// 		Role:     "pikachu",
// 	})
// 	if err != nil {
// 		t.Errorf("Error creating user: %v", err)
// 	}
// 	t.Errorf("Created users. 1:%v\n2:%v\n3.%v", createdUser1, createdUser2, createdUser3)

// 	// Build inheritances
// 	inheritance1 := models.GRecord{Role: createdUser1.Role, InheritsFrom: createdUser2.Role}
// 	inheritance2 := models.GRecord{Role: createdUser2.Role, InheritsFrom: createdUser3.Role}

// 	// Create inheritance
// 	err = testModule.auth.serv.CreateInheritance(inheritance1)
// 	if err != nil {
// 		t.Errorf("Error creating inheritance: %v", err)
// 	}
// 	err = testModule.auth.serv.CreateInheritance(inheritance2)
// 	if err != nil {
// 		t.Errorf("Error creating inheritance: %v", err)
// 	}

// 	// Test function
// 	inheritances, err := testModule.auth.serv.FindAllRoleInheritance()
// 	if err != nil {
// 		t.Errorf("Error finding inheritance: %v", err)
// 	}

// 	// Test if all inheritances are found (2 on top of the default 3 roles)
// 	if len(inheritances) != 5 {
// 		t.Errorf("Expected 2 inheritances, got %d", len(inheritances))
// 	}

// 	// Iterate through inheritances
// 	for _, inheritance := range inheritances {

// 		// Check if role matches with inheritance 1
// 		if inheritance.Role == inheritance1.Role {
// 			if inheritance.InheritsFrom != inheritance1.InheritsFrom {
// 				t.Errorf("Expected inheritance %v, got %v", inheritance1, inheritance)
// 			}

// 			// Else check if matches with inheritance 2
// 		} else if inheritance.Role == inheritance2.Role {
// 			if inheritance.InheritsFrom != inheritance2.InheritsFrom {
// 				t.Errorf("Expected inheritance %v, got %v", inheritance2, inheritance)
// 			}
// 		}
// 	}

// 	// Clean up
// 	err = testModule.users.serv.Delete(int(createdUser1.ID))
// 	if err != nil {
// 		t.Errorf("Error deleting user: %v", err)
// 	}
// 	err = testModule.users.serv.Delete(int(createdUser2.ID))
// 	if err != nil {
// 		t.Errorf("Error deleting user: %v", err)
// 	}
// 	err = testModule.users.serv.Delete(int(createdUser3.ID))
// 	if err != nil {
// 		t.Errorf("Error deleting user: %v", err)
// 	}

// 	// Delete inheritances
// 	err = testModule.auth.serv.DeleteInheritance(inheritance1)
// 	if err != nil {
// 		t.Errorf("Error deleting inheritance: %v", err)
// 	}
// 	err = testModule.auth.serv.DeleteInheritance(inheritance2)
// 	if err != nil {
// 		t.Errorf("Error deleting inheritance: %v", err)
// 	}
// }

// func TestAuthPolicyService_CreateInheritance(t *testing.T) {
// 	// Create a user
// 	createdUser1, err := testModule.users.serv.Create(&models.CreateUser{
// 		Email:    "yummyjam@gmial.com",
// 		Password: "password",
// 		Role:     "user",
// 	})
// 	if err != nil {
// 		t.Errorf("Error creating user: %v", err)
// 	}
// 	createdUser2, err := testModule.users.serv.Create(&models.CreateUser{
// 		Email:    "swaqmmyamy@gmial.com",
// 		Password: "password",
// 		Role:     "admin",
// 	})
// 	if err != nil {
// 		t.Errorf("Error creating user: %v", err)
// 	}

// 	// Test function
// 	inheritance := models.GRecord{Role: "admin", InheritsFrom: "user"}
// 	err = testModule.auth.serv.CreateInheritance(inheritance)
// 	if err != nil {
// 		t.Errorf("Error creating inheritance: %v", err)
// 	}

// 	// Test if inheritance is created
// 	inheritances, err := testModule.auth.serv.FindAllRoleInheritance()
// 	if err != nil {
// 		t.Errorf("Error finding inheritance: %v", err)
// 	}

// 	if len(inheritances) != 1 {
// 		t.Errorf("Expected 1 inheritance, got %d", len(inheritances))
// 	}

// 	// Add naming convention as expected
// 	expectedInheritance := models.GRecord{Role: "role:admin", InheritsFrom: "role:user"}

// 	helpers.CompareObjects(inheritances[0], expectedInheritance, t, []string{"Role", "InheritsFrom"})

// 	// Clean up
// 	// Delete inheritance
// 	app.Auth.Enforcer.RemoveNamedGroupingPolicy("g", []string{expectedInheritance.Role, expectedInheritance.InheritsFrom})
// 	if err != nil {
// 		t.Errorf("Error deleting inheritance: %v", err)
// 	}

// 	err = testModule.users.serv.Delete(int(createdUser1.ID))
// 	if err != nil {
// 		t.Errorf("Error deleting user: %v", err)
// 	}
// 	err = testModule.users.serv.Delete(int(createdUser2.ID))
// 	if err != nil {
// 		t.Errorf("Error deleting user: %v", err)
// 	}
// }

// func TestAuthPolicyService_DeleteInheritance(t *testing.T) {
// 	// Create a user
// 	createdUser1, err := testModule.users.serv.Create(&models.CreateUser{
// 		Email:    "billbellamy@gmial.com",
// 		Password: "password",
// 		Role:     "user",
// 	})
// 	if err != nil {
// 		t.Errorf("Error creating user: %v", err)
// 	}
// 	createdUser2, err := testModule.users.serv.Create(&models.CreateUser{
// 		Email:    "swankingsont@gmial.com",
// 		Password: "password",
// 		Role:     "admin",
// 	})
// 	if err != nil {
// 		t.Errorf("Error creating user: %v", err)
// 	}

// 	// Test function
// 	inheritance := models.GRecord{Role: "admin", InheritsFrom: "user"}
// 	err = testModule.auth.serv.CreateInheritance(inheritance)
// 	if err != nil {
// 		t.Errorf("Error creating inheritance: %v", err)
// 	}

// 	// Test if inheritance is created
// 	inheritances, err := testModule.auth.serv.FindAllRoleInheritance()
// 	if err != nil {
// 		t.Errorf("Error finding inheritance: %v", err)
// 	}

// 	if len(inheritances) != 1 {
// 		t.Errorf("Expected 1 inheritance, got %d", len(inheritances))
// 	}

// 	// Test function
// 	err = testModule.auth.serv.DeleteInheritance(inheritance)
// 	if err != nil {
// 		t.Errorf("Error deleting inheritance: %v", err)
// 	}

// 	// Test if inheritance is deleted
// 	inheritances, err = testModule.auth.serv.FindAllRoleInheritance()
// 	if err != nil {
// 		t.Errorf("Error finding inheritance: %v", err)
// 	}

// 	if len(inheritances) != 0 {
// 		t.Errorf("Expected 0 inheritance, got %d", len(inheritances))
// 	}

//		// Clean up
//		err = testModule.users.serv.Delete(int(createdUser1.ID))
//		if err != nil {
//			t.Errorf("Error deleting user: %v", err)
//		}
//		err = testModule.users.serv.Delete(int(createdUser2.ID))
//		if err != nil {
//			t.Errorf("Error deleting user: %v", err)
//		}
//	}
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
