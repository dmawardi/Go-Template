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
	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/email"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/modules"
	"github.com/dmawardi/Go-Template/internal/queue"
	"github.com/dmawardi/Go-Template/internal/repository"
	corerepositories "github.com/dmawardi/Go-Template/internal/repository/core"
	"github.com/dmawardi/Go-Template/internal/routes"
	"github.com/dmawardi/Go-Template/internal/seed"
	"github.com/dmawardi/Go-Template/internal/service"
	coreservices "github.com/dmawardi/Go-Template/internal/service/core"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Init state
var app config.AppConfig

// Connect email service used in API setup to connect email services (forgot password, etc.)
var connectEmailService = true

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
		Handler: api.Routes(),
	}

	// Listen and serve using server settings above
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// Edit this to use the entire appconfig instead of just the client
// Build API and store the services and repos in the config
func ApiSetup(client *gorm.DB, connectEmail bool) routes.Api {
	var mail email.Email
	// If connectEmail is true, use SMTP email
	if connectEmail {
		mail = email.NewSMTPEmail()
	} else {
		// Else, use email mock
		mail = &helpers.EmailMock{}
	}

	// Create job queue
	jobQueue := queue.NewQueue(client, mail)

	// Establish async job processing
	go jobQueue.Worker()

	// Authorization
	groupRepo := corerepositories.NewAuthPolicyRepository(client)
	groupService := coreservices.NewAuthPolicyService(groupRepo)
	groupController := core.NewAuthPolicyController(groupService)
	// user
	userRepo := corerepositories.NewUserRepository(client)
	userService := coreservices.NewUserService(userRepo, groupRepo, mail, jobQueue)
	userController := core.NewUserController(userService)

	// Build selector service is used for selector boxes in Admin panel
	selectorService := adminpanel.NewSelectorService(client, groupService)

	// Setup basic modules with new implementation (including admin controllers if available)
	moduleMap := modules.SetupModules(modules.ModulesToSetup, client, selectorService)

	// Admin panel
	//
	// Create admin controller
	adminController := adminpanel.NewAdminController(
		// Basic ADMIN modules
		adminpanel.NewAdminBaseController(userService),
		adminpanel.NewAdminUserController(userService, selectorService),
		adminpanel.NewAdminAuthPolicyController(groupService, selectorService),
		// ADD ADDITIONAL MODULES HERE
		moduleMap,
	)

	// Generate admin sidebar list from admin controller
	adminpanel.GenerateAndSetAdminSidebar(adminController)

	// Build API using controllers
	api := routes.NewApi(adminController, userController, groupController,
		// Created modules contained in moduleMap
		moduleMap,
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

// Slice of state setters for each module
var stateFuncs = []StateFunc{
	controller.SetStateInHandlers,
	auth.SetStateInAuth,
	adminpanel.SetStateInAdminPanel,
	service.SetAppConfig,
	repository.SetAppConfig,
	routes.BuildRouteState,
}
