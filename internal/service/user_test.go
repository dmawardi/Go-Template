package service_test

import (
	"testing"

	"github.com/dmawardi/Go-Template/internal/db"
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

	// Test function
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

	// Check if default role applied
	if createdUser.Role != "role:user" {
		t.Errorf("created user has incorrect role: expected 'role:user', got %s", createdUser.Role)
	}

	// Clean up: Delete created user
	result := testModule.dbClient.Delete(&db.User{}, createdUser.ID)
	if result.Error != nil {
		t.Fatalf("failed to delete created user: %v", result.Error)
	}
}

func TestUserService_FindById(t *testing.T) {
	// Build test user
	userToCreate := &db.User{
		Username: "Jabar",
		Email:    "juba@ymail.com",
		Password: "password",
		Name:     "Bamba",
	}

	// Create user
	createdUser, err := helpers.HashPassAndGenerateUserInDb(userToCreate, testModule.dbClient, t)
	if err != nil {
		t.Fatalf("failed to create test user for find by id user service testr: %v", err)
	}

	// Test function
	// Find created user by id
	foundUser, err := testModule.users.serv.FindById(int(userToCreate.ID))
	if err != nil {
		t.Fatalf("failed to find created user: %v", err)
	}
	// Check details
	helpers.CompareObjects(foundUser, userToCreate, t, []string{"ID", "Name", "Username", "Email"})

	// Clean up: Delete created user
	result := testModule.dbClient.Delete(createdUser)
	if result.Error != nil {
		t.Fatalf("failed to delete created user: %v", result.Error)
	}
}

func TestUserService_FindByEmail(t *testing.T) {
	// Build test user
	userToCreate := &db.User{
		Username: "Jabar",
		Email:    "juba@findmymail.com",
		Password: "password",
		Name:     "Bamba",
	}

	// Create user
	createdUser, err := helpers.HashPassAndGenerateUserInDb(userToCreate, testModule.dbClient, t)
	if err != nil {
		t.Fatalf("failed to create test user for find by id user service testr: %v", err)
	}
	// Find created user by id
	foundUser, err := testModule.users.serv.FindByEmail(createdUser.Email)
	if err != nil {
		t.Fatalf("failed to find created user: %v", err)
	}

	helpers.CompareObjects(foundUser, userToCreate, t, []string{"ID", "Name", "Username", "Email"})

	// Clean up: Delete created user
	result := testModule.dbClient.Delete(createdUser)
	if result.Error != nil {
		t.Fatalf("failed to delete created user: %v", result.Error)
	}
}

func TestUserService_Delete(t *testing.T) {
	createdUser, err := helpers.HashPassAndGenerateUserInDb(&db.User{
		Username: "Jabar",
		Email:    "dollar@ymail.com",
		Password: "password",
		Name:     "Crimson",
	}, testModule.dbClient, t)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Test function
	// Delete the created user
	err = testModule.users.serv.Delete(int(createdUser.ID))
	if err != nil {
		t.Fatalf("failed to delete created user: %v", err)
	}

	// Check to see if user has been deleted
	foundUser := &db.User{}
	result := testModule.dbClient.First(foundUser, createdUser.ID)
	if result.Error == nil {
		t.Fatal("Expected an error but got none")
	}

	// Clean up: Delete created user
	result = testModule.dbClient.Delete(createdUser)
	if result.Error != nil {
		t.Fatalf("failed to delete created user: %v", result.Error)
	}
}

func TestUserService_BulkDelete(t *testing.T) {
}

func TestUserService_Update(t *testing.T) {
	// Create test user
	createdUser, err := helpers.HashPassAndGenerateUserInDb(&db.User{
		Username: "Jabar",
		Email:    "update-test@ymail.com",
		Password: "password",
		Name:     "Crimson",
	}, testModule.dbClient, t)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Create updated details in update user DTO
	userToUpdate := &models.UpdateUser{Username: "Hullabaloo",
		Email:    "update-twist@ymail.com",
		Password: "squash",
		Name:     "Crazy"}

	// Update the created user
	updatedUser, err := testModule.users.serv.Update(int(createdUser.ID), userToUpdate)
	if err != nil {
		t.Fatalf("failed to update created user in service: %v", err)
	}

	helpers.CompareObjects(updatedUser, userToUpdate, t, []string{"Username", "Email", "Name"})

	// Verify that the user ID matches as expected
	if createdUser.ID != updatedUser.ID {
		t.Error("created user should have same ID as updated user")
	}

	// Clean up: Delete created user
	result := testModule.dbClient.Delete(createdUser)
	if result.Error != nil {
		t.Fatalf("failed to delete created user: %v", result.Error)
	}
}

func TestUserService_FindAll(t *testing.T) {
	createdUser1, err := helpers.HashPassAndGenerateUserInDb(&db.User{
		Username: "Joe",
		Email:    "crazy@gmail.com",
		Password: "password",
		Name:     "Bamba",
	}, testModule.dbClient, t)
	if err != nil {
		t.Fatalf("failed to create test user1: %v", err)
	}

	createdUser2, err := helpers.HashPassAndGenerateUserInDb(&db.User{
		Username: "Jabar",
		Email:    "scuba@ymail.com",
		Password: "password",
		Name:     "Bamba",
	}, testModule.dbClient, t)
	if err != nil {
		t.Fatalf("failed to create test user2: %v", err)
	}

	users, err := testModule.users.serv.FindAll(10, 0, "", []models.QueryConditionParameters{})
	if err != nil {
		t.Fatalf("failed to find all: %v", err)
	}

	// Make sure both users are in database
	if len(*users.Data) != 2 {
		t.Errorf("Length of []users is not as expected. Got: %+v", *users.Data)
	}

	// Iterate through results checking user 1 and 2 results
	for _, u := range *users.Data {
		// If it's the first user
		if int(u.ID) == int(createdUser1.ID) {
			helpers.CompareObjects(u, createdUser1, t, []string{"ID", "Name", "Username", "Email"})
		} else {
			helpers.CompareObjects(u, createdUser2, t, []string{"ID", "Name", "Username", "Email"})
		}
	}

	// Clean up created users
	usersToDelete := []db.User{{ID: createdUser1.ID}, {ID: createdUser2.ID}}
	result := testModule.dbClient.Delete(usersToDelete)
	if result.Error != nil {
		t.Fatalf("failed to delete created users: %v", result.Error)
	}
}

func TestUserService_ResetPasswordAndSendEmail(t *testing.T) {
	// Create test user
	createdUser, err := helpers.HashPassAndGenerateUserInDb(&db.User{
		Username: "Jabar",
		Email:    "update-smash@ymail.com",
		Password: "password",
		Name:     "Crimson",
	}, testModule.dbClient, t)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Test function
	// Reset password and send email
	err = testModule.users.serv.ResetPasswordAndSendEmail(createdUser.Email)
	if err != nil {
		t.Fatalf("failed to reset password and send email: %v", err)
	}

	// Clean up: Delete created user
	result := testModule.dbClient.Delete(createdUser)
	if result.Error != nil {
		t.Fatalf("failed to delete created user: %v", result.Error)
	}
}

func TestUserService_LoginUser(t *testing.T) {
	// create password
	password := "password"
	// Create test user
	createdUser, err := helpers.HashPassAndGenerateUserInDb(&db.User{
		Username: "Jabar",
		Email:    "update-zegula@ymail.com",
		Password: password,
		Name:     "Crimson",
	}, testModule.dbClient, t)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Test function
	token, err := testModule.users.serv.LoginUser(&models.Login{Email: createdUser.Email, Password: password})
	if err != nil {
		t.Fatalf("failed to login user: %v", err)
	}

	// Verify that the token is not empty
	if token == "" {
		t.Error("token should not be empty")
	}

	// Clean up: Delete created user
	result := testModule.dbClient.Delete(createdUser)
	if result.Error != nil {
		t.Fatalf("failed to delete created user: %v", result.Error)
	}
}

func TestUserService_CheckPasswordMatch(t *testing.T) {
	// Create test user
	createdUser, err := helpers.HashPassAndGenerateUserInDb(&db.User{
		Username: "Jabar",
		Email:    "update-cigz@ymail.com",
		Password: "password",
		Name:     "Crimson",
	}, testModule.dbClient, t)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Test function
	// Check password match
	matchFound := testModule.users.serv.CheckPasswordMatch(int(createdUser.ID), []byte("password"))
	if !matchFound {
		t.Error("password match should be found")
	}

	// Clean up: Delete created user
	result := testModule.dbClient.Delete(createdUser)
	if result.Error != nil {
		t.Fatalf("failed to delete created user: %v", result.Error)
	}
}

func TestUserService_VerifyEmailCode(t *testing.T) {
	// Create test user
	createdUser, err := helpers.HashPassAndGenerateUserInDb(&db.User{
		Username: "Jabar",
		Email:    "update-wgz@ymail.com",
		Password: "password",
		Name:     "Crimson",
	}, testModule.dbClient, t)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	// Generate verification code and set expiry
	helpers.GenerateVerificationCodeAndSetExpiry(createdUser)
	// Update user
	result := testModule.dbClient.Save(createdUser)
	if result.Error != nil {
		t.Fatalf("failed to update created user: %v", result.Error)
	}

	// Test function
	// Verify email code
	err = testModule.users.serv.VerifyEmailCode(createdUser.VerificationCode)
	if err != nil {
		t.Fatalf("failed to verify email code: %v", err)
	}

	// Clean up: Delete created user
	result = testModule.dbClient.Delete(createdUser)
	if result.Error != nil {
		t.Fatalf("failed to delete created user: %v", result.Error)
	}
}

func TestUserService_ResendEmailVerification(t *testing.T) {
	// Create test user
	createdUser, err := helpers.HashPassAndGenerateUserInDb(&db.User{
		Username: "Jabar",
		Email:    "update-tadow@ymail.com",
		Password: "password",
		Name:     "Crimson",
	}, testModule.dbClient, t)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Test function
	// Resend email verification
	err = testModule.users.serv.ResendEmailVerification(int(createdUser.ID))
	if err != nil {
		t.Fatalf("failed to resend email verification: %v", err)
	}

	// Clean up: Delete created user
	result := testModule.dbClient.Delete(createdUser)
	if result.Error != nil {
		t.Fatalf("failed to delete created user: %v", result.Error)
	}
}
