package adminpanel

import (
	"fmt"
	"strings"

	"github.com/dmawardi/Go-Template/internal/models"
)

// Go Template components

type FormData struct {
	FormFields  []FormField
	FormDetails FormDetails
}

// Data for each form field
type FormField struct {
	// Label above input
	Label   string
	DbLabel string
	// Used for form submission
	Name string
	// Is this field required?
	Required bool
	// Is this field disabled?
	Disabled bool

	// Current value
	Value string
	// Silhouette
	Placeholder string
	Type        string
	Errors      []ErrorMessage
	Selectors   []FormFieldSelector
}

type FormFieldSelector struct {
	Value    string
	Label    string
	Selected bool
}

// Display of errors in form
type ErrorMessage string

type UserEditForm struct {
	Title  string
	Fields []FormField
}

// Form Details
type FormDetails struct {
	FormAction string
	FormMethod string
}

// Sets the Errors field in each field of a form
func SetValidationErrorsInForm(form []FormField, validationErrors models.ValidationError) {
	// Iterate through fields
	for i, field := range form {
		// Check if field name is in validation errors
		if errors, ok := validationErrors.Validation_errors[field.Name]; ok {
			// If found, iterate through errors
			for _, err := range errors {
				// Get error message
				errorMessage := ErrorMessage(err)
				// If error contains does not, remove validated text
				if strings.Contains(err, "does not") {
					// Split string using does
					split := strings.Split(err, "does")
					// Update error message with rebuilt string
					errorMessage = ErrorMessage(fmt.Sprintf("Does %s", split[1]))
				}
				// Append error message to field
				form[i].Errors = append(form[i].Errors, errorMessage)
			}
		}
	}
}

// Adds an error to a specific field in a form
func addErrorToField(fields []FormField, fieldName string, newError ErrorMessage) {
	// Iterate through fields
	for i, field := range fields {
		// Check if field name matches
		if field.Name == fieldName {
			// If found, append error message to field
			fields[i].Errors = append(fields[i].Errors, newError)
			break // assuming only one field can have the matching name
		}
	}
}

func addDefaultSelectedToSelector(selector []FormFieldSelector, currentValue string) []FormFieldSelector {
	// Iterate through selector
	for i := range selector {
		// Check if value matches current value
		if selector[i].Value == currentValue {
			// If found, set selected to true
			selector[i].Selected = true

		} else {
			// Else, set selected to false
			selector[i].Selected = false
		}
	}
	// Return completed selector
	return selector
}
