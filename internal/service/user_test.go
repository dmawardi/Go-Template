package service_test

import (
	"testing"

	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_Create(t *testing.T) {
	userToCreate := &models.CreateUser{
		Name:     "Wigwam",
		Username: "Celebration",
		Email:    "wallow@smail.com",
		Password: "HoolaHoops",
	}

	createdUser, err := testModule.users.serv.Create(userToCreate)
	if err != nil {
		t.Fatalf("Failed to create user in service test: %v", err)
	}

	// Verify that the created user has an ID
	if createdUser.ID == 0 {
		t.Error("created user should have an ID")
	}

	// Compare objects
	fieldsToCheck := []string{"Name", "Username", "Email"}
	helpers.CompareObjects(createdUser, userToCreate, t, fieldsToCheck)

	// Verify that the created user has a hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte(userToCreate.Password)); err != nil {
		t.Errorf("created user has incorrect password hash: %v", err)
	}

	// Clean up: Delete created user
	testModule.dbClient.Delete(createdUser)
}

// func TestUserService_FindById(t *testing.T) {
// 	// Build test user
// 	userToCreate := &db.User{
// 		Username: "Jabar",
// 		Email:    "juba@ymail.com",
// 		Password: "password",
// 		Name:     "Bamba",
// 	}

// 	// Create user
// 	createdUser, err := hashPassAndGenerateUserInDb(userToCreate, t)
// 	if err != nil {
// 		t.Fatalf("failed to create test user for find by id user service testr: %v", err)
// 	}
// 	// Find created user by id
// 	foundUser, err := testModule.users.serv.FindById(int(userToCreate.ID))
// 	if err != nil {
// 		t.Fatalf("failed to find created user: %v", err)
// 	}

// 	// Verify that the found user matches the original user
// 	if foundUser.ID != userToCreate.ID {
// 		t.Errorf("found createdUser has incorrect ID: expected %d, got %d", userToCreate.ID, foundUser.ID)
// 	}
// 	if foundUser.Email != userToCreate.Email {
// 		t.Errorf("found createdUser has incorrect email: expected %s, got %s", userToCreate.Email, foundUser.Email)
// 	}

// 	// Clean up: Delete created user
// 	testModule.dbClient.Delete(createdUser)
// }

// func TestUserService_FindByEmail(t *testing.T) {
// 	// Build test user
// 	userToCreate := &db.User{
// 		Username: "Jabar",
// 		Email:    "juba@findmymail.com",
// 		Password: "password",
// 		Name:     "Bamba",
// 	}

// 	// Create user
// 	createdUser, err := hashPassAndGenerateUserInDb(userToCreate, t)
// 	if err != nil {
// 		t.Fatalf("failed to create test user for find by id user service testr: %v", err)
// 	}
// 	// Find created user by id
// 	foundUser, err := testModule.users.serv.FindByEmail(createdUser.Email)
// 	if err != nil {
// 		t.Fatalf("failed to find created user: %v", err)
// 	}

// 	// Verify that the found user matches the original user
// 	if foundUser.ID != createdUser.ID {
// 		t.Errorf("found createdUser has incorrect ID: expected %d, got %d", userToCreate.ID, foundUser.ID)
// 	}
// 	if foundUser.Email != createdUser.Email {
// 		t.Errorf("found createdUser has incorrect email: expected %s, got %s", userToCreate.Email, foundUser.Email)
// 	}

// 	// Clean up: Delete created user
// 	testModule.dbClient.Delete(createdUser)
// }

// func TestUserService_Delete(t *testing.T) {
// 	createdUser, err := hashPassAndGenerateUserInDb(&db.User{
// 		Username: "Jabar",
// 		Email:    "dollar@ymail.com",
// 		Password: "password",
// 		Name:     "Crimson",
// 	}, t)
// 	if err != nil {
// 		t.Fatalf("failed to create test user: %v", err)
// 	}

// 	// Delete the created user
// 	err = testModule.users.serv.Delete(int(createdUser.ID))
// 	if err != nil {
// 		t.Fatalf("failed to delete created user: %v", err)
// 	}

// 	// Check to see if user has been deleted
// 	_, err = testModule.repo.FindById(int(createdUser.ID))
// 	if err == nil {
// 		t.Fatal("Expected an error but got none")
// 	}
// }

// func TestUserService_Update(t *testing.T) {
// 	// Create test user
// 	createdUser, err := hashPassAndGenerateUserInDb(&db.User{
// 		Username: "Jabar",
// 		Email:    "update-test@ymail.com",
// 		Password: "password",
// 		Name:     "Crimson",
// 	}, t)
// 	if err != nil {
// 		t.Fatalf("failed to create test user: %v", err)
// 	}

// 	// Create updated details in update user DTO
// 	userToUpdate := &models.UpdateUser{Username: "Hullabaloo",
// 		Email:    "update-twist@ymail.com",
// 		Password: "squash",
// 		Name:     "Crazy"}

// 	// Update the created user
// 	updatedUser, err := testModule.users.serv.Update(int(createdUser.ID), userToUpdate)
// 	if err != nil {
// 		t.Fatalf("failed to update created user in service: %v", err)
// 	}

// 	// Verify that the created user;s ID is as expected
// 	if createdUser.ID != updatedUser.ID {
// 		t.Error("created user should have same ID as updated user")
// 	}

// 	// Verify other details
// 	if userToUpdate.Email != updatedUser.Email {
// 		t.Errorf("Updated user does not have exp. updated email of: %v. Value: %v", userToUpdate.Email, updatedUser.Email)
// 	}
// 	if userToUpdate.Name != updatedUser.Name {
// 		t.Errorf("Updated user does not have exp. updated email of: %v. Value: %v", userToUpdate.Name, updatedUser.Name)
// 	}
// 	if userToUpdate.Username != updatedUser.Username {
// 		t.Errorf("Updated user does not have exp. updated email of: %v. Value: %v", userToUpdate.Username, updatedUser.Username)
// 	}

// 	// Clean up: Delete created user
// 	testModule.dbClient.Delete(updatedUser)
// }

// func TestUserService_FindAll(t *testing.T) {
// 	createdUser1, err := hashPassAndGenerateUserInDb(&db.User{
// 		Username: "Joe",
// 		Email:    "crazy@gmail.com",
// 		Password: "password",
// 		Name:     "Bamba",
// 	}, t)
// 	if err != nil {
// 		t.Fatalf("failed to create test user1: %v", err)
// 	}

// 	createdUser2, err := hashPassAndGenerateUserInDb(&db.User{
// 		Username: "Jabar",
// 		Email:    "scuba@ymail.com",
// 		Password: "password",
// 		Name:     "Bamba",
// 	}, t)
// 	if err != nil {
// 		t.Fatalf("failed to create test user2: %v", err)
// 	}

// 	users, err := testModule.users.serv.FindAll(10, 0, "", []string{})
// 	if err != nil {
// 		t.Fatalf("failed to find all: %v", err)
// 	}

// 	// Make sure both users are in database
// 	if len(*users.Data) != 2 {
// 		t.Errorf("Length of []users is not as expected. Got: %v", len(*users.Data))
// 	}

// 	// Iterate through results checking user 1 and 2 results
// 	for _, u := range *users.Data {
// 		// If it's the first user
// 		if int(u.ID) == int(createdUser1.ID) {
// 			// check details of first created user
// 			if createdUser1.Email != u.Email {
// 				t.Errorf("Email of user1 doesn't match. Got: %v, expected %v", u.Email, createdUser1.Email)
// 			}
// 			if createdUser1.Username != u.Username {
// 				t.Errorf("Email of user1 doesn't match. Got: %v, expected %v", u.Username, createdUser1.Username)
// 			}
// 			if createdUser1.Name != u.Name {
// 				t.Errorf("Email of user1 doesn't match. Got: %v, expected %v", u.Name, createdUser1.Name)
// 			}
// 		} else {
// 			// check details of second user
// 			if createdUser2.Email != u.Email {
// 				t.Errorf("Email of user1 doesn't match. Got: %v, expected %v", u.Email, createdUser2.Email)
// 			}
// 			if createdUser2.Username != u.Username {
// 				t.Errorf("Email of user1 doesn't match. Got: %v, expected %v", u.Username, createdUser2.Username)
// 			}
// 			if createdUser2.Name != u.Name {
// 				t.Errorf("Email of user1 doesn't match. Got: %v, expected %v", u.Name, createdUser2.Name)
// 			}

// 		}
// 	}
// 	// Clean up created users
// 	usersToDelete := []db.User{{ID: createdUser1.ID}, {ID: createdUser2.ID}}
// 	testModule.dbClient.Delete(usersToDelete)
// }
