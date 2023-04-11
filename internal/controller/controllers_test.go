package controller_test

// import (
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"

// 	"github.com/dmawardi/Go-Template/internal/controller"
// 	"github.com/dmawardi/Go-Template/internal/models"
// 	"github.com/go-chi/chi"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// type MockDB struct{}

// func TestGetMyUserDetails(t *testing.T) {
// 	var testTable = []struct {
// 		name             string
// 		expectedResponse interface{}
// 		expectedStatus   int
// 	}{
// 		{"my-details-not-logged-in", &models.PartialUser{}, http.StatusBadRequest},
// 	}

// 	// Create a new request
// 	req := httptest.NewRequest("GET", "/api/me", nil)
// 	// Set a dummy token in the request header
// 	req.Header.Set("Authorization", "Bearer dummy_token")

// 	// Create a new response recorder
// 	rRec := httptest.NewRecorder()

// 	// Create a new router
// 	r := chi.NewRouter()
// 	// Add the handler to the router

// 	r.Get("/api/me", controller.GetMyUserDetails)

// 	for _, tt := range testTable {
// 		// Call the handler and use recorder
// 		r.ServeHTTP(rRec, req)

// 		// Check the status code is as expected
// 		if status := rRec.Code; status != tt.expectedStatus {
// 			t.Errorf("handler returned wrong status code: got %v want %v",
// 				status, tt.expectedStatus)
// 		}
// 		//
// 		if rRec.Body.String() != tt.expectedResponse {
// 			t.Errorf("handler returned unexpected body: got %v want %v",
// 				rRec.Body.String(), tt.expectedResponse)
// 		}

// 	}

// }

// func DbConnectForTesting(t *testing.T) *gorm.DB {
// 	// Grab environment variables for connection
// 	var DB_USER string = os.Getenv("DB_USER")
// 	var DB_PASS string = os.Getenv("DB_PASS")
// 	var DB_HOST string = os.Getenv("DB_HOST")
// 	var DB_PORT string = os.Getenv("DB_PORT")
// 	var DB_NAME string = os.Getenv("DB_NAME")

// 	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", DB_HOST, DB_USER, DB_PASS, DB_NAME, DB_PORT)

// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		t.Fatal("failed to connect database")
// 	}

// 	// // Migrate the schema
// 	// db.AutoMigrate(&User{})

// 	return db
// }
