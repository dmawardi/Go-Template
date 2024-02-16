package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
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
	// Setup DB
	testConnection.dbClient = setupTestDatabase()

	// Build enforcer
	enforcer, err := auth.EnforcerSetup(testConnection.dbClient)
	if err != nil {
		fmt.Println("Error building enforcer")
	}

	// Set in app state
	app.BaseURL = syncBaseUrl()

	// Set Gorm client
	app.DbClient = testConnection.dbClient
	// Set enforcer in state
	app.Auth.Enforcer = enforcer.Enforcer
	app.Auth.Adapter = enforcer.Adapter

	// Sync app in authentication package for usage in authentication functions
	SetAppWideState(app)

	// build API for serving requests
	testConnection.api = testConnection.TestApiSetup(testConnection.dbClient)
	testConnection.router = testConnection.api.Routes()

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

	// Run the rest of the tests
	exitCode := m.Run()
	// exit with the same exit code as the tests
	os.Exit(exitCode)
}

// Builds new API using routes package
func (t *controllerTestModule) TestApiSetup(client *gorm.DB) routes.Api {
	// Setup module stack
	// Auth
	t.auth.repo = repository.NewAuthPolicyRepository(client)
	t.auth.serv = service.NewAuthPolicyService(t.auth.repo)
	t.auth.cont = controller.NewAuthPolicyController(t.auth.serv)
	// Users
	t.users.repo = repository.NewUserRepository(client)
	t.users.serv = service.NewUserService(t.users.repo, t.auth.repo)
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

func syncBaseUrl() string {
	// Extract environment variables
	serverUrl := os.Getenv("SERVER_BASE_URL")
	portNumber := os.Getenv("SERVER_PORT")

	// Get BASE URL from environment variables
	baseURL := fmt.Sprintf("%s%s", serverUrl, portNumber)
	return baseURL
}

// Sets app config state to all packages for usage
func SetAppWideState(appConfig config.AppConfig) {
	controller.SetStateInHandlers(&appConfig)
	auth.SetStateInAuth(&appConfig)
	adminpanel.SetStateInAdminPanel(&appConfig)
	service.BuildServiceState(&appConfig)
	repository.SetAppConfig(&appConfig)
	routes.BuildRouteState(&appConfig)
}

// Helper functions
//
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

// Setup database connection
func setupTestDatabase() *gorm.DB {
	// Open a new, temporary database for testing
	dbClient, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		fmt.Printf("failed to open database: %v", err)
	}

	// Migrate the database schema
	if err := dbClient.AutoMigrate(&db.User{}, &db.Post{}); err != nil {
		fmt.Printf("failed to migrate database schema: %v", err)
	}

	return dbClient
}

// buildApiRequest is a helper function to build an API request that starts with url of '/api/'
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

// CompareObjects compares the specified fields of two interface{} objects.
// It uses reflection to dynamically compare the field values of both objects.
func CompareObjects(actualObject interface{}, expectedObject interface{}, t *testing.T, fieldsToCheck []string) {
	// Convert both objects to reflect.Value to facilitate comparison.
	actualValue := reflect.ValueOf(actualObject)
	if actualValue.Kind() == reflect.Ptr {
		actualValue = actualValue.Elem()
	}
	expectedValue := reflect.ValueOf(expectedObject)
	if expectedValue.Kind() == reflect.Ptr {
		expectedValue = expectedValue.Elem()
	}

	// Iterate over the specified fields to compare their values.
	for _, field := range fieldsToCheck {
		actualFieldValue := actualValue.FieldByName(field)
		expectedFieldValue := expectedValue.FieldByName(field)

		// Check if both fields are valid.
		if !actualFieldValue.IsValid() {
			t.Errorf("actual object does not have field %s", field)
			continue
		}
		if !expectedFieldValue.IsValid() {
			t.Errorf("expected object does not have field %s", field)
			continue
		}

		// Compare the actual and expected field values.
		if !reflect.DeepEqual(actualFieldValue.Interface(), expectedFieldValue.Interface()) {
			t.Errorf("field %s does not match: expected %v, got %v", field, expectedFieldValue.Interface(), actualFieldValue.Interface())
		}
	}
}

// UpdateModelFields updates the fields of a GORM model based on a map[string]string.
// The model parameter is expected to be a pointer to a struct that's a GORM model.
// The updates parameter is a map where keys are field names and values are new values for those fields, as strings.
func UpdateModelFields(model interface{}, updates map[string]string) error {
	// Ensure the model is a pointer to a struct.
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() != reflect.Ptr || modelValue.Elem().Kind() != reflect.Struct {
		return errors.New("model must be a pointer to a struct")
	}

	// Get the underlying struct value.
	structValue := modelValue.Elem()

	// Iterate through the updates map to update struct fields.
	for field, newValue := range updates {
		// Find the struct field.
		structField := structValue.FieldByName(field)
		if !structField.IsValid() {
			return fmt.Errorf("no such field: %s in model", field)
		}

		// Ensure the field can be set.
		if !structField.CanSet() {
			return fmt.Errorf("cannot set field: %s", field)
		}

		// Convert and set the field value based on its kind.
		switch structField.Kind() {
		case reflect.String:
			structField.SetString(newValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.ParseInt(newValue, 10, 64)
			if err != nil {
				return fmt.Errorf("cannot convert %s to int for field: %s", newValue, field)
			}
			structField.SetInt(intVal)
		case reflect.Float32, reflect.Float64:
			floatVal, err := strconv.ParseFloat(newValue, 64)
			if err != nil {
				return fmt.Errorf("cannot convert %s to float for field: %s", newValue, field)
			}
			structField.SetFloat(floatVal)
		// Add more cases here for other types as needed.
		default:
			return fmt.Errorf("unsupported field type: %s for field: %s", structField.Type(), field)
		}
	}

	return nil
}
