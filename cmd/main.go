package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dmawardi/Go-Template/ent"
	"github.com/dmawardi/Go-Template/ent/migrate"
	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/handlers"
	"github.com/dmawardi/Go-Template/internal/services"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const portNumber = ":8080"

// Init state
var app config.AppConfig

var store *sessions.CookieStore

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

	// Create sessions storage for authentication
	secretKey := os.Getenv("SESSIONS_SECRET_KEY")
	// Create new Session storage
	store = sessions.NewCookieStore([]byte(secretKey))
	// Set session in state to store
	app.Session = store

	// Initialize RBAC Authorization
	// e, _ := casbin.NewEnforcer("../internal/auth/model.conf", "path/to/policy.csv")
	// Set in state
	// app.RBEnforcer = e

	// Set state in other packages
	handlers.SetStateInHandlers(&app)
	auth.SetStateInAuth(&app)
	services.BuildServiceRepo(&app)

	// Create client using DbConnect
	client := DbConnect()
	app.DbClient = client
	// close the client once not operational
	defer client.Close()

	// _, err = handlers.CreateUser(ctx, client)
	// if err != nil {
	// 	fmt.Println("Unable to create user.")
	// }

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
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")

	// Create Postgres connection client
	client, err := ent.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, "TestGo", DB_PASS))

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
