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
	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/repository"
	"github.com/dmawardi/Go-Template/internal/routes"
	"github.com/dmawardi/Go-Template/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Set default port number
const portNumber = ":8080"

// Init state
var app config.AppConfig

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
	controller.SetStateInHandlers(&app)
	auth.SetStateInAuth(&app)
	adminpanel.SetStateInAdminPanel(&app)
	service.SetAppConfig(&app)
	repository.SetAppConfig(&app)
	routes.BuildRouteState(&app)

	// Create api
	api := ApiSetup(client)

	fmt.Printf("Starting application on port: %s\n", portNumber)

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
func ApiSetup(client *gorm.DB) routes.Api {
	mail := &helpers.EmailMock{}
	// Authorization
	groupRepo := repository.NewAuthPolicyRepository(client)
	groupService := service.NewAuthPolicyService(groupRepo)
	groupController := controller.NewAuthPolicyController(groupService)
	// user
	userRepo := repository.NewUserRepository(client)
	userService := service.NewUserService(userRepo, groupRepo, mail)
	userController := controller.NewUserController(userService)
	// post
	postRepo := repository.NewPostRepository(client)
	postService := service.NewPostService(postRepo)
	postController := controller.NewPostController(postService)

	// Admin panel
	selectorService := adminpanel.NewSelectorService(client, groupService)
	adminController := adminpanel.NewAdminController(
		adminpanel.NewAdminBaseController(userService),
		adminpanel.NewAdminUserController(userService, selectorService),
		adminpanel.NewAdminPostController(postService, selectorService),
		adminpanel.NewAdminAuthPolicyController(groupService, selectorService))

	// Generate admin sidebar list from admin controller
	adminpanel.GenerateAndSetAdminSidebar(adminController)

	// Build API using controllers
	api := routes.NewApi(adminController, userController, groupController, postController)
	return api
}
