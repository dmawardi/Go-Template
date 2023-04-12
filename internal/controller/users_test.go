package controller_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/controller"
	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/go-chi/chi"

	"github.com/dmawardi/Go-Template/internal/repository"
	"github.com/dmawardi/Go-Template/internal/service"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type testDbRepo struct {
	dbClient *gorm.DB
	repo     repository.UserRepository
	serv     service.UserService
	cont     controller.UserController
	router   *chi.Mux
}

var testConnection testDbRepo

func init() {
	testConnection.dbClient = setupDatabase()
	// Create test modules
	testConnection.repo = repository.NewUserRepository(testConnection.dbClient)
	testConnection.serv = service.NewUserService(testConnection.repo)
	testConnection.cont = controller.NewUserController(testConnection.serv)
	// Create router
	testConnection.router = buildRouter(testConnection.cont)
}

func buildRouter(c controller.UserController) *chi.Mux {
	// Create a new chi router
	r := chi.NewRouter()

	// Use Authenticator
	r.Use(auth.AuthenticateJWT)
	// Basic user paths
	r.Get("/api/users", c.FindAll)
	r.Get("/api/users/{id}", c.Find)
	r.Put("/api/users/{id}", c.Update)
	r.Delete("/api/users/{id}", c.Delete)

	// My profile
	r.Get("/api/me", c.GetMyUserDetails)
	r.Put("/api/me", c.UpdateMyProfile)

	return r
}

func setupDatabase() *gorm.DB {
	// Open a new, temporary database for testing
	dbClient, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		fmt.Errorf("failed to open database: %v", err)
	}

	// Migrate the database schema
	if err := dbClient.AutoMigrate(&db.User{}); err != nil {
		fmt.Errorf("failed to migrate database schema: %v", err)
	}

	return dbClient
}

func TestUserController_Find(t *testing.T) {
	// Build test user
	userToCreate := &db.User{
		Username: "Jabar",
		Email:    "juba@ymail.com",
		Password: "password",
		Name:     "Bamba",
	}

	// Create user
	createdUser, err := hashPassAndGenerateUserInDb(userToCreate, t)
	if err != nil {
		t.Fatalf("failed to create test user for find by id user service testr: %v", err)
	}
	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/users/%v", createdUser.ID)
	// t.Fatalf("for url: %v\n. Created user iD: %v\n", requestUrl, createdUser.ID)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		t.Fatal(err)
	}

	foundUser, _ := testConnection.serv.FindById(int(createdUser.ID))
	fmt.Printf("for url: %v\n. Service result: %v \n", req.URL, foundUser)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Use handler with recorder and created request
	testConnection.router.ServeHTTP(rr, req)

	fmt.Printf("resp body: %v for url: %v\n. Service result: %v \n", rr.Body, requestUrl, foundUser)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// Test helper function: Hashes password and generates a new user in the database
func hashPassAndGenerateUserInDb(user *db.User, t *testing.T) (*db.User, error) {
	// Hash password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		t.Fatalf("Couldn't hash password")
	}
	user.Password = string(hashedPass)

	// Create user
	createResult := testConnection.dbClient.Create(user)
	if createResult.Error != nil {
		t.Fatalf("Couldn't create user: %v", user.Email)
	}

	return user, nil
}
