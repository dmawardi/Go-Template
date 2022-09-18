package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/casbin/casbin/v2"
	entadapter "github.com/casbin/ent-adapter"
	"github.com/dmawardi/Go-Template/ent"
	"github.com/dmawardi/Go-Template/ent/migrate"
	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/handlers"
	"github.com/dmawardi/Go-Template/internal/services"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const portNumber = ":8080"

// Init state
var app config.AppConfig

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
	handlers.SetStateInHandlers(&app)
	auth.SetStateInAuth(&app)
	services.BuildServiceState(&app)

	// Create client using DbConnect
	client := DbConnect()
	app.DbClient = client
	// close the client once not operational
	defer client.Close()

	// Setup enforcer
	e, err := EnforcerSetup()
	if err != nil {
		log.Fatal("Couldn't setup RBAC Authorization Enforcer")
	}
	// Set enforcer in state
	app.RBEnforcer = e

	fmt.Printf("Starting application on port: %s\n", portNumber)

	// Server settings
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}

	// Listen and serve using server settings above
	err = srv.ListenAndServe()
	if err != nil {

		log.Fatal(err)
	}
}

// DbConnect connects to database using ent
func DbConnect() *ent.Client {
	// Grab environment variables for connection
	var DB_USER string = os.Getenv("DB_USER")
	var DB_PASS string = os.Getenv("DB_PASS")
	var DB_HOST string = os.Getenv("DB_HOST")
	var DB_PORT string = os.Getenv("DB_PORT")

	// Create Postgres connection client
	client, err := ent.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, "TestGo", DB_PASS))

	fmt.Println("DB PORT in db connect", DB_PORT)
	// Handle error
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	// If no error detected
	fmt.Println("Successfully connected to DB")
	// close client at end of function
	// defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background(), migrate.WithDropIndex(true),
		migrate.WithDropColumn(true)); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client
}

// Setup RBAC enforcer based on local model. Connects to DB and builds base policy
func EnforcerSetup() (*casbin.Enforcer, error) {
	// Grab environment variables for connection
	var DB_USER string = os.Getenv("DB_USER")
	var DB_PASS string = os.Getenv("DB_PASS")
	var DB_HOST string = os.Getenv("DB_HOST")
	var DB_PORT string = os.Getenv("DB_PORT")

	// Create new adapter
	enforcerAdapter, err := entadapter.NewAdapter("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, "TestGo", DB_PASS))
	// If error
	if err != nil {
		log.Fatal("Couldn't connect Enforcer to DB: ", err, "\nDB PORT", DB_PORT)
		return nil, err
	}

	// Initialize RBAC Authorization
	enforcer, err := casbin.NewEnforcer("./internal/auth/rbac_model.conf", enforcerAdapter)
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
