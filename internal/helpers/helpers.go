package helpers

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/dmawardi/Go-Template/internal/models"
)

// Takes struct data and returns as JSON to Response writer
func WriteAsJSON(w http.ResponseWriter, data interface{}) error {
	// Edit content type
	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	// Write data as response
	w.Write(jsonData)
	return nil
}

// Extract base path from request
func ExtractBasePath(r *http.Request) string {
	// Extract current URL being accessed
	extractedPath := r.URL.Path
	// Split path
	fullPathArray := strings.Split(extractedPath, "/")

	// If the final item in the slice is determined to be numeric
	if govalidator.IsNumeric(fullPathArray[len(fullPathArray)-1]) {
		// Remove final element from slice
		fullPathArray = fullPathArray[:len(fullPathArray)-1]
	}
	// Join strings in slice for clean URL
	pathWithoutParameters := strings.Join(fullPathArray, "/")
	return pathWithoutParameters
}

// Takes in a list of errors from Go Validator and formats into JSON ready struct
func CreateStructFromValidationErrorString(errs []error) *models.ValidationError {
	// Prepare validation model for appending errors
	validation := &models.ValidationError{
		Validation_errors: make(map[string][]string),
	}
	// Loop through slice of errors
	for _, e := range errs {
		// Grab error strong
		errorString := e.Error()
		// Split by colon
		errorArray := strings.Split(errorString, ": ")

		// Prepare err message array for map preparation
		var errMessageArray []string
		// Append the error to the array
		errMessageArray = append(errMessageArray, errorArray[1])

		// Assign to validation struct
		validation.Validation_errors[errorArray[0]] = errMessageArray
	}

	return validation
}

// Uses a DTO struct's key value "valid" config to assess whether it's valid
// then returns a struct ready for JSON marshal
func GoValidateStruct(objectToValidate interface{}) (bool, *models.ValidationError) {
	// Validate the incoming DTO
	_, err := govalidator.ValidateStruct(objectToValidate)

	// if no error found
	if err != nil {
		// Prepare slice of errors
		errs := err.(govalidator.Errors).Errors()

		// Grabs the error slice and creates a front-end ready validation error
		validationResponse := CreateStructFromValidationErrorString(errs)
		// return failure on validation and validation response
		return false, validationResponse
	}
	// Return pass on validation and empty validation response
	return true, &models.ValidationError{}
}

// Update a struct field dynamically
func UpdateStructField(structPtr interface{}, fieldName string, fieldValue interface{}) error {
	value := reflect.ValueOf(structPtr)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return fmt.Errorf("invalid struct pointer")
	}

	structValue := value.Elem()
	if !structValue.CanSet() {
		return fmt.Errorf("cannot set struct field value")
	}

	field := structValue.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("invalid struct field name")
	}

	if !field.CanSet() {
		return fmt.Errorf("cannot set struct field value")
	}

	fieldValueRef := reflect.ValueOf(fieldValue)
	if !fieldValueRef.Type().AssignableTo(field.Type()) {
		return fmt.Errorf("field value type mismatch")
	}

	field.Set(fieldValueRef)
	return nil
}

// Used for pagination to return a nil value if the incoming int is 0
func ReturnNilIfZero(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

// Accepts request and a slice of conditions to extract from the request
func ExtractConditionParams(r *http.Request, userConditions []string) ([]string, error) {
	// Grab URL query parameters
	queryParams := r.URL.Query()

	// extract conditions from query params
	extractedConditions := []string{}
	// Iterate through conditions
	for _, condition := range userConditions {
		// Check if condition is present in query params
		if queryParams.Get(condition) != "" {
			// Extract value from query params
			queryValue := queryParams.Get(condition)
			// Replace double quotes with single quotes
			modQueryValue := strings.Replace(queryValue, `"`, "'", -1)
			// Build string condition
			stringCondition := fmt.Sprintf("%s = %s", condition, modQueryValue)
			// If present, append to slice
			extractedConditions = append(extractedConditions, stringCondition)
		}
	}
	return extractedConditions, nil
}

// Generates random string with n characters
func GenerateRandomString(n int) (string, error) {
	const lettersAndDigits = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Make a byte slice of n length
	bytes := make([]byte, n)

	// Fill byte slice with random bytes
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Replace each byte with a letter or digit
	for i, b := range bytes {
		bytes[i] = lettersAndDigits[b%byte(len(lettersAndDigits))]
	}

	// Return the random string
	return string(bytes), nil
}

// LoadTemplate parses an HTML template, executes it with the provided data, and returns the result as a string.
func LoadTemplate(templateFilePath string, data interface{}) (string, error) {
	// Parse the template file
	t, err := template.ParseFiles(templateFilePath)
	if err != nil {
		return "", err
	}

	// Build the template with the injected data
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

// Grabs a query parameter from the request, if not present, returns default value
func GrabQueryParamOrDefault(r *http.Request, param string, defaultValue string) string {
	// Grab query parameters
	queryParam := r.URL.Query().Get(param)
	// Check if limit is available, if not, set to default
	if queryParam == "" {
		queryParam = defaultValue
	}
	return queryParam
}

func GrabIntQueryParamOrDefault(r *http.Request, param string, defaultValue int) (int, error) {
	// Grab query parameters
	queryParam := r.URL.Query().Get(param)
	// Check if limit is available, if not, set to default
	if queryParam == "" {
		return defaultValue, nil
	}
	// Convert to int
	intQuery, err := strconv.Atoi(queryParam)
	if err != nil {
		return 0, err
	}
	return intQuery, nil
}
