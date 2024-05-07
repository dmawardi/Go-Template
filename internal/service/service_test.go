package service_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/helpers"
	webapi "github.com/dmawardi/Go-Template/internal/helpers/webApi"
	"github.com/dmawardi/Go-Template/internal/repository"
	"github.com/dmawardi/Go-Template/internal/service"
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
	serv service.UserService
}
type authModule struct {
	repo repository.AuthPolicyRepository
	serv service.AuthPolicyService
}

type postModule struct {
	repo repository.PostRepository
	serv service.PostService
}

// Initial setup before running tests in package
func TestMain(m *testing.M) {

	fmt.Printf("Setting up test connection\n")
	// Set URL in app state
	app.BaseURL = webapi.BuildBaseUrl()

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
	service.SetAppConfig(&app)
	auth.SetStateInAuth(&app)

	testModule.TestServSetup(testModule.dbClient)

	// Set up roles
	success, err := app.Auth.Enforcer.AddGroupingPolicies([][]string{{"role:admin", "role:moderator"}, {"role:moderator", "role:user"}})
	if !success {
		fmt.Println("Error setting up roles")
	}
	if err != nil {
		fmt.Println("Error setting up roles")
	}

	// Run the tests
	exitCode := m.Run()
	// exit with the same exit code as the tests
	os.Exit(exitCode)
}

// Builds new API using routes package
func (t *repositoryTestModule) TestServSetup(client *gorm.DB) {
	mail := &helpers.EmailMock{}
	// Setup module stack
	// Auth
	t.auth.repo = repository.NewAuthPolicyRepository(client)
	t.auth.serv = service.NewAuthPolicyService(t.auth.repo)
	// Users
	t.users.repo = repository.NewUserRepository(client)
	t.users.serv = service.NewUserService(t.users.repo, t.auth.repo, mail)
	// Posts
	t.posts.repo = repository.NewPostRepository(client)
	t.posts.serv = service.NewPostService(t.posts.repo)

}
