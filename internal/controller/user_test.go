package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmawardi/Go-Template/internal/models"
)

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
	req, err := buildApiRequest("GET", fmt.Sprintf("users/%v", createdUser.ID), nil, true, testConnection.accounts.admin.token)
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
	req, err := buildApiRequest("GET", "users?limit=10&offset=0&order=id_desc", nil, true, testConnection.accounts.admin.token)
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
		req, err := buildApiRequest("GET", fmt.Sprintf("users?limit=%v&offset=%v&order=%v", v.limit, v.offset, v.order), nil, true, testConnection.accounts.admin.token)
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
		req, err := buildApiRequest("DELETE", fmt.Sprintf("users/%v", createdUser.ID), nil, true, v.tokenToUse)
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
		req, err := buildApiRequest("PUT", fmt.Sprintf("users/%v", createdUser.ID), buildReqBody(v.data), true, v.tokenToUse)
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

			// Convert response JSON to struct
			var body models.UserWithRole
			json.Unmarshal(rr.Body.Bytes(), &body)

			// check user details for match
			CompareObjects(body, createdUser, t, []string{"ID", "Username", "Email", "Name"})
		}
	}

	// Delete the created user
	err = testConnection.users.serv.Delete(int(createdUser.ID))
	if err != nil {
		t.Fatalf("Error clearing created user")
	}
}

func TestUserController_Create(t *testing.T) {
	var tests = []struct {
		testName               string
		data                   models.CreateUser
		expectedResponseStatus int
	}{
		{"Successful user creation", models.CreateUser{
			Username: "Jabarnam",
			Email:    "gabor@ymail.com",
			Password: "password",
			Name:     "Bambaliya",
		}, http.StatusCreated},
		{"Successful user creation", models.CreateUser{
			Username: "Swalanim",
			Email:    "salvia@ymail.com",
			Password: "seradfasdf",
			Name:     "CreditTomyaA",
		}, http.StatusCreated},
		{"Failure: Not email", models.CreateUser{
			Username: "Yukon",
			Email:    "Sylvio",
			Password: "wowogsdfg",
			Name:     "Sosawsdfgsdfg",
		}, http.StatusBadRequest},
		{"Failure: Pass/Name field length", models.CreateUser{
			Username: "Jabarnam",
			Email:    "Cakawu@ymail.com",
			Password: "as",
			Name:     "df",
		}, http.StatusBadRequest},
		{"Failure: Duplicate user", models.CreateUser{
			Username: "Jabarnam",
			Email:    "Jabal@ymail.com",
			Password: "as",
			Name:     "df",
		}, http.StatusBadRequest},
	}

	for _, v := range tests {
		req, err := buildApiRequest("POST", "users", buildReqBody(v.data), false, "")

		// Make new request with user update in body
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Send update request to mock server
		testConnection.router.ServeHTTP(rr, req)
		// Check response is failed for normal user to update another
		if status := rr.Code; status != v.expectedResponseStatus {

			t.Errorf("%s: Got %v want %v.", v.data.Name,
				status, v.expectedResponseStatus)
		}

		// Init body for response extraction
		var body models.UserWithRole
		// Grab ID from response body
		json.Unmarshal(rr.Body.Bytes(), &body)

		// Delete the created user
		err = testConnection.users.serv.Delete(int(body.ID))
		if err != nil {
			t.Fatalf("Error clearing created user")
		}
	}
}

func TestUserController_GetMyUserDetails(t *testing.T) {
	var tests = []struct {
		testName               string
		checkDetails           bool
		tokenToUse             string
		userToCheck            models.UserWithRole
		expectedResponseStatus int
	}{
		{"User checking own profile", true, testConnection.accounts.user.token, *testConnection.accounts.user.details, http.StatusOK},
		{"Admin checking own profile", true, testConnection.accounts.admin.token, *testConnection.accounts.admin.details, http.StatusOK},
		// Deny access to user that doesn't have authentication
		{"Logged out user checking profile", false, "", models.UserWithRole{}, http.StatusForbidden},
	}
	// Create a request url with an "id" URL parameter
	requestUrl := "/api/me"

	for _, v := range tests {
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
			t.Errorf("Got %v want %v.",
				status, v.expectedResponseStatus)
		}

		// If need to check details
		if v.checkDetails == true {
			// Convert response JSON to struct
			var body models.UserWithRole
			json.Unmarshal(rr.Body.Bytes(), &body)

			// Check user details using updated object
			CompareObjects(body, &v.userToCheck, t, []string{"ID", "Username", "Email", "Name"})
		}
	}
}

func TestUserController_UpdateMyProfile(t *testing.T) {
	var tests = []struct {
		testName               string
		data                   map[string]string
		tokenToUse             string
		expectedResponseStatus int
		checkDetails           bool
		loggedInDetails        models.UserWithRole
	}{
		{"Admin self update", map[string]string{
			"Username": "JabarCindi",
			"Name":     "Bambaloonie",
		}, testConnection.accounts.admin.token, http.StatusOK, true, *testConnection.accounts.admin.details},
		{"User self update", map[string]string{
			"Username": "JabarHindi",
			"Name":     "Bambaloonie",
			"Password": "YeezusChris",
		}, testConnection.accounts.user.token, http.StatusOK, true, *testConnection.accounts.user.details},
		{"User self update with invalid email", map[string]string{
			"Email": "JabarHindi",
		}, testConnection.accounts.user.token, http.StatusBadRequest, false, *testConnection.accounts.user.details},
		{"Fail: User self update with duplicate email", map[string]string{
			"Username": "Swahili",
			"Email":    testConnection.accounts.admin.details.Email,
		}, testConnection.accounts.user.token, http.StatusBadRequest, false, *testConnection.accounts.user.details},
		{"Fail: User update without token", map[string]string{
			"Username": "JabarHindi",
			"Name":     "Bambaloonie",
			"Password": "YeezusChris",
		}, "", http.StatusForbidden, false, *testConnection.accounts.user.details},
		{"Fail: Admin update with invalid validation", map[string]string{
			"Username": "Gobod",
			"Name":     "solu",
		}, testConnection.accounts.admin.token, http.StatusBadRequest, false, *testConnection.accounts.admin.details},
		{"Fail: User update with invalid validation", map[string]string{
			"Username": "Gabor",
			"Name":     "solu",
		}, testConnection.accounts.user.token, http.StatusBadRequest, false, *testConnection.accounts.user.details},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := "/api/me"

	for _, v := range tests {
		// Make new request with user update in body
		req, err := http.NewRequest("PUT", requestUrl, buildReqBody(v.data))
		if err != nil {
			t.Fatal(err)
		}
		// Add user auth token to header
		req.Header.Set("Authorization", fmt.Sprintf("bearer %v", v.tokenToUse))

		// Create a response recorder
		rr := httptest.NewRecorder()
		// Send update request to mock server
		testConnection.router.ServeHTTP(rr, req)
		// Check response is failed for normal user to update another
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("%v: got %v want %v.", v.testName,
				status, v.expectedResponseStatus)
		}

		// If need to check details
		if v.checkDetails == true {
			// Update created user struct with the changes pushed through API
			UpdateModelFields(&v.loggedInDetails, v.data)

			// Check user details using updated object
			checkUserDetails(rr, &v.loggedInDetails, t, true)
			// Convert response JSON to struct
			var body models.UserWithRole
			json.Unmarshal(rr.Body.Bytes(), &body)

			CompareObjects(body, v.loggedInDetails, t, []string{"ID", "Username", "Email", "Name"})
		}

		// Return updates to original state
		testConnection.users.serv.Update(int(testConnection.accounts.admin.details.ID), &models.UpdateUser{
			Username: testConnection.accounts.admin.details.Username,
			Password: testConnection.accounts.admin.details.Password,
			Email:    testConnection.accounts.admin.details.Email,
			Name:     testConnection.accounts.admin.details.Name,
		})
		testConnection.users.serv.Update(int(testConnection.accounts.user.details.ID), &models.UpdateUser{
			Username: testConnection.accounts.user.details.Username,
			Password: testConnection.accounts.user.details.Password,
			Email:    testConnection.accounts.user.details.Email,
			Name:     testConnection.accounts.user.details.Name,
		})
	}
}

func TestUserController_Login(t *testing.T) {
	var tests = []struct {
		testName               string
		data                   models.Login
		expectedResponseStatus int
		failureExpected        bool
		expectedMessage        string
	}{
		{"Admin user login", models.Login{
			Email:    testConnection.accounts.admin.details.Email,
			Password: testConnection.accounts.admin.details.Password,
		}, http.StatusOK, false, ""},
		{"Fail: Admin user incorrect details", models.Login{
			Email:    testConnection.accounts.admin.details.Email,
			Password: "wrongPassword",
		}, http.StatusUnauthorized, true, "Invalid Credentials\n"},
		{"Basic user login", models.Login{
			Email:    testConnection.accounts.user.details.Email,
			Password: testConnection.accounts.user.details.Password,
		}, http.StatusOK, false, ""},
		{"Fail: Basic user incorrect details", models.Login{
			Email:    testConnection.accounts.user.details.Email,
			Password: "VeryWrongPassword",
		}, http.StatusUnauthorized, true, "Invalid Credentials\n"},
		{"Fail: Non existent user login", models.Login{
			Email:    "jester@gmail.com",
			Password: "VeryWrongPassword",
		}, http.StatusUnauthorized, true, "Invalid Credentials\n"},
		{"Fail: Invalid email user login", models.Login{
			Email:    "jester",
			Password: "VeryWrongPassword",
		}, http.StatusBadRequest, false, ""},
		{"Fail: Empty credentials", models.Login{
			Email:    "jester",
			Password: "",
		}, http.StatusBadRequest, false, ""},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := "/api/users/login"

	for _, v := range tests {
		// Make request with update in body
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
			t.Errorf("%v: %v/%v. got %v want %v. Resp: %v", v.testName, v.data.Email, v.data.Password,
				status, v.expectedResponseStatus, rr.Body)
		}

		// If failure is expected
		if v.failureExpected {
			// Form req body
			reqBody := rr.Body.String()
			// Check if matches with expectation
			if reqBody != v.expectedMessage {
				t.Errorf("%v: The body is: %v. expected: %v.", v.testName, rr.Body.String(), v.expectedMessage)
			}

		}

	}
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
