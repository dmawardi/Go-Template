package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		fmt.Printf("The error string for this line is: %s", errorString)
		// Split by colon
		errorArray := strings.Split(errorString, ": ")

		// Prepare err message array for map preparation
		var errMessageArray []string
		// Append the error to the array
		errMessageArray = append(errMessageArray, errorArray[1])

		// Assign to validation struct
		validation.Validation_errors[errorArray[0]] = errMessageArray
	}

	fmt.Printf("validation errors: %v", validation.Validation_errors)
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
