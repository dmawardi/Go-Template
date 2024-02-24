package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	adminpanel "github.com/dmawardi/Go-Template/internal/admin-panel"
	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/routes"

	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/controller"
	"github.com/dmawardi/Go-Template/internal/repository"
	"github.com/dmawardi/Go-Template/internal/service"
	"gorm.io/gorm"
)

var testModule controllerTestModule

var app config.AppConfig

type controllerTestModule struct {
	dbClient *gorm.DB
	users    userModule
	admin    adminpanel.AdminController
	auth     authModule
	posts    postModule
	router   http.Handler
	api      routes.Api
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
	fmt.Printf("Setting up test connection\n")
	// Set URL in app state
	app.BaseURL = helpers.BuildBaseUrl()

	// Setup DB
	testModule.dbClient = helpers.SetupTestDatabase()
	// Set Gorm client
	app.DbClient = testModule.dbClient

	// Build enforcer
	enforcer, err := auth.EnforcerSetup(testModule.dbClient, true)
	if err != nil {
		fmt.Println("Error building enforcer")
	}
	// Set enforcer in state
	app.Auth.Enforcer = enforcer.Enforcer
	app.Auth.Adapter = enforcer.Adapter

	// Sync app in authentication package for usage in authentication functions
	SetAppWideState(app)

	// build API for serving requests
	testModule.api = testModule.TestApiSetup(testModule.dbClient)
	testModule.router = testModule.api.Routes()

	// Setup accounts for mocking authentication
	testModule.setupDummyAccounts(&models.CreateUser{
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

	// Run the rest of the tests
	exitCode := m.Run()
	// exit with the same exit code as the tests
	os.Exit(exitCode)
}

// Builds new API using routes package
func (t *controllerTestModule) TestApiSetup(client *gorm.DB) routes.Api {
	mail := &helpers.EmailMock{}
	// Setup module stack
	// Auth
	t.auth.repo = repository.NewAuthPolicyRepository(client)
	t.auth.serv = service.NewAuthPolicyService(t.auth.repo)
	t.auth.cont = controller.NewAuthPolicyController(t.auth.serv)
	// Users
	t.users.repo = repository.NewUserRepository(client)
	t.users.serv = service.NewUserService(t.users.repo, t.auth.repo, mail)
	t.users.cont = controller.NewUserController(t.users.serv)
	// Posts
	t.posts.repo = repository.NewPostRepository(client)
	t.posts.serv = service.NewPostService(t.posts.repo)
	t.posts.cont = controller.NewPostController(t.posts.serv)

	// Admin panel
	selectorService := adminpanel.NewSelectorService(client, t.auth.serv)
	t.admin = adminpanel.NewAdminController(
		adminpanel.NewAdminBaseController(t.users.serv),
		adminpanel.NewAdminUserController(t.users.serv, selectorService),
		adminpanel.NewAdminPostController(t.posts.serv, selectorService),
		adminpanel.NewAdminAuthPolicyController(t.auth.serv, selectorService))

	// Generate admin sidebar list from admin controller
	adminpanel.GenerateAndSetAdminSidebar(t.admin)

	// Setup API using controllers
	api := routes.NewApi(
		t.admin,
		t.users.cont,
		t.auth.cont,
		t.posts.cont,
	)

	return api
}

// Setup functions
//
// Setup dummy admin and user account and apply to test connection
func (t *controllerTestModule) setupDummyAccounts(adminUser *models.CreateUser, basicUser *models.CreateUser) {
	adminUser.Role = "role:admin"
	// Build admin user
	createdAdminUser, adminToken := t.generateUserWithRoleAndToken(
		adminUser)
	// Store credentials
	t.accounts.admin.details = createdAdminUser
	t.accounts.admin.token = adminToken

	basicUser.Role = "role:user"
	// Build normal user
	createdBasicUser, userToken := t.generateUserWithRoleAndToken(
		basicUser)
	// Store credentials
	t.accounts.user.details = createdBasicUser
	t.accounts.user.token = userToken
}

// Generates a new user, changes its role to admin and returns it with token
func (t *controllerTestModule) generateUserWithRoleAndToken(user *models.CreateUser) (*models.UserWithRole, string) {
	// Create user
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

// Helper functions
//
// Sets app config state to all packages for usage
func SetAppWideState(appConfig config.AppConfig) {
	controller.SetStateInHandlers(&appConfig)
	auth.SetStateInAuth(&appConfig)
	adminpanel.SetStateInAdminPanel(&appConfig)
	service.BuildServiceState(&appConfig)
	repository.SetAppConfig(&appConfig)
	routes.BuildRouteState(&appConfig)
}

// A helper function to build an API request that starts with url of '/api/'
func buildApiRequest(method string, urlSuffix string, body io.Reader, authHeaderRequired bool, token string) (request *http.Request, err error) {
	req, err := http.NewRequest(method, fmt.Sprintf("/api/%v", urlSuffix), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	// If authorization header required
	if authHeaderRequired {
		req.Header.Set("Authorization", fmt.Sprintf("bearer %v", token))
	}
	return req, nil
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
