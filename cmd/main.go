package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"

	_ "github.com/swaggo/http-swagger/example/go-chi/docs"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/controller"
	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/repository"
	"github.com/dmawardi/Go-Template/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

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

	// Set state in other packages
	controller.SetStateInHandlers(&app)
	auth.SetStateInAuth(&app)
	service.BuildServiceState(&app)

	// Create client using DbConnect
	client := db.DbConnect()
	// Set in state
	app.DbClient = client

	// userToCreate := db.User{Name: "Goba", Username: "Walow", Password: "certainly", Email: "gustav@mail.com"}
	// createdUser, err := services.CreateUser(&userToCreate)

	// Create api
	api := ApiSetup(client)

	// Setup enforcer
	e, err := EnforcerSetup(client)
	if err != nil {
		log.Fatal("Couldn't setup RBAC Authorization Enforcer")
	}
	// Set enforcer in state
	app.RBEnforcer = e

	fmt.Printf("Starting application on port: %s\n", portNumber)

	// Server settings
	srv := &http.Server{
		Addr:    portNumber,
		Handler: api.routes(),
	}

	// Listen and serve using server settings above
	err = srv.ListenAndServe()
	if err != nil {

		log.Fatal(err)
	}
}

// Setup RBAC enforcer based using gorm client. Connects to DB and builds base policy
func EnforcerSetup(db *gorm.DB) (*casbin.Enforcer, error) {
	// Grab environment variables for connection
	var DB_PORT string = os.Getenv("DB_PORT")

	// Build adapter
	adapter, err := gormadapter.NewAdapterByDB(db)
	// If error
	if err != nil {
		log.Fatal("Couldn't build adapter for enforcer: ", err, "\nDB PORT", DB_PORT)
		return nil, err
	}

	// Initialize RBAC Authorization
	enforcer, err := casbin.NewEnforcer("./internal/auth/rbac_model.conf", adapter)

	// If error
	if err != nil {
		log.Fatal("Couldn't build RBAC enforcer: ", err)
		return nil, err
	}

	// Create default policies if not already detected within system
	auth.SetupCasbinPolicy(enforcer, auth.DefaultPolicyList)

	// else
	return enforcer, nil
}

func ApiSetup(client *gorm.DB) Api {
	// user
	userRepo := repository.NewUserRepository(client)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	// Build API using controllers
	api := NewApi(userController)
	return api
}
