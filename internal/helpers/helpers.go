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

// Special handling for search query if found
func addSearchQueryToConditions(r *http.Request, conditionsToExtract map[string]string, currentConditions []interface{}) []interface{} {
	// Prepare URL query parameters
	queryParams := r.URL.Query()
	// Prepare search query
	searchQuery := ""
	if searchValue := queryParams.Get("search"); searchValue != "" {
		searchQuery = "%" + searchValue + "%"
		// Iterate through query list to add the search condition to each as a LIKE query
		for param, conditionType := range conditionsToExtract {
			// If query parameter is string
			if conditionType == "string" {
				// Add query to conditions with param name and make case insensitive
				lowerCaseValue := strings.ToLower(fmt.Sprintf("%v", searchQuery))
				currentConditions = append(currentConditions, fmt.Sprintf("LOWER(%s) LIKE ?", param), lowerCaseValue)
			}
		}
	}

	return currentConditions
}

// Accepts request and a slice of conditions to extract from the request
// Extracts as a slice of interfaces that are structured as [condition, value]
func ExtractConditionParams(r *http.Request, conditionsToExtract map[string]string) ([]interface{}, error) {
	// Prepare URL query parameters
	queryParams := r.URL.Query()
	// Prepare slice of conditions
	var extractedConditions []interface{}

	// Special handling for search query if found
	extractedConditions = addSearchQueryToConditions(r, conditionsToExtract, extractedConditions)

	// Iterate through query list
	for param, conditionType := range conditionsToExtract {
		// If query parameter is present and not empty
		if queryValue := queryParams.Get(param); queryValue != "" {
			// Prepare variables
			var condition string
			var value interface{}
			var err error

			// Detecting prefixes for operators
			switch {
			// If greater than ie. age=gt:20
			case strings.HasPrefix(queryValue, "gt:"):
				condition = param + " > ?"
				value, err = parseValue(queryValue[3:], conditionType)
			// If less than ie. age=lt:20
			case strings.HasPrefix(queryValue, "lt:"):
				condition = param + " < ?"
				value, err = parseValue(queryValue[3:], conditionType)
			// If greater than or equal to ie. age=gte:20
			case strings.HasPrefix(queryValue, "gte:"):
				condition = param + " >= ?"
				value, err = parseValue(queryValue[4:], conditionType)
			// If less than or equal to ie. age=lte:20
			case strings.HasPrefix(queryValue, "lte:"):
				condition = param + " <= ?"
				value, err = parseValue(queryValue[4:], conditionType)
			// If default equal condition
			default:
				condition = param + " = ?"
				value, err = parseValue(queryValue, conditionType)
			}

			// If an issue found, return error
			if err != nil {
				return nil, fmt.Errorf("invalid value for %s: %v", param, err)
			}

			// Append to slice
			extractedConditions = append(extractedConditions, condition, value)
		}
	}

	return extractedConditions, nil
}

// Helper function to parse the value based on type
func parseValue(value, conditionType string) (interface{}, error) {
	switch conditionType {
	case "int":
		return strconv.Atoi(value)
	case "string":
		return value, nil
	case "bool":
		return strconv.ParseBool(value)
	default:
		return nil, fmt.Errorf("unknown condition type: %s", conditionType)
	}
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

// Grabs an INT type query parameter from the request, if not present, returns default value
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
