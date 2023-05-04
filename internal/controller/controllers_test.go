package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/routes"

	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/controller"
	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/repository"
	"github.com/dmawardi/Go-Template/internal/service"
	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var testConnection TestDbRepo

var app config.AppConfig

type TestDbRepo struct {
	dbClient *gorm.DB
	users    userDB
	router   http.Handler
	// For authentication mocking
	accounts userAccounts
}

// DB structures
type userDB struct {
	repo repository.UserRepository
	serv service.UserService
	cont controller.UserController
}

// Account structures
type userAccounts struct {
	admin dummyAccount
	user  dummyAccount
}
type dummyAccount struct {
	details *db.User
	token   string
}

// Initial setup before running e2e tests in controllers_test package
func TestMain(m *testing.M) {
	// Setup a new/reset connection
	testConnection.setupDBAuthAppModels()

	// Setup accounts for mocking authentication
	testConnection.setupDummyAccounts(&db.User{
		Username: "Jabar",
		Email:    "Jabal@ymail.com",
		Password: "password",
		Name:     "Bamba",
	}, &db.User{
		Username: "Jabar",
		Email:    "Juba@ymail.com",
		Password: "password",
		Name:     "Bamba",
	})

	// Furnish the database with assets for testing
	testConnection.furnishDB()

	// build API for serving requests
	testConnection.router = testConnection.buildAPI()

	// Run the rest of the tests
	exitCode := m.Run()
	// exit with the same exit code as the tests
	os.Exit(exitCode)

}

// Builds new API using routes package
func (t TestDbRepo) buildAPI() http.Handler {
	api := routes.NewApi(
		t.users.cont,
	)
	// Extract handlers from api
	handler := api.Routes()

	return handler
}

// Setup functions
//
// Setup dummy admin and user account and apply to test connection
func (t *TestDbRepo) setupDummyAccounts(adminUser *db.User, basicUser *db.User) {
	// Build admin user
	adminUser, adminToken := t.generateUserWithRoleAndToken(
		adminUser, "admin")
	// Store credentials
	t.accounts.admin.details = adminUser
	t.accounts.admin.token = adminToken

	// Build normal user
	normalUser, userToken := t.generateUserWithRoleAndToken(
		basicUser, "user")
	// Store credentials
	t.accounts.user.details = normalUser
	t.accounts.user.token = userToken
}

// Setup Database, repos, services, controllers, dummy accounts for auth, and auth enforcer
func (t *TestDbRepo) setupDBAuthAppModels() {
	// Setup DB
	t.dbClient = setupDatabase()
	// Create test modules
	// Users
	t.users.repo = repository.NewUserRepository(t.dbClient)
	t.users.serv = service.NewUserService(t.users.repo)
	t.users.cont = controller.NewUserController(t.users.serv)

	// Setup the enforcer for usage as middleware
	setupTestEnforcer(t.dbClient)
}

// Fill this with any data required for tests
func (t *TestDbRepo) furnishDB() {

}

// Setup database connection
func setupDatabase() *gorm.DB {
	// Open a new, temporary database for testing
	dbClient, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		fmt.Errorf("failed to open database: %v", err)
	}

	// Migrate the database schema
	if err := dbClient.AutoMigrate(&db.User{}); err != nil {
		fmt.Errorf("failed to migrate database schema: %v", err)
	}

	return dbClient
}

// Setup enforcer and sync app state
func setupTestEnforcer(dbClient *gorm.DB) {
	// Build enforcer
	enforcer, err := auth.EnforcerSetup(dbClient)
	if err != nil {
		fmt.Println("Error building enforcer")
	}

	// Assign values in app for authentication
	app.DbClient = dbClient
	app.RBEnforcer = enforcer
	// Sync app in authentication package for usage in authentication functions
	auth.SetStateInAuth(&app)
}

// Helper functions
//
// Generates a new user, changes its role to admin and returns it with token
func (t TestDbRepo) generateUserWithRoleAndToken(user *db.User, role string) (*db.User, string) {
	unhashedPass := user.Password
	createdUser, err := t.hashPassAndGenerateUserInDb(user)
	if err != nil {
		fmt.Print("Problem creating admin user for tests.")
	}
	// Update user to admin
	createdUser.Role = role
	updatedUser, err := t.users.repo.Update(int(createdUser.ID), createdUser)
	// If match found (no errors)
	if err == nil {
		fmt.Println("Generating token for: ", updatedUser.Email)
		// Set login status to true
		tokenString, err := auth.GenerateJWT(int(updatedUser.ID), updatedUser.Email, updatedUser.Role)
		if err != nil {
			fmt.Println("Failed to create JWT")
		}

		// Add unhashed password to returned object
		updatedUser.Password = unhashedPass
		// Send to user in body
		return updatedUser, tokenString
	}
	return nil, ""
}

// Test helper function: Hashes password and generates a new user in the database
func (t TestDbRepo) hashPassAndGenerateUserInDb(user *db.User) (*db.User, error) {
	// Hash password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		fmt.Print("Couldn't hash password")
	}
	user.Password = string(hashedPass)

	// Create user
	createResult := t.dbClient.Create(user)
	if createResult.Error != nil {
		fmt.Printf("Couldn't create user: %v", user.Email)
	}

	return user, nil
}

// Build a struct object to a type of bytes.reader to fulfill io.reader interface
func buildReqBody(data interface{}) *bytes.Reader {
	// Marshal to JSON
	marshalled, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Failed to marshal JSON")
	}
	// Make into reader
	readerReqBody := bytes.NewReader(marshalled)
	return readerReqBody
}
