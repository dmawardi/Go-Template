package controller_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
)

// Builds request for user module. Method (GET/POST/etc.) url suffix is for any values that may come after /api/users
// body is for the request body (nil for none), authHeaderRequired is for if the request requires an authorization header and token is the token to use
func buildUserRequest(method string, urlSuffix string, body io.Reader, authHeaderRequired bool, token string) (request *http.Request, err error) {
	req, err := http.NewRequest(method, fmt.Sprintf("/api/users%v", urlSuffix), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	// If authorization header required
	if authHeaderRequired {
		req.Header.Set("Authorization", fmt.Sprintf("bearer %v", token))
	}
	return req, nil
}

func TestUserController_Find(t *testing.T) {
	// Create user
	createdUser, err := testConnection.users.serv.Create(&models.CreateUser{
		Username: "Jabar",
		Email:    "greenie@ymail.com",
		Password: "password",
		Name:     "Bamba",
	})
	if err != nil {
		t.Fatalf("failed to create test user for test: %v", err)
	}

	// Create a request with an "id" URL parameter
	req, err := buildUserRequest("GET", fmt.Sprintf("/%v", createdUser.ID), nil, true, testConnection.accounts.admin.token)
	if err != nil {
		t.Fatal(err)
	}
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
	var body models.UserWithRole
	json.Unmarshal(rr.Body.Bytes(), &body)

	// check user details for match
	CompareObjects(body, createdUser, t, []string{"ID", "Username", "Email", "Name"})

	// Delete the created user
	delResult := testConnection.users.serv.Delete(int(createdUser.ID))
	if delResult != nil {
		t.Fatalf("Error clearing created user")
	}
}

func TestUserController_FindAll(t *testing.T) {
	// Test finding two already created users for authentication mocking
	// Create a new request
	req, err := buildUserRequest("GET", "?limit=10&offset=0&order=id_desc", nil, true, testConnection.accounts.admin.token)
	if err != nil {
		t.Fatal(err)
	}

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
	var body *models.PaginatedUsers
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Check length of user array
	if len(*body.Data) != 2 {
		t.Errorf("Users array in findAll failed: expected %d, got %d", 2, len(*body.Data))
	}

	// Iterate through users array received
	for _, item := range *body.Data {
		// If id is admin id
		if item.ID == testConnection.accounts.admin.details.ID {
			CompareObjects(item, testConnection.accounts.admin.details, t, []string{"ID", "Username", "Email", "Name"})

		} else {
			CompareObjects(item, testConnection.accounts.user.details, t, []string{"ID", "Username", "Email", "Name"})

		}
	}

	// Test parameter input
	var failParameterTests = []struct {
		test_name              string
		limit                  string
		offset                 string
		order                  string
		expectedResponseStatus int
	}{
		// Only limit
		{test_name: "Only limit", limit: "5", offset: "", order: "", expectedResponseStatus: http.StatusOK},
		// Check normal parameters functional with order by
		{test_name: "Normal test", limit: "20", offset: "1", order: "id", expectedResponseStatus: http.StatusOK},
		// Descending order
		{test_name: "Normal test", limit: "20", offset: "1", order: "id_desc", expectedResponseStatus: http.StatusOK},
	}
	for _, v := range failParameterTests {
		req, err := buildUserRequest("GET", fmt.Sprintf("?limit=%v&offset=%v&order=%v", v.limit, v.offset, v.order), nil, true, testConnection.accounts.admin.token)
		if err != nil {
			t.Fatal(err)
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Use handler with recorder and created request
		testConnection.router.ServeHTTP(rr, req)

		// Check the response status code
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("In test '%s': handler returned wrong status code: got %v want %v", v.test_name,
				status, v.expectedResponseStatus)
		}
	}
}

func TestUserController_Delete(t *testing.T) {
	// Create user
	createdUser, err := testConnection.users.serv.Create(&models.CreateUser{
		Username: "Jabar",
		Email:    "zubayle@ymail.com",
		Password: "password",
		Name:     "Bamba",
	})
	if err != nil {
		t.Fatalf("failed to create user for test: %v", err)
	}

	// Test parameter input
	var tests = []struct {
		testName               string
		tokenToUse             string
		expectedResponseStatus int
	}{
		{testName: "Normal user delete failure", tokenToUse: testConnection.accounts.user.token, expectedResponseStatus: http.StatusForbidden},
		// Put last to also replace test user deletion
		{testName: "Admin user delete success", tokenToUse: testConnection.accounts.admin.token, expectedResponseStatus: http.StatusOK},
	}

	for _, v := range tests {
		// Create a request
		req, err := buildUserRequest("DELETE", fmt.Sprintf("/%v", createdUser.ID), nil, true, v.tokenToUse)
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Use handler with recorder and created request
		testConnection.router.ServeHTTP(rr, req)

		// Check response is failed for normal user
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("Deletion test: got %v want %v.",
				status, v.expectedResponseStatus)
		}

	}
}

func TestUserController_Update(t *testing.T) {
	// Create user
	createdUser, err := testConnection.users.serv.Create(&models.CreateUser{
		Username: "Jabar",
		Email:    "greenthumb@ymail.com",
		Password: "password",
		Name:     "Bamba",
	})
	if err != nil {
		t.Fatalf("failed to create test user for test: %v", err)
	}

	var updateTests = []struct {
		testName string
		// To be converted to string for URL
		urlExtension           interface{}
		data                   map[string]string
		tokenToUse             string
		expectedResponseStatus int
		checkDetails           bool
	}{
		{"Fail case: User updating another user", createdUser.ID, map[string]string{
			"Username": "JabarHindi",
			"Name":     "Bambaloonie",
		}, testConnection.accounts.user.token, http.StatusForbidden, false},
		{"Admin updating a user", createdUser.ID, map[string]string{
			"Username": "JabarHindi",
			"Name":     "Bambaloonie",
		}, testConnection.accounts.admin.token, http.StatusOK, true},
		// Update should be disallowed due to being too short
		{"Fail case: Update with validation errors", createdUser.ID, map[string]string{
			"Username": "Gobod",
			"Name":     "solu",
		}, testConnection.accounts.admin.token, http.StatusBadRequest, false},
		// User should be forbidden before validating
		{"Fail case: User invalid update should fail due to permissions", createdUser.ID, map[string]string{
			"Username": "Gobod",
			"Name":     "solu",
		}, testConnection.accounts.user.token, http.StatusForbidden, false},
		// Should fail as url extension is incorrect
		{"Fail case: Bad url parameter", "x", map[string]string{
			"Username": "Gobod",
			"Name":     "solu",
		}, testConnection.accounts.admin.token, http.StatusBadRequest, false},
		// Should fail as id is above currently created
		{"Fail case: Bad url parameter", "99", map[string]string{
			"Username": "Gobod",
			"Name":     "solu",
		}, testConnection.accounts.admin.token, http.StatusBadRequest, false},
	}

	for _, v := range updateTests {
		req, err := buildUserRequest("PUT", fmt.Sprintf("/%v", createdUser.ID), buildReqBody(v.data), true, v.tokenToUse)
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()
		// Send update request to mock server
		testConnection.router.ServeHTTP(rr, req)
		// Check response expected vs received
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("Got %v want %v.",
				status, v.expectedResponseStatus)
		}

		// If need to check details
		if v.checkDetails == true {
			// Update created user struct with the changes pushed through API
			UpdateModelFields(createdUser, v.data)

			// Check user details using updated object
			checkUserDetails(rr, createdUser, t, true)
		}
	}

	// Delete the created user
	testConnection.dbClient.Delete(createdUser)
}

// func TestUserController_Create(t *testing.T) {
// 	var updateTests = []struct {
// 		data                   models.CreateUser
// 		expectedResponseStatus int
// 	}{
// 		{models.CreateUser{
// 			Username: "Jabarnam",
// 			Email:    "gabor@ymail.com",
// 			Password: "password",
// 			Name:     "Bambaliya",
// 		}, http.StatusCreated},
// 		{models.CreateUser{
// 			Username: "Swalanim",
// 			Email:    "salvia@ymail.com",
// 			Password: "seradfasdf",
// 			Name:     "CreditTomyaA",
// 		}, http.StatusCreated},
// 		// Create should be disallowed due to not being email
// 		{models.CreateUser{
// 			Username: "Yukon",
// 			Email:    "Sylvio",
// 			Password: "wowogsdfg",
// 			Name:     "Sosawsdfgsdfg",
// 		}, http.StatusBadRequest},
// 		// Should be a bad request due to pass/name length
// 		{models.CreateUser{
// 			Username: "Jabarnam",
// 			Email:    "Cakawu@ymail.com",
// 			Password: "as",
// 			Name:     "df",
// 		}, http.StatusBadRequest},
// 		// Should be a bad request due to duplicate user (created in init)
// 		{models.CreateUser{
// 			Username: "Jabarnam",
// 			Email:    "Jabal@ymail.com",
// 			Password: "as",
// 			Name:     "df",
// 		}, http.StatusBadRequest},
// 	}

// 	// Create a request url with an "id" URL parameter
// 	requestUrl := "/api/users"

// 	for _, v := range updateTests {
// 		// Make new request with user update in body
// 		req, err := http.NewRequest("POST", requestUrl, buildReqBody(v.data))
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		// Create a response recorder
// 		rr := httptest.NewRecorder()

// 		// Send update request to mock server
// 		testConnection.router.ServeHTTP(rr, req)
// 		// Check response is failed for normal user to update another
// 		if status := rr.Code; status != v.expectedResponseStatus {

// 			t.Errorf("User update test (%v): got %v want %v.", v.data.Name,
// 				status, v.expectedResponseStatus)
// 		}

// 		// Init body for response extraction
// 		var body db.User
// 		// Grab ID from response body
// 		json.Unmarshal(rr.Body.Bytes(), &body)

// 		// Delete the created user
// 		testConnection.dbClient.Delete(&db.User{ID: uint(body.ID)})
// 		// testConnection.users.serv.Delete(int(body.ID))
// 	}
// }

// func TestUserController_GetMyUserDetails(t *testing.T) {
// 	var updateTests = []struct {
// 		checkDetails           bool
// 		tokenToUse             string
// 		userToCheck            db.User
// 		expectedResponseStatus int
// 	}{
// 		{true, testConnection.accounts.user.token, *testConnection.accounts.user.details, http.StatusOK},
// 		{true, testConnection.accounts.admin.token, *testConnection.accounts.admin.details, http.StatusOK},
// 		// Deny access to user that doesn't have authentication
// 		{false, "", db.User{}, http.StatusForbidden},
// 	}
// 	// Create a request url with an "id" URL parameter
// 	requestUrl := "/api/me"

// 	for _, v := range updateTests {
// 		// Make new request with user update in body
// 		req, err := http.NewRequest("GET", requestUrl, nil)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		// Create a response recorder
// 		rr := httptest.NewRecorder()

// 		// If you need to check details for successful requests, set token
// 		if v.checkDetails {
// 			// Add user auth token to header
// 			req.Header.Set("Authorization", fmt.Sprintf("bearer %v", v.tokenToUse))
// 		}
// 		// Send update request to mock server
// 		testConnection.router.ServeHTTP(rr, req)
// 		// Check response is failed for normal user to update another
// 		if status := rr.Code; status != v.expectedResponseStatus {
// 			t.Errorf("User update test: got %v want %v.",
// 				status, v.expectedResponseStatus)
// 		}

// 		// If need to check details
// 		if v.checkDetails == true {
// 			// Check user details using updated object
// 			checkUserDetails(rr, &v.userToCheck, t, true)
// 		}
// 	}
// }

// func TestUserController_UpdateMyProfile(t *testing.T) {
// 	var updateTests = []struct {
// 		data                   map[string]string
// 		tokenToUse             string
// 		expectedResponseStatus int
// 		checkDetails           bool
// 		loggedInDetails        db.User
// 	}{
// 		// Admin test
// 		{map[string]string{
// 			"Username": "JabarCindi",
// 			"Name":     "Bambaloonie",
// 		}, testConnection.accounts.admin.token, http.StatusOK, true, *testConnection.accounts.admin.details},
// 		// User test
// 		{map[string]string{
// 			"Username": "JabarHindi",
// 			"Name":     "Bambaloonie",
// 			"Password": "YeezusChris",
// 		}, testConnection.accounts.user.token, http.StatusOK, true, *testConnection.accounts.user.details},
// 		// User update Email with non-email
// 		{map[string]string{
// 			"Email": "JabarHindi",
// 		}, testConnection.accounts.user.token, http.StatusBadRequest, false, *testConnection.accounts.user.details},
// 		// User update Email with duplicate email
// 		{map[string]string{
// 			"Username": "Swahili",
// 			"Email":    testConnection.accounts.admin.details.Email,
// 		}, testConnection.accounts.user.token, http.StatusBadRequest, false, *testConnection.accounts.user.details},
// 		// User updates without token should fail
// 		{map[string]string{
// 			"Username": "JabarHindi",
// 			"Name":     "Bambaloonie",
// 			"Password": "YeezusChris",
// 		}, "", http.StatusForbidden, false, *testConnection.accounts.user.details},
// 		// Update for 2 tests below should be disallowed due to being too short
// 		{map[string]string{
// 			"Username": "Gobod",
// 			"Name":     "solu",
// 		}, testConnection.accounts.admin.token, http.StatusBadRequest, false, *testConnection.accounts.admin.details},
// 		{map[string]string{
// 			"Username": "Gabor",
// 			"Name":     "solu",
// 		}, testConnection.accounts.user.token, http.StatusBadRequest, false, *testConnection.accounts.user.details},
// 	}

// 	// Create a request url with an "id" URL parameter
// 	requestUrl := "/api/me"

// 	for _, v := range updateTests {
// 		// Make new request with user update in body
// 		req, err := http.NewRequest("PUT", requestUrl, buildReqBody(v.data))
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		// Create a response recorder
// 		rr := httptest.NewRecorder()
// 		// Add user auth token to header
// 		req.Header.Set("Authorization", fmt.Sprintf("bearer %v", v.tokenToUse))
// 		// Send update request to mock server
// 		testConnection.router.ServeHTTP(rr, req)
// 		// Check response is failed for normal user to update another
// 		if status := rr.Code; status != v.expectedResponseStatus {
// 			t.Errorf("User update test (%v): got %v want %v.", v.data["Username"],
// 				status, v.expectedResponseStatus)
// 		}

// 		// If need to check details
// 		if v.checkDetails == true {
// 			// Update created user struct with the changes pushed through API
// 			updateChangesOnly(&v.loggedInDetails, v.data)

// 			// Check user details using updated object
// 			checkUserDetails(rr, &v.loggedInDetails, t, true)
// 		}

// 		// Return updates to original state
// 		testConnection.users.serv.Update(int(testConnection.accounts.admin.details.ID), &models.UpdateUser{
// 			Username: testConnection.accounts.admin.details.Username,
// 			Password: testConnection.accounts.admin.details.Password,
// 			Email:    testConnection.accounts.admin.details.Email,
// 			Name:     testConnection.accounts.admin.details.Name,
// 		})
// 		testConnection.users.serv.Update(int(testConnection.accounts.user.details.ID), &models.UpdateUser{
// 			Username: testConnection.accounts.user.details.Username,
// 			Password: testConnection.accounts.user.details.Password,
// 			Email:    testConnection.accounts.user.details.Email,
// 			Name:     testConnection.accounts.user.details.Name,
// 		})
// 	}
// }

// func TestUserController_Login(t *testing.T) {
// 	var loginTests = []struct {
// 		title                  string
// 		data                   models.Login
// 		expectedResponseStatus int
// 		failureExpected        bool
// 		expectedMessage        string
// 	}{
// 		// Admin user login
// 		{"Admin user login", models.Login{
// 			Email:    testConnection.accounts.admin.details.Email,
// 			Password: testConnection.accounts.admin.details.Password,
// 		}, http.StatusOK, false, ""},
// 		// Admin user incorrect login
// 		{"Admin user incorrect", models.Login{
// 			Email:    testConnection.accounts.admin.details.Email,
// 			Password: "wrongPassword",
// 		}, http.StatusUnauthorized, true, "Incorrect username/password\n"},
// 		// Basic user login
// 		{"Basic user login", models.Login{
// 			Email:    testConnection.accounts.user.details.Email,
// 			Password: testConnection.accounts.user.details.Password,
// 		}, http.StatusOK, false, ""},
// 		// Basic user incorrect login
// 		{"Basic user incorrect", models.Login{
// 			Email:    testConnection.accounts.user.details.Email,
// 			Password: "VeryWrongPassword",
// 		}, http.StatusUnauthorized, true, "Incorrect username/password\n"},
// 		// Completely made up email for user login
// 		{"Non existent user login", models.Login{
// 			Email:    "jester@gmail.com",
// 			Password: "VeryWrongPassword",
// 		}, http.StatusUnauthorized, true, "Invalid Credentials\n"},
// 		// Email is not an email (Validation error, can't be checked below)
// 		// Should result in bad request
// 		{"Non existent user login", models.Login{
// 			Email:    "jester",
// 			Password: "VeryWrongPassword",
// 		}, http.StatusBadRequest, false, ""},
// 		// Empty credentials
// 		{"Non existent user login", models.Login{
// 			Email:    "jester",
// 			Password: "",
// 		}, http.StatusBadRequest, false, ""},
// 	}

// 	// Create a request url with an "id" URL parameter
// 	requestUrl := "/api/users/login"

// 	for _, v := range loginTests {
// 		// Build request body
// 		reqBody := buildReqBody(v.data)
// 		// Make new request with user update in body
// 		req, err := http.NewRequest("POST", requestUrl, reqBody)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		// Create a response recorder
// 		rr := httptest.NewRecorder()

// 		// Send update request to mock server
// 		testConnection.router.ServeHTTP(rr, req)

// 		// Check response is failed for normal user to update another
// 		if status := rr.Code; status != v.expectedResponseStatus {
// 			t.Errorf("User login test (%v)\nDetails: %v/%v. got %v want %v. Resp: %v", v.title, v.data.Email, v.data.Password,
// 				status, v.expectedResponseStatus, rr.Body)
// 		}

// 		// If failure is expected
// 		if v.failureExpected {
// 			// Form req body
// 			reqBody := rr.Body.String()
// 			// Check if matches with expectation
// 			if reqBody != v.expectedMessage {
// 				t.Errorf("The body is: %v. expected: %v.", rr.Body.String(), v.expectedMessage)
// 			}

// 		}

// 	}
// }

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
func checkUserDetails(rr *httptest.ResponseRecorder, createdUser *models.UserWithRole, t *testing.T, checkId bool) {
	// Convert response JSON to struct
	var body models.UserWithRole
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
