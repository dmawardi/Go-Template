package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/controller"
	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
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
	// For authentication mocking
	admin      *db.User
	adminToken string
	user       *db.User
	userToken  string
}

var testConnection testDbRepo

var app config.AppConfig

func init() {
	testConnection.dbClient = setupDatabase()
	// Create test modules
	testConnection.repo = repository.NewUserRepository(testConnection.dbClient)
	testConnection.serv = service.NewUserService(testConnection.repo)
	testConnection.cont = controller.NewUserController(testConnection.serv)
	// Create router
	testConnection.router = buildRouter(testConnection.cont)

	// Build admin user
	adminUser, adminToken := generateUserWithRoleAndToken(
		&db.User{
			Username: "Jabar",
			Email:    "juba@ymail.com",
			Password: "password",
			Name:     "Bamba",
		}, "admin")
	// Store credentials
	testConnection.admin = adminUser
	testConnection.adminToken = adminToken

	// Build normal user
	normalUser, userToken := generateUserWithRoleAndToken(
		&db.User{
			Username: "Jabar",
			Email:    "Jabal@ymail.com",
			Password: "password",
			Name:     "Bamba",
		}, "user")
	// Store credentials
	testConnection.user = normalUser
	testConnection.userToken = userToken

	// Setup the enforcer for usage as middleware
	setupTestEnforcer(testConnection.dbClient)
}

func TestUserController_Find(t *testing.T) {
	// Build test user
	userToCreate := &db.User{
		Username: "Jabar",
		Email:    "greenie@ymail.com",
		Password: "password",
		Name:     "Bamba",
	}

	// Create user
	createdUser, err := hashPassAndGenerateUserInDb(userToCreate)
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
	// Add auth token to header
	req.Header.Set("Authorization", fmt.Sprintf("bearer %v", testConnection.adminToken))
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Use handler with recorder and created request
	testConnection.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Convert response JSON to struct
	var body db.User
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Verify that the found user matches the original user
	if body.ID != createdUser.ID {
		t.Errorf("found createdUser has incorrect ID: expected %d, got %d", userToCreate.ID, body.ID)
	}
	if body.Email != userToCreate.Email {
		t.Errorf("found createdUser has incorrect email: expected %s, got %s", userToCreate.Email, body.Email)
	}
	if body.Username != userToCreate.Username {
		t.Errorf("found createdUser has incorrect username: expected %s, got %s", userToCreate.Username, body.Username)
	}
	if body.Name != userToCreate.Name {
		t.Errorf("found createdUser has incorrect name: expected %s, got %s", userToCreate.Name, body.Name)
	}

	// Delete the created user
	testConnection.dbClient.Delete(createdUser)
}

func TestUserController_FindAll(t *testing.T) {
	// Create a new request
	req, err := http.NewRequest("GET", "/api/users?limit=10&offset=0&order=role DESC", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Add auth token to header
	req.Header.Set("Authorization", fmt.Sprintf("bearer %v", testConnection.adminToken))
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Use handler with recorder and created request
	testConnection.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Convert response JSON to struct
	var body []db.User
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Check length of user array
	if len(body) != 2 {
		t.Errorf("Users array in findAll failed: expected %d, got %d", 2, len(body))
	}

	// Iterate through users array received
	for _, item := range body {
		// If id is admin id
		if item.ID == testConnection.admin.ID {
			// Check details
			if item.Email != testConnection.admin.Email {
				t.Errorf("found createdUser has incorrect email: expected %s, got %s", testConnection.admin.Email, item.Email)
			}
			if item.Username != testConnection.admin.Username {
				t.Errorf("found createdUser has incorrect username: expected %s, got %s", testConnection.admin.Username, item.Username)
			}
			if item.Name != testConnection.admin.Name {
				t.Errorf("found createdUser has incorrect name: expected %s, got %s", testConnection.admin.Name, item.Name)
			}
		} else {
			// Else check user details
			if item.Email != testConnection.user.Email {
				t.Errorf("found createdUser has incorrect email: expected %s, got %s", testConnection.user.Email, item.Email)
			}
			if item.Username != testConnection.user.Username {
				t.Errorf("found createdUser has incorrect username: expected %s, got %s", testConnection.user.Username, item.Username)
			}
			if item.Name != testConnection.user.Name {
				t.Errorf("found createdUser has incorrect name: expected %s, got %s", testConnection.user.Name, item.Name)
			}
		}
	}
}

func TestUserController_Delete(t *testing.T) {
	// Build test user for deletion
	userToCreate := &db.User{
		Username: "Jabar",
		Email:    "greenie@ymail.com",
		Password: "password",
		Name:     "Bamba",
	}

	// Create user
	createdUser, err := hashPassAndGenerateUserInDb(userToCreate)
	if err != nil {
		t.Fatalf("failed to create test user for delete user controller test: %v", err)
	}
	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/users/%v", createdUser.ID)
	// t.Fatalf("for url: %v\n. Created user iD: %v\n", requestUrl, createdUser.ID)
	req, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Add user auth token to header
	req.Header.Set("Authorization", fmt.Sprintf("bearer %v", testConnection.userToken))
	// Send deletion requestion to mock server
	testConnection.router.ServeHTTP(rr, req)
	// Check response is failed for normal user
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("User deletion test: got %v want %v.",
			status, http.StatusForbidden)
	}

	// Check admin delete works
	// Reset http recorder
	rr = httptest.NewRecorder()

	// Set replacement header with admin credentials
	req.Header.Set("Authorization", fmt.Sprintf("bearer %v", testConnection.adminToken))
	// Perform GET request to mock server (using admin token)
	testConnection.router.ServeHTTP(rr, req)
	// Check the response status code for user deletion success
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("User deletion test: got %v want %v.",
			status, http.StatusOK)
	}
}

func TestUserController_Update(t *testing.T) {
	var updateTests = []struct {
		data                   map[string]string
		tokenToUse             string
		expectedResponseStatus int
		checkDetails           bool
	}{
		{map[string]string{
			"Username": "JabarHindi",
			"Name":     "Bambaloonie",
		}, testConnection.userToken, http.StatusForbidden, false},
		{map[string]string{
			"Username": "JabarHindi",
			"Name":     "Bambaloonie",
		}, testConnection.adminToken, http.StatusOK, true},
		// Update should be disallowed due to being too short
		{map[string]string{
			"Username": "Gobod",
			"Name":     "solu",
		}, testConnection.adminToken, http.StatusBadRequest, false},
		// User should be forbidden before validating
		{map[string]string{
			"Username": "Gobod",
			"Name":     "solu",
		}, testConnection.userToken, http.StatusForbidden, false},
	}

	// Build test user for update
	userToCreate := &db.User{
		Username: "Jabarnam",
		Email:    "sweenie@ymail.com",
		Password: "password",
		Name:     "Bambaliya",
	}
	// Create user
	createdUser, err := hashPassAndGenerateUserInDb(userToCreate)
	if err != nil {
		t.Fatalf("failed to create test user for delete user controller test: %v", err)
	}

	// Create a request url with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/users/%v", createdUser.ID)

	for _, v := range updateTests {
		// Make new request with user update in body
		req, err := http.NewRequest("PUT", requestUrl, buildReqBody(v.data))
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()
		// Add user auth token to header
		req.Header.Set("Authorization", fmt.Sprintf("bearer %v", v.tokenToUse))
		// Send update request to mock server
		testConnection.router.ServeHTTP(rr, req)
		// Check response is failed for normal user to update another
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("User update test: got %v want %v.",
				status, v.expectedResponseStatus)
		}

		// If need to check details
		if v.checkDetails == true {
			// Update created user struct with the changes pushed through API
			updateChangesOnly(createdUser, v.data)

			// Check user details using updated object
			checkUserDetails(rr, createdUser, t, true)
		}
	}

	// Delete the created user
	testConnection.dbClient.Delete(createdUser)
}

func TestUserController_Create(t *testing.T) {
	var updateTests = []struct {
		data                   models.CreateUser
		expectedResponseStatus int
	}{
		{models.CreateUser{
			Username: "Jabarnam",
			Email:    "gabor@ymail.com",
			Password: "password",
			Name:     "Bambaliya",
		}, http.StatusCreated},
		{models.CreateUser{
			Username: "Swalanim",
			Email:    "salvia@ymail.com",
			Password: "seradfasdf",
			Name:     "CreditTomyaA",
		}, http.StatusCreated},
		// Create should be disallowed due to not being email
		{models.CreateUser{
			Username: "Yukon",
			Email:    "Sylvio",
			Password: "wowogsdfg",
			Name:     "Sosawsdfgsdfg",
		}, http.StatusBadRequest},
		// Should be a bad request due to pass/name length
		{models.CreateUser{
			Username: "Jabarnam",
			Email:    "Cakawu@ymail.com",
			Password: "as",
			Name:     "df",
		}, http.StatusBadRequest},
		// Should be a bad request due to duplicate user (created in init)
		{models.CreateUser{
			Username: "Jabarnam",
			Email:    "Jabal@ymail.com",
			Password: "as",
			Name:     "df",
		}, http.StatusBadRequest},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := "/api/users"

	for _, v := range updateTests {
		// Make new request with user update in body
		req, err := http.NewRequest("POST", requestUrl, buildReqBody(v.data))
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Send update request to mock server
		testConnection.router.ServeHTTP(rr, req)
		// Check response is failed for normal user to update another
		if status := rr.Code; status != v.expectedResponseStatus {

			t.Errorf("User update test (%v): got %v want %v.", v.data.Name,
				status, v.expectedResponseStatus)
		}

		// Init body for response extraction
		var body db.User
		// Grab ID from response body
		json.Unmarshal(rr.Body.Bytes(), &body)

		// Delete the created user
		testConnection.serv.Delete(int(body.ID))
	}
}

func TestUserController_GetMyUserDetails(t *testing.T) {
	var updateTests = []struct {
		checkDetails           bool
		tokenToUse             string
		userToCheck            db.User
		expectedResponseStatus int
	}{
		{true, testConnection.userToken, *testConnection.user, http.StatusOK},
		{true, testConnection.adminToken, *testConnection.admin, http.StatusOK},
		// Deny access to user that doesn't have authentication
		{false, "", db.User{}, http.StatusForbidden},
	}
	// Create a request url with an "id" URL parameter
	requestUrl := "/api/me"

	for _, v := range updateTests {
		// Make new request with user update in body
		req, err := http.NewRequest("GET", requestUrl, nil)
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()

		// If you need to check details for successful requests, set token
		if v.checkDetails {
			// Add user auth token to header
			req.Header.Set("Authorization", fmt.Sprintf("bearer %v", v.tokenToUse))
		}
		// Send update request to mock server
		testConnection.router.ServeHTTP(rr, req)
		// Check response is failed for normal user to update another
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("User update test: got %v want %v.",
				status, v.expectedResponseStatus)
		}

		// If need to check details
		if v.checkDetails == true {
			// Check user details using updated object
			checkUserDetails(rr, &v.userToCheck, t, true)
		}
	}
}

// Note: This must be the final test as it updates details within the test connection
func TestUserController_UpdateMyProfile(t *testing.T) {
	var updateTests = []struct {
		data                   map[string]string
		tokenToUse             string
		expectedResponseStatus int
		checkDetails           bool
		loggedInDetails        db.User
	}{
		// Admin test
		{map[string]string{
			"Username": "JabarCindi",
			"Name":     "Bambaloonie",
		}, testConnection.adminToken, http.StatusOK, true, *testConnection.admin},
		// User test
		{map[string]string{
			"Username": "JabarHindi",
			"Name":     "Bambaloonie",
			"Password": "YeezusChris",
		}, testConnection.userToken, http.StatusOK, true, *testConnection.user},
		// User update Email with non-email
		{map[string]string{
			"Email": "JabarHindi",
		}, testConnection.userToken, http.StatusBadRequest, false, *testConnection.user},
		// User update Email with duplicate email
		{map[string]string{
			"Email": "juba@ymail.com",
		}, testConnection.userToken, http.StatusBadRequest, false, *testConnection.user},
		// User updates without token should fail
		{map[string]string{
			"Username": "JabarHindi",
			"Name":     "Bambaloonie",
			"Password": "YeezusChris",
		}, "", http.StatusForbidden, false, *testConnection.user},
		// Update for 2 tests below should be disallowed due to being too short
		{map[string]string{
			"Username": "Gobod",
			"Name":     "solu",
		}, testConnection.adminToken, http.StatusBadRequest, false, *testConnection.admin},
		{map[string]string{
			"Username": "Gabor",
			"Name":     "solu",
		}, testConnection.userToken, http.StatusBadRequest, false, *testConnection.user},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := "/api/me"

	for _, v := range updateTests {
		// Make new request with user update in body
		req, err := http.NewRequest("PUT", requestUrl, buildReqBody(v.data))
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()
		// Add user auth token to header
		req.Header.Set("Authorization", fmt.Sprintf("bearer %v", v.tokenToUse))
		// Send update request to mock server
		testConnection.router.ServeHTTP(rr, req)
		// Check response is failed for normal user to update another
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("User update test (%v): got %v want %v.", v.data["Username"],
				status, v.expectedResponseStatus)
		}

		// If need to check details
		if v.checkDetails == true {
			// Update created user struct with the changes pushed through API
			updateChangesOnly(&v.loggedInDetails, v.data)

			// Check user details using updated object
			checkUserDetails(rr, &v.loggedInDetails, t, true)
		}

		// Return updates to original state
		testConnection.serv.Update(int(testConnection.admin.ID), &models.UpdateUser{
			Username: testConnection.admin.Username,
			Password: testConnection.admin.Password,
			Email:    testConnection.admin.Email,
			Name:     testConnection.admin.Name,
		})
		testConnection.serv.Update(int(testConnection.user.ID), &models.UpdateUser{
			Username: testConnection.user.Username,
			Password: testConnection.user.Password,
			Email:    testConnection.user.Email,
			Name:     testConnection.user.Name,
		})
	}
}

func TestUserController_Login(t *testing.T) {
	var loginTests = []struct {
		title                  string
		data                   models.Login
		expectedResponseStatus int
		failureExpected        bool
		expectedMessage        string
	}{
		// Admin user login
		{"Admin user login", models.Login{
			Email:    testConnection.admin.Email,
			Password: testConnection.admin.Password,
		}, http.StatusOK, false, ""},
		// Admin user incorrect login
		{"Admin user incorrect", models.Login{
			Email:    testConnection.admin.Email,
			Password: "wrongPassword",
		}, http.StatusUnauthorized, true, "Incorrect username/password\n"},
		// Basic user login
		{"Basic user login", models.Login{
			Email:    testConnection.user.Email,
			Password: testConnection.user.Password,
		}, http.StatusOK, false, ""},
		// Basic user incorrect login
		{"Basic user incorrect", models.Login{
			Email:    testConnection.user.Email,
			Password: "VeryWrongPassword",
		}, http.StatusUnauthorized, true, "Incorrect username/password\n"},
		// Completely made up email for user login
		{"Non existent user login", models.Login{
			Email:    "jester@gmail.com",
			Password: "VeryWrongPassword",
		}, http.StatusUnauthorized, true, "Invalid Credentials\n"},
		// Email is not an email (Validation error, can't be checked below)
		// Should result in bad request
		{"Non existent user login", models.Login{
			Email:    "jester",
			Password: "VeryWrongPassword",
		}, http.StatusBadRequest, false, ""},
		// Empty credentials
		{"Non existent user login", models.Login{
			Email:    "jester",
			Password: "",
		}, http.StatusBadRequest, false, ""},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := "/api/users/login"

	for _, v := range loginTests {
		// Build request body
		reqBody := buildReqBody(v.data)
		// Make new request with user update in body
		req, err := http.NewRequest("POST", requestUrl, reqBody)
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Send update request to mock server
		testConnection.router.ServeHTTP(rr, req)

		// Check response is failed for normal user to update another
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("User login test (%v:%v): got %v want %v. Resp: %v", v.title, v.data.Password,
				status, v.expectedResponseStatus, rr.Body)
		}

		// If failure is expected
		if v.failureExpected {
			// Form req body
			reqBody := rr.Body.String()
			// Check if matches with expectation
			if reqBody != v.expectedMessage {
				t.Errorf("The body is: %v. expected: %v.", rr.Body.String(), v.expectedMessage)
			}

		}

	}
}

// Updates the parameter user struct with the updated values in updated user
func updateChangesOnly(createdUser *db.User, updatedUser map[string]string) error {
	// Iterate through map and change struct values
	for k, v := range updatedUser {
		// Update each struct field using map
		err := helpers.UpdateStructField(createdUser, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// Check the user details (username, name, email and ID)
func checkUserDetails(rr *httptest.ResponseRecorder, createdUser *db.User, t *testing.T, checkId bool) {
	// Convert response JSON to struct
	var body db.User
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Only check ID if parameter checkId is true
	if checkId == true {
		// Verify that the found user matches the original created user
		if body.ID != createdUser.ID {
			t.Errorf("found createdUser has incorrect ID: expected %d, got %d", createdUser.ID, body.ID)
		}
	}
	// Check updated details
	if body.Email != createdUser.Email {
		t.Errorf("found createdUser has incorrect email: expected %s, got %s", createdUser.Email, body.Email)
	}
	if body.Username != createdUser.Username {
		t.Errorf("found createdUser has incorrect username: expected %s, got %s", createdUser.Username, body.Username)
	}
	if body.Name != createdUser.Name {
		t.Errorf("found createdUser has incorrect name: expected %s, got %s", createdUser.Name, body.Name)
	}
}

// Build a struct object to a type of bytes.reader to fulfill io.reader interface
func buildReqBody(data interface{}) *bytes.Reader {
	// Marshal to JSON
	marshalled, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Failed to marshal JSON")
	}
	// Make into reader
	readerReqBody := bytes.NewReader(marshalled)
	return readerReqBody
}

// SETUP FUNCTIONS
//
// Setup enforcer and sync app state
func setupTestEnforcer(dbClient *gorm.DB) {
	// Build enforcer
	enforcer, err := auth.EnforcerSetup(dbClient)
	if err != nil {
		fmt.Println("Error building enforcer")
	}

	// Assign values in app for authentication
	app.DbClient = dbClient
	app.RBEnforcer = enforcer
	// Sync app in authentication package for usage in authentication functions
	auth.SetStateInAuth(&app)
}

// Build api router for testing
func buildRouter(c controller.UserController) *chi.Mux {
	// Create a new chi router
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {

		// Login
		r.Post("/api/users/login", c.Login)

		// Create new user
		r.Post("/api/users", c.Create)

		r.Group(func(r chi.Router) {
			r.Use(auth.AuthenticateJWT)

			// users
			r.Get("/api/users", c.FindAll)
			r.Get("/api/users/{id}", c.Find)
			r.Put("/api/users/{id}", c.Update)
			r.Delete("/api/users/{id}", c.Delete)

			// My profile
			r.Get("/api/me", c.GetMyUserDetails)
			r.Put("/api/me", c.UpdateMyProfile)

		})
	})
	return r
}

// setup database connection
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

// Generates a new user, changes its role to admin and returns it with token
func generateUserWithRoleAndToken(user *db.User, role string) (*db.User, string) {
	unhashedPass := user.Password
	createdUser, err := hashPassAndGenerateUserInDb(user)
	if err != nil {
		fmt.Print("Problem creating admin user for tests.")
	}
	// Update user to admin
	createdUser.Role = role
	updatedUser, err := testConnection.repo.Update(int(createdUser.ID), createdUser)
	// If match found (no errors)
	if err == nil {
		fmt.Println("Generating token for: ", updatedUser.Email)
		// Set login status to true
		tokenString, err := auth.GenerateJWT(int(updatedUser.ID), updatedUser.Email, updatedUser.Role)
		if err != nil {
			fmt.Println("Failed to create JWT")
		}

		// Add unhashed password to returned object
		updatedUser.Password = unhashedPass
		// Send to user in body
		return updatedUser, tokenString
	}
	return nil, ""
}

// Test helper function: Hashes password and generates a new user in the database
func hashPassAndGenerateUserInDb(user *db.User) (*db.User, error) {
	// Hash password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		fmt.Print("Couldn't hash password")
	}
	user.Password = string(hashedPass)

	// Create user
	createResult := testConnection.dbClient.Create(user)
	if createResult.Error != nil {
		fmt.Printf("Couldn't create user: %v", user.Email)
	}

	return user, nil
}
