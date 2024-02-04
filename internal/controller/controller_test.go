package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	adminpanel "github.com/dmawardi/Go-Template/internal/admin-panel"
	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/routes"

	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/controller"
	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/repository"
	"github.com/dmawardi/Go-Template/internal/service"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var testConnection controllerTestModule

var app config.AppConfig

type controllerTestModule struct {
	dbClient *gorm.DB
	users    userModule
	admin    adminpanel.AdminController
	auth     authModule
	posts    postModule
	router   http.Handler
	// For authentication mocking
	accounts userAccounts
}

// Module structures
type userModule struct {
	repo repository.UserRepository
	serv service.UserService
	cont controller.UserController
}
type authModule struct {
	repo repository.AuthPolicyRepository
	serv service.AuthPolicyService
	cont controller.AuthPolicyController
}

type postModule struct {
	repo repository.PostRepository
	serv service.PostService
	cont controller.PostController
}

// Account structures
type userAccounts struct {
	admin dummyAccount
	user  dummyAccount
}
type dummyAccount struct {
	details *models.UserWithRole
	token   string
}

// Initial setup before running e2e tests in controllers_test package
func TestMain(m *testing.M) {
	// Setup DB
	testConnection.dbClient = setupTestDatabase()

	// Build enforcer
	enforcer, err := auth.EnforcerSetup(testConnection.dbClient)
	if err != nil {
		fmt.Println("Error building enforcer")
	}

	// Set app state
	// Set Gorm client
	app.DbClient = testConnection.dbClient
	// Set enforcer in state
	app.Auth.Enforcer = enforcer.Enforcer
	app.Auth.Adapter = enforcer.Adapter
	// Sync app in authentication package for usage in authentication functions
	auth.SetStateInAuth(&app)

	// Setup accounts for mocking authentication
	testConnection.setupDummyAccounts(&models.CreateUser{
		Username: "Jabar",
		Email:    "Jabal@ymail.com",
		Password: "password",
		Name:     "Bamba",
	}, &models.CreateUser{
		Username: "Jabar",
		Email:    "Juba@ymail.com",
		Password: "password",
		Name:     "Bamba",
	})

	// build API for serving requests
	testConnection.router = testConnection.buildAPI()

	// Run the rest of the tests
	exitCode := m.Run()
	// exit with the same exit code as the tests
	os.Exit(exitCode)
}

// Builds new API using routes package
func (t controllerTestModule) buildAPI() http.Handler {
	// Setup module stack
	t.setupModuleStack()
	// Setup API using controllers
	api := routes.NewApi(
		t.admin,
		t.users.cont,
		t.auth.cont,
		t.posts.cont,
	)
	// Extract handlers from api
	handler := api.Routes()

	return handler
}

// Setup functions
//
// Setup dummy admin and user account and apply to test connection
func (t *controllerTestModule) setupDummyAccounts(adminUser *models.CreateUser, basicUser *models.CreateUser) {
	adminUser.Role = "admin"
	// Build admin user
	createdAdminUser, adminToken := t.generateUserWithRoleAndToken(
		adminUser)
	// Store credentials
	t.accounts.admin.details = createdAdminUser
	t.accounts.admin.token = adminToken

	// Build normal user
	normalUser, userToken := t.generateUserWithRoleAndToken(
		basicUser)
	// Store credentials
	t.accounts.user.details = normalUser
	t.accounts.user.token = userToken
}

// Setup Database, repos, services, controllers, dummy accounts for auth, and auth enforcer
func (t *controllerTestModule) setupModuleStack() {
	// Create test modules
	// Auth
	t.auth.repo = repository.NewAuthPolicyRepository(t.dbClient)
	t.auth.serv = service.NewAuthPolicyService(t.auth.repo)
	t.auth.cont = controller.NewAuthPolicyController(t.auth.serv)
	// Users
	t.users.repo = repository.NewUserRepository(t.dbClient)
	t.users.serv = service.NewUserService(t.users.repo, t.auth.repo)
	t.users.cont = controller.NewUserController(t.users.serv)
	// Posts
	t.posts.repo = repository.NewPostRepository(t.dbClient)
	t.posts.serv = service.NewPostService(t.posts.repo)
	t.posts.cont = controller.NewPostController(t.posts.serv)
}

// Helper functions
//
// Generates a new user, changes its role to admin and returns it with token
func (t controllerTestModule) generateUserWithRoleAndToken(user *models.CreateUser) (*models.UserWithRole, string) {
	createdUser, err := t.users.serv.Create(user)

	// If match found (no errors)
	if err == nil {
		fmt.Println("Generating token for: ", createdUser.Email)
		// Set login status to true
		tokenString, err := auth.GenerateJWT(int(createdUser.ID), createdUser.Email, createdUser.Role)
		if err != nil {
			fmt.Println("Failed to create JWT")
		}

		// Add unhashed password to returned object
		createdUser.Password = user.Password
		// Send to user in body
		return createdUser, tokenString
	}
	return nil, ""
}

// Setup database connection
func setupTestDatabase() *gorm.DB {
	// Open a new, temporary database for testing
	dbClient, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		fmt.Errorf("failed to open database: %v", err)
	}

	// Migrate the database schema
	if err := dbClient.AutoMigrate(&db.User{}, &db.Post{}); err != nil {
		fmt.Errorf("failed to migrate database schema: %v", err)
	}

	return dbClient
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
