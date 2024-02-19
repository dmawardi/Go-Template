package repository_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/repository"
	"gorm.io/gorm"
)

var testModule repositoryTestModule

var app config.AppConfig

type repositoryTestModule struct {
	dbClient *gorm.DB
	users    userModule
	auth     authModule
	posts    postModule
}

// Module structures
type userModule struct {
	repo repository.UserRepository
}
type authModule struct {
	repo repository.AuthPolicyRepository
}

type postModule struct {
	repo repository.PostRepository
}

// Initial setup before running tests in package
func TestMain(m *testing.M) {

	fmt.Printf("Setting up test connection\n")
	// Set URL in app state
	app.BaseURL = helpers.BuildBaseUrl()

	// Setup DB
	testModule.dbClient = helpers.SetupTestDatabase()
	// Set Gorm client
	app.DbClient = testModule.dbClient

	// Build enforcer
	enforcer, err := auth.EnforcerSetup(testModule.dbClient, false)
	if err != nil {
		fmt.Println("Error building enforcer")
	}
	// Set enforcer in state
	app.Auth.Enforcer = enforcer.Enforcer
	app.Auth.Adapter = enforcer.Adapter

	// Set app config in repository
	repository.SetAppConfig(&app)

	testModule.TestRepoSetup(testModule.dbClient)

	// Run the tests
	exitCode := m.Run()
	// exit with the same exit code as the tests
	os.Exit(exitCode)
}

// Builds new API using routes package
func (t *repositoryTestModule) TestRepoSetup(client *gorm.DB) {
	// Setup module stack
	// Auth
	t.auth.repo = repository.NewAuthPolicyRepository(client)
	// Users
	t.users.repo = repository.NewUserRepository(client)
	// Posts
	t.posts.repo = repository.NewPostRepository(client)
}
