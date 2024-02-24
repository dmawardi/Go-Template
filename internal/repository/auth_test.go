package repository_test

import (
	"errors"
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
func TestAuthPolicyRepository_AssignUserRole(t *testing.T) {
	createdUser, err := testModule.users.repo.Create(&db.User{Email: "ratbag@gmail.com", Password: "password"})
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
	// Test function
	success, err := testModule.auth.repo.AssignUserRole(fmt.Sprint(createdUser.ID), "role:admin")
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
func TestAuthPolicyRepository_FindAllRoles(t *testing.T) {
	createdUser1, err := testModule.users.repo.Create(&db.User{Email: "catman@gmail.com", Password: "password"})
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
	createdUser2, err := testModule.users.repo.Create(&db.User{Email: "dogman@gmail.com", Password: "password"})
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	role1 := "role:admin"
	role2 := "role:user"
	// Setup
	success, err := testModule.auth.repo.AssignUserRole(fmt.Sprint(createdUser1.ID), role1)
	if err != nil {
		t.Errorf("Error assigning role to user: %v", err)
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
	}
	success, err = testModule.auth.repo.AssignUserRole(fmt.Sprint(createdUser2.ID), role2)
	if err != nil {
		t.Errorf("Error assigning role to user: %v", err)
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
	}
	// Test function
	roles, err := testModule.auth.repo.FindAllRoles()
	if err != nil {
		t.Errorf("Error finding roles: %v", err)
	}
	if len(roles) != 2 {
		t.Errorf("Expected 2 roles, found %v", len(roles))
	}

	// Cleanup
	err = testModule.users.repo.Delete(int(createdUser1.ID))
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}
	err = testModule.users.repo.Delete(int(createdUser2.ID))
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}
	success, err = testModule.auth.repo.DeleteRolesForUser(fmt.Sprint(createdUser1.ID))
	if err != nil {
		t.Errorf("Error deleting roles for user: %v", err)
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
	}
	success, err = testModule.auth.repo.DeleteRolesForUser(fmt.Sprint(createdUser2.ID))
	if err != nil {
		t.Errorf("Error deleting roles for user: %v", err)
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
	}
}
func TestAuthPolicyRepository_RoleByUserId(t *testing.T) {
	// Create user
	createdUser, err := testModule.users.repo.Create(&db.User{Email: "pikachi@gmail.com", Password: "password"})
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
	// Setup
	success, err := testModule.auth.repo.AssignUserRole(fmt.Sprint(createdUser.ID), "role:admin")
	if err != nil {
		t.Errorf("Error assigning role to user: %v", err)
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
	}

	// Test function
	role, err := testModule.auth.repo.FindRoleByUserId(fmt.Sprint(createdUser.ID))
	if err != nil {
		t.Errorf("Error finding role: %v", err)
	}
	if role != "role:admin" {
		t.Errorf("Expected admin, found %v", role)
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
func TestAuthRoleRepository_DeleteRolesForUser(t *testing.T) {
	// Create user
	createdUser, err := testModule.users.repo.Create(&db.User{Email: "pikachu@gmail.com", Password: "password"})
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
	// Setup
	success, err := testModule.auth.repo.AssignUserRole(fmt.Sprint(createdUser.ID), "role:admin")
	if err != nil {
		t.Errorf("Error assigning role to user: %v", err)
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
	}
	// Test function
	success, err = testModule.auth.repo.DeleteRolesForUser(fmt.Sprint(createdUser.ID))
	if err != nil {
		t.Errorf("Error deleting roles for user: %v", err)
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
	}

	// Cleanup
	err = testModule.users.repo.Delete(int(createdUser.ID))
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}
}

// // Role inheritances
// func TestAuthPolicyRepository_CreateInheritance(t *testing.T) {
// 	createdUser1, createdUser1Role := createUserAndSetRole(db.User{Email: "smartman@gmail.com", Password: "password"}, "admin", t)
// 	createdUser2, createdUser2Role := createUserAndSetRole(db.User{Email: "fartman@gmail.com", Password: "password"}, "user", t)
// 	if createdUser1 == nil || createdUser2 == nil {
// 		t.Errorf("Error creating users")
// 	}

// 	// Test function
// 	err := testModule.auth.repo.CreateInheritance(models.G2Record{Role: createdUser1Role, InheritsFrom: createdUser2Role})
// 	if err != nil {
// 		t.Errorf("Error adding role inheritance: %v", err)
// 	}

// 	// Cleanup
// 	err = testModule.auth.repo.DeleteInheritance(models.G2Record{Role: createdUser1Role, InheritsFrom: createdUser2Role})
// 	if err != nil {
// 		t.Errorf("Error deleting role inheritance: %v", err)
// 	}

// 	err = deleteUserAndRole(createdUser1, t)
// 	if err != nil {
// 		t.Errorf("Error deleting user: %v", err)
// 	}
// 	err = deleteUserAndRole(createdUser2, t)
// 	if err != nil {
// 		t.Errorf("Error deleting user: %v", err)
// 	}
// }

// func TestAuthPolicyRepository_DeleteInheritance(t *testing.T) {
// 	createdUser1, createdUser1Role := createUserAndSetRole(db.User{Email: "sparta@gmail.com", Password: "password"}, "admin", t)
// 	createdUser2, createdUser2Role := createUserAndSetRole(db.User{Email: "wingman@gmail.com", Password: "password"}, "user", t)
// 	if createdUser1 == nil || createdUser2 == nil {
// 		t.Fatalf("Error creating users")
// 	}
// 	err := testModule.auth.repo.CreateInheritance(models.G2Record{Role: createdUser1Role, InheritsFrom: createdUser2Role})
// 	if err != nil {
// 		t.Fatalf("Error adding role inheritance: %v", err)
// 	}

// 	// Test function
// 	err = testModule.auth.repo.DeleteInheritance(models.G2Record{Role: createdUser1Role, InheritsFrom: createdUser2Role})
// 	if err != nil {
// 		t.Errorf("Error deleting role inheritance: %v", err)
// 	}

// 	// Cleanup
// 	err = deleteUserAndRole(createdUser1, t)
// 	if err != nil {
// 		t.Errorf("Error deleting user: %v", err)
// 	}
// 	err = deleteUserAndRole(createdUser2, t)
// 	if err != nil {
// 		t.Errorf("Error deleting user: %v", err)
// 	}
// }

// func TestAuthPolicyRepository_FindAllInheritances(t *testing.T) {
// 	// Setup
// 	//  Create users
// 	createdUser1, createdUser1Role := createUserAndSetRole(db.User{Email: "scoresetlle@gmail.com", Password: "password"}, "admin", t)
// 	createdUser2, createdUser2Role := createUserAndSetRole(db.User{Email: "scaryscong@gmail.com", Password: "password"}, "user", t)
// 	createdUser3, createdUser3Role := createUserAndSetRole(db.User{Email: "snailman@smail.com", Password: "password"}, "snail", t)
// 	if createdUser1 == nil || createdUser2 == nil || createdUser3 == nil {
// 		t.Fatalf("Error creating users")
// 	}
// 	// Create inheritances
// 	inheritance1 := models.G2Record{Role: createdUser1Role, InheritsFrom: createdUser2Role}
// 	inheritance2 := models.G2Record{Role: createdUser2Role, InheritsFrom: createdUser3Role}
// 	err := testModule.auth.repo.CreateInheritance(inheritance1)
// 	if err != nil {
// 		t.Fatalf("Error adding role inheritance: %v", err)
// 	}
// 	err = testModule.auth.repo.CreateInheritance(inheritance2)
// 	if err != nil {
// 		t.Fatalf("Error adding role inheritance: %v", err)
// 	}

// 	// Test function
// 	inheritances, err := testModule.auth.repo.FindAllRoleInheritance()
// 	if err != nil {
// 		t.Errorf("Error finding role inheritances: %v", err)
// 	}
// 	if len(inheritances) != 2 {
// 		t.Errorf("Expected 2 inheritances, found %v", len(inheritances))
// 	}

// 	// Check details of each inheritance
// 	for _, inheritance := range inheritances {
// 		// If match found to inheritance 1
// 		if inheritance[0] == inheritance1.Role {
// 			// Check if inherits from matches as well
// 			if inheritance[1] != inheritance1.InheritsFrom {
// 				t.Errorf("Expected %v, found %v", inheritance1.InheritsFrom, inheritance[1])
// 			}

// 			// Else if match found to inheritance 2
// 		} else if inheritance[0] == inheritance2.Role {
// 			// Check if inherits from matches as well
// 			if inheritance[1] != inheritance2.InheritsFrom {
// 				t.Errorf("Expected %v, found %v", inheritance2.InheritsFrom, inheritance[1])
// 			}

// 		}
// 	}

// 	// Cleanup
// 	err = testModule.auth.repo.DeleteInheritance(inheritance1)
// 	if err != nil {
// 		t.Errorf("Error deleting role inheritance: %v", err)
// 	}
// 	err = testModule.auth.repo.DeleteInheritance(inheritance2)
// 	if err != nil {
// 		t.Errorf("Error deleting role inheritance: %v", err)
// 	}
// 	err = deleteUserAndRole(createdUser1, t)
// 	if err != nil {
// 		t.Errorf("Error deleting user: %v", err)
// 	}
// 	err = deleteUserAndRole(createdUser2, t)
// 	if err != nil {
// 		t.Errorf("Error deleting user: %v", err)
// 	}
// 	err = deleteUserAndRole(createdUser3, t)
// 	if err != nil {
// 		t.Errorf("Error deleting user: %v", err)
// 	}
// }

// Creates a user and assigns a role. Returns the user and the role if successful
func createUserAndSetRole(user db.User, role string, t *testing.T) (*db.User, string) {
	createdUser, err := testModule.users.repo.Create(&user)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
		return nil, ""
	}
	success, err := testModule.auth.repo.AssignUserRole(fmt.Sprint(createdUser.ID), role)
	if err != nil {
		t.Errorf("Error assigning role to user: %v", err)
		return nil, ""
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
		return nil, ""
	}
	return createdUser, role
}

func deleteUserAndRole(user *db.User, t *testing.T) error {
	err := testModule.users.repo.Delete(int(user.ID))
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
		return err
	}
	success, err := testModule.auth.repo.DeleteRolesForUser(fmt.Sprint(user.ID))
	if err != nil {
		t.Errorf("Error deleting roles for user: %v", err)
		return err
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
		// Return new error
		return errors.New("Error deleting roles for user")
	}
	return nil
}

// Helper functions
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
