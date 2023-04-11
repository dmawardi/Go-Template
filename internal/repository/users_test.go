package repository_test

import (
	"fmt"
	"testing"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/repository"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type testDbRepo struct {
	dbClient *gorm.DB
	repo     repository.UserRepository
}

var testConnection testDbRepo

func init() {
	setupDatabase()
}

func setupDatabase() {
	// Open a new, temporary database for testing
	dbClient, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		fmt.Errorf("failed to open database: %v", err)
	}

	// Migrate the database schema
	if err := dbClient.AutoMigrate(&db.User{}); err != nil {
		fmt.Errorf("failed to migrate database schema: %v", err)
	}
	// Create a new user repository
	repo := repository.NewUserRepository(dbClient)

	// Setup test connection
	testConnection.dbClient = dbClient
	testConnection.repo = repo
}

func TestUserRepository_Create(t *testing.T) {
	// Build hashed password from user password input
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), 10)
	if err != nil {
		t.Fatalf("failed to encrypt password: %v", err)
	}

	// Create a new user
	user := &db.User{
		Email: "test@example.com",
		// Imitate bcrypt encryption from user service
		Password: string(hashedPassword),
	}
	createdUser, err := testConnection.repo.Create(user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Verify that the created user has an ID
	if createdUser.ID == 0 {
		t.Error("created user should have an ID")
	}

	// Verify that the created user has a hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte("password")); err != nil {
		t.Errorf("created user has incorrect password hash: %v", err)
	}

	// Attempt to create duplicate user
	duplicateUser := &db.User{
		Email: "test@example.com",
		// Imitate bcrypt encryption from user service
		Password: string(hashedPassword),
	}
	_, err = testConnection.repo.Create(duplicateUser)
	if err == nil {
		t.Fatalf("Creating duplicate email should have failed but it didn't: %v", err)
	}

	// Clean up: Delete created user
	testConnection.dbClient.Delete(createdUser)
	// In case duplicate was created, delete
	testConnection.dbClient.Delete(duplicateUser)
}

func TestUserRepository_FindById(t *testing.T) {
	createdUser, err := hashPassAndGenerateUserInDb(&db.User{
		Username: "Jabar",
		Email:    "juba@ymail.com",
		Password: "password",
		Name:     "Bamba",
	}, t)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	foundUser, err := testConnection.repo.FindById(int(createdUser.ID))
	if err != nil {
		t.Fatalf("failed to find created user: %v", err)
	}

	// Verify that the found user matches the original user
	if foundUser.ID != createdUser.ID {
		t.Errorf("found createdUser has incorrect ID: expected %d, got %d", createdUser.ID, foundUser.ID)
	}
	if foundUser.Email != createdUser.Email {
		t.Errorf("found createdUser has incorrect email: expected %s, got %s", createdUser.Email, foundUser.Email)
	}
	// check default role applied
	if foundUser.Role != "user" {
		t.Errorf("found createdUser has incorrect default role: expected user, got %s", foundUser.Email)
	}

	// Clean up: Delete created user
	testConnection.dbClient.Delete(createdUser)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	createdUser, err := hashPassAndGenerateUserInDb(&db.User{
		Username: "Jabar",
		Email:    "elon@ymail.com",
		Password: "password",
		Name:     "Bamba",
	}, t)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	foundUser, err := testConnection.repo.FindByEmail(createdUser.Email)
	if err != nil {
		t.Fatalf("failed to find created user: %v", err)
	}

	// Verify that the found user matches the original user
	if foundUser.ID != createdUser.ID {
		t.Errorf("found createdUser has incorrect ID: expected %d, got %d", createdUser.ID, foundUser.ID)
	}
	if foundUser.Email != createdUser.Email {
		t.Errorf("found createdUser has incorrect email: expected %s, got %s", createdUser.Email, foundUser.Email)
	}
	// check default role applied
	if foundUser.Role != "user" {
		t.Errorf("found createdUser has incorrect default role: expected user, got %s", foundUser.Email)
	}

	// Clean up: Delete created user
	testConnection.dbClient.Delete(createdUser)
}

func TestUserRepository_Delete(t *testing.T) {
	createdUser, err := hashPassAndGenerateUserInDb(&db.User{
		Username: "Jabar",
		Email:    "delete@ymail.com",
		Password: "password",
		Name:     "Bamba",
	}, t)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Delete the created user
	err = testConnection.repo.Delete(int(createdUser.ID))
	if err != nil {
		t.Fatalf("failed to delete created user: %v", err)
	}

	// Check to see if user has been deleted
	_, err = testConnection.repo.FindById(int(createdUser.ID))
	if err == nil {
		t.Fatal("Expected an error but got none")
	}
}

func TestUserRepository_Update(t *testing.T) {
	createdUser, err := hashPassAndGenerateUserInDb(&db.User{
		Username: "Jabar",
		Email:    "update@ymail.com",
		Password: "password",
		Name:     "Bamba",
	}, t)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	createdUser.Username = "Al-Amal"

	updatedUser, err := testConnection.repo.Update(int(createdUser.ID), createdUser)
	if err != nil {
		t.Fatalf("An error was encountered while updating: %v", err)
	}

	foundUser, err := testConnection.repo.FindById(int(updatedUser.ID))
	if err != nil {
		t.Errorf("An error was encountered while finding updated user: %v", err)
	}

	assert.Equal(t, foundUser, updatedUser, "Found user did is not equal to updated user")

	// Clean up: Delete created user
	testConnection.dbClient.Delete(updatedUser)
}

func TestUserRepository_FindAll(t *testing.T) {
	createdUser1, err := hashPassAndGenerateUserInDb(&db.User{
		Username: "Joe",
		Email:    "crazy@gmail.com",
		Password: "password",
		Name:     "Bamba",
	}, t)
	if err != nil {
		t.Fatalf("failed to create test user1: %v", err)
	}

	createdUser2, err := hashPassAndGenerateUserInDb(&db.User{
		Username: "Jabar",
		Email:    "scuba@ymail.com",
		Password: "password",
		Name:     "Bamba",
	}, t)
	if err != nil {
		t.Fatalf("failed to create test user2: %v", err)
	}

	users, err := testConnection.repo.FindAll(10, 0, "")
	if err != nil {
		t.Fatalf("failed to find all: %v", err)
	}
	// Make sure both users are in database
	if len(*users) != 2 {
		t.Errorf("Length of []users is not as expected. Got: %v", len(*users))
	}

	// t.Fatal(createdUser1.ID)
	for _, u := range *users {
		// If it's the first user
		if int(u.ID) == int(createdUser1.ID) {
			// check details of first created user
			if createdUser1.Email != u.Email {
				t.Errorf("Email of user1 doesn't match. Got: %v, expected %v", u.Email, createdUser1.Email)
			}
			if createdUser1.Username != u.Username {
				t.Errorf("Email of user1 doesn't match. Got: %v, expected %v", u.Username, createdUser1.Username)
			}
			if createdUser1.Name != u.Name {
				t.Errorf("Email of user1 doesn't match. Got: %v, expected %v", u.Name, createdUser1.Name)
			}
		} else {
			// check details of second user
			if createdUser2.Email != u.Email {
				t.Errorf("Email of user1 doesn't match. Got: %v, expected %v", u.Email, createdUser2.Email)
			}
			if createdUser2.Username != u.Username {
				t.Errorf("Email of user1 doesn't match. Got: %v, expected %v", u.Username, createdUser2.Username)
			}
			if createdUser2.Name != u.Name {
				t.Errorf("Email of user1 doesn't match. Got: %v, expected %v", u.Name, createdUser2.Name)
			}

		}
	}
	// Clean up created users
	usersToDelete := []db.User{{ID: createdUser1.ID}, {ID: createdUser2.ID}}
	testConnection.dbClient.Delete(usersToDelete)
}

// Test helper function
func hashPassAndGenerateUserInDb(user *db.User, t *testing.T) (*db.User, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		t.Fatalf("Couldn't create user")
	}
	user.Password = string(hashedPass)
	return testConnection.repo.Create(user)
}
