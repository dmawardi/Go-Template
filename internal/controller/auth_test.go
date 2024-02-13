package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
)

func TestAuthController_FindAll(t *testing.T) {
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
	var body []models.PolicyRuleCombinedActions
	json.Unmarshal(rr.Body.Bytes(), &body)

	match := helpers.CheckSliceType(body, reflect.TypeOf(models.PolicyRuleCombinedActions{}))

	if match == false {
		t.Errorf("Expected %v, got %v", reflect.TypeOf(models.PolicyRuleCombinedActions{}), reflect.TypeOf(body))
	}

}

func TestAuthController_FindByResource(t *testing.T) {
	policy1 := models.PolicyRule{
		Role:     "admin",
		Resource: "/api/gustav",
		Action:   "read",
	}
	// Create policy
	err := testConnection.auth.serv.Create(policy1)
	if err != nil {
		t.Error(err)
	}
	// Build slug
	slug := helpers.SlugifyResourceName(policy1.Resource)
	requestUrl := fmt.Sprintf("auth/%s", slug)

	req, err := buildApiRequest("GET", requestUrl, nil, true, testConnection.accounts.admin.token)
	if err != nil {
		t.Error(err)
	}
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Use handler with recorder and created request
	testConnection.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("%v: got %v want %v", "Auth Find by resource",
			status, http.StatusOK)
	}

	// Convert response JSON to struct
	var body []models.PolicyRuleCombinedActions
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Check details
	if body[0].Resource != policy1.Resource {
		t.Errorf("Body: %+v. Expected %v, got %v", rr.Body.String(), policy1.Resource, body[0].Resource)
	}
	if body[0].Role != policy1.Role {
		t.Errorf("Expected %v, got %v", policy1.Role, body[0].Role)
	}
	// Check if array of strings contains the created record
	if helpers.ArrayContainsString(body[0].Action, policy1.Action) == false {
		t.Errorf("Expected %v, got %v", policy1.Action, body[0].Action)
	}
}
