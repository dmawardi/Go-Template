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
	checkPolicyDetails(t, body[0], policy1)

	// Delete policy
	err = testConnection.auth.serv.Delete(policy1)
	if err != nil {
		t.Error(err)
	}
}

func TestAuthController_Delete(t *testing.T) {
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
	requestUrl := "auth"

	req, err := buildApiRequest("DELETE", requestUrl, buildReqBody(policy1), true, testConnection.accounts.admin.token)
	if err != nil {
		t.Error(err)
	}
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Use handler with recorder and created request
	testConnection.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("%v: got %v want %v.\nResp:%s", "Auth Delete",
			status, http.StatusOK, rr.Body.String())
	}

	// Check if the record is deleted
	found, err := testConnection.auth.serv.FindByResource(policy1.Resource)

	if err != nil {
		t.Errorf("Error detected when finding resource: %v", err)
	}
	if len(found) > 0 {
		t.Errorf("Expected to not find resource, however, found: %v", found)
	}
}

func TestAuthController_Create(t *testing.T) {
	policy1 := models.PolicyRule{
		Role:     "admin",
		Resource: "/api/gustav",
		Action:   "read",
	}

	// Build slug
	requestUrl := "auth"

	req, err := buildApiRequest("POST", requestUrl, buildReqBody(policy1), true, testConnection.accounts.admin.token)
	if err != nil {
		t.Error(err)
	}
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Use handler with recorder and created request
	testConnection.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("%v: got %v want %v.\nResp:%s", "Auth Create",
			status, http.StatusOK, rr.Body.String())
	}

	// Check if the record is created
	found, err := testConnection.auth.serv.FindByResource(policy1.Resource)

	if err != nil {
		t.Errorf("Error detected when finding resource: %v", err)
	}
	if len(found) == 0 {
		t.Errorf("Expected to find resource, however, not found: %v", found)
	}

	// Check details
	checkPolicyDetails(t, found[0], policy1)

	// Delete policy
	err = testConnection.auth.serv.Delete(policy1)
	if err != nil {
		t.Error(err)
	}
}

func TestAuthController_Update(t *testing.T) {
	policy1 := models.PolicyRule{
		Role:     "admin",
		Resource: "/api/gustav",
		Action:   "read",
	}
	policy2 := models.PolicyRule{
		Role:     "admin",
		Resource: "/api/gustav",
		Action:   "update",
	}

	// Create policy
	err := testConnection.auth.serv.Create(policy1)
	if err != nil {
		t.Error(err)
	}

	// Build slug
	requestUrl := "auth"

	req, err := buildApiRequest("PUT", requestUrl, buildReqBody(models.UpdateCasbinRule{OldPolicy: policy1, NewPolicy: policy2}), true, testConnection.accounts.admin.token)
	if err != nil {
		t.Error(err)
	}
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Use handler with recorder and created request
	testConnection.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("%v: got %v want %v.\nResp:%s", "Auth Update",
			status, http.StatusOK, rr.Body.String())
	}

	// Check if the record is updated
	found, err := testConnection.auth.serv.FindByResource(policy2.Resource)

	if err != nil {
		t.Errorf("Error detected when finding resource: %v", err)
	}
	if len(found) == 0 {
		t.Errorf("Expected to find resource, however, not found: %v", found)
	}

	// Check details
	checkPolicyDetails(t, found[0], policy2)

	// Delete policy
	err = testConnection.auth.serv.Delete(policy2)
	if err != nil {
		t.Error(err)
	}
}

// Checks if the policy details are a match
func checkPolicyDetails(t *testing.T, body models.PolicyRuleCombinedActions, policy models.PolicyRule) {
	// Check details
	if body.Resource != policy.Resource {
		t.Errorf("Body: %+v. Expected %v, got %v", body, policy.Resource, body.Resource)
	}
	if body.Role != policy.Role {
		t.Errorf("Expected %v, got %v", policy.Role, body.Role)
	}
	// Check if array of strings contains the created record
	if helpers.ArrayContainsString(body.Action, policy.Action) == false {
		t.Errorf("Expected %v, got %v", policy.Action, body.Action)
	}
}
