package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmawardi/Go-Template/internal/models"
)

func TestAuthController_FindAll(t *testing.T) {
	policy1 := &models.PolicyRule{
		Role:     "user",
		Resource: "/api/v1/user",
		Action:   "read",
	}

	policy2 := &models.PolicyRule{
		Role:     "admin",
		Resource: "/api/v1/user",
		Action:   "create",
	}
	// Create 2 policies
	err := testConnection.auth.serv.Create(*policy1)
	if err != nil {
		t.Error(err)
	}
	err = testConnection.auth.serv.Create(*policy2)
	if err != nil {
		t.Error(err)
	}

	req, err := buildApiRequest("GET", "auth", nil, true, testConnection.accounts.admin.token)
	if err != nil {
		t.Error(err)
	}
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Use handler with recorder and created request
	testConnection.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("%v: got %v want %v", "Auth Find all",
			status, http.StatusOK)
	}

	// Convert response JSON to struct
	var body []map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &body)

	t.Fatalf(fmt.Sprintf("%v", body))

}
