package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/gorm"

	_ "github.com/swaggo/http-swagger/example/go-chi/docs"

	adminpanel "github.com/dmawardi/Go-Template/internal/admin-panel"
	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/controller"
	"github.com/dmawardi/Go-Template/internal/controller/core"
	modulecontrollers "github.com/dmawardi/Go-Template/internal/controller/moduleControllers"
	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/email"
	"github.com/dmawardi/Go-Template/internal/helpers"
	webapi "github.com/dmawardi/Go-Template/internal/helpers/webApi"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/repository"
	corerepositories "github.com/dmawardi/Go-Template/internal/repository/core"
	modulerepositories "github.com/dmawardi/Go-Template/internal/repository/module"
	"github.com/dmawardi/Go-Template/internal/routes"
	"github.com/dmawardi/Go-Template/internal/seed"
	"github.com/dmawardi/Go-Template/internal/service"
	coreservices "github.com/dmawardi/Go-Template/internal/service/core"
	moduleservices "github.com/dmawardi/Go-Template/internal/service/module"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Init state
var app config.AppConfig

// Connect email service used in API setup to connect email services (forgot password, etc.)
var connectEmailService = true

// Slice of state setters for each module
var stateFuncs = []StateFunc{
	controller.SetStateInHandlers,
	auth.SetStateInAuth,
	adminpanel.SetStateInAdminPanel,
	service.SetAppConfig,
	repository.SetAppConfig,
	routes.BuildRouteState,
}

// Define setup configurations (to use in setupModules within API setup function)
var modulesToSetup = []models.EntityConfig{
	{
		Name:          "Post",
		NewRepo:       webapi.NewRepository(modulerepositories.NewPostRepository),
		NewService:    webapi.NewService(moduleservices.NewPostService),
		NewController: webapi.NewController(modulecontrollers.NewPostController),
	},
	// ADD ADDITIONAL BASIC MODULES HERE
}

// API Details
// @title           Go Template
// @version         1.0
// @description     This is a template API server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/

// @securityDefinitions.apikey BearerToken
// @in header
// @name Authorization

func main() {
	// Build context
	ctx := context.Background()
	// Set context in app config
	app.Ctx = ctx
	// Load env variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Unable to load environment variables.")
	}
	// Extract environment variables
	serverUrl := os.Getenv("SERVER_BASE_URL")
	portNumber := os.Getenv("SERVER_PORT")
	// If port number is empty, set to default
	if portNumber == "" {
		portNumber = ":8080"
	}

	// Get BASE URL from environment variables
	baseURL := fmt.Sprintf("%s%s", serverUrl, portNumber)
	// Set in app state
	app.BaseURL = baseURL

	// Parse the template files in the templates directory
	tmpl, err := adminpanel.ParseAdminTemplates()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Set template in state
	app.AdminTemplates = tmpl

	// Create client using DbConnect
	client := db.DbConnect()
	// Set in state
	app.DbClient = client

	// Setup enforcer
	e, err := auth.EnforcerSetup(client, true)
	if err != nil {
		log.Fatal("Couldn't setup RBAC Authorization Enforcer")
	}
	// Set enforcer in state
	app.Auth.Enforcer = e.Enforcer
	app.Auth.Adapter = e.Adapter

	// Set state in other packages
	setAppState(&app, stateFuncs)

	// Seed the database
	err = seed.Boot(client)
	if err != nil {
		log.Fatal(err)
	}

	// Create api
	api := ApiSetup(client, connectEmailService)

	fmt.Printf("Starting application: http://%s%s\n", serverUrl, portNumber)

	// Server settings
	srv := &http.Server{
		Addr:    portNumber,
		Handler: api.Routes(routes.CRUDRouteSet{}, routes.AdminRouteSet{}),
	}

	// Listen and serve using server settings above
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// Edit this to use the entire appconfig instead of just the client
// Build API and store the services and repos in the config
func ApiSetup(client *gorm.DB, emailMock bool) routes.Api {
	var mail email.Email
	// If emailMock is false, use SMTP email
	if !emailMock {
		mail = email.NewSMTPEmail()
	} else {
		// Else, use email mock
		mail = &helpers.EmailMock{}
	}

	// Authorization
	groupRepo := corerepositories.NewAuthPolicyRepository(client)
	groupService := coreservices.NewAuthPolicyService(groupRepo)
	groupController := core.NewAuthPolicyController(groupService)
	// user
	userRepo := corerepositories.NewUserRepository(client)
	userService := coreservices.NewUserService(userRepo, groupRepo, mail)
	userController := core.NewUserController(userService)

	// Setup basic modules with new implementation
	moduleMap := webapi.SetupModules(modulesToSetup, client)

	// Admin panel
	selectorService := adminpanel.NewSelectorService(client, groupService)
	adminController := adminpanel.NewAdminController(
		adminpanel.NewAdminBaseController(userService),
		adminpanel.NewAdminUserController(userService, selectorService),
		adminpanel.NewAdminAuthPolicyController(groupService, selectorService),
		// Basic modules
		adminpanel.NewAdminPostController(moduleMap["Post"].Service.(moduleservices.PostService), selectorService),
		// ADD ADDITIONAL BASIC MODULES HERE
	)

	// Generate admin sidebar list from admin controller
	adminpanel.GenerateAndSetAdminSidebar(adminController)

	// Build API using controllers
	api := routes.NewApi(adminController, userController, groupController,
		// ADD BASIC MODULES HERE
		moduleMap["Post"].Controller.(modulecontrollers.PostController),
	)
	return api
}

// STATE MANAGEMENT
//
// The state of the app is set using a series of state functions.
// StateFunc defines the type of function that can set state on App.
// ie. Any function that takes app as an argument and returns nothing.
type StateFunc func(*config.AppConfig)

// setAppState sets the state of the app using the provided StateFuncs.
func setAppState(app *config.AppConfig, funcs []StateFunc) {
	for _, fn := range funcs {
		fn(app)
	}
}
