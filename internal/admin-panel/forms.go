package adminpanel

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dmawardi/Go-Template/internal/models"
)

// Used to build form in Go Templates
// Primary form data used in PagerenderData
type FormData struct {
	// Form title
	FormFields []FormField
	// Contains form action and method
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

// Used to store data to render form selectors in Go Templates
type FormFieldSelector struct {
	Value    string
	Label    string
	Selected bool
}

// Display of errors in form
type ErrorMessage string

// Form Details
type FormDetails struct {
	FormAction string
	FormMethod string
}

// Map used to group form selectors for a schema (eg. FormSelector["field_name"])
// Form selectors are used to index the form selectors as functions that return map[string]string
type FormSelectors map[string]func() []FormFieldSelector

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

// Used to populate FormField values with placeholder values found from form in request
func populateValuesWithForm(r *http.Request, form *[]FormField, fieldMap map[string]string) error {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		return errors.New("Error parsing form")
	}

	// Loop through fields and populate placeholders
	for i := range *form {
		// Get pointer to field
		field := &(*form)[i]
		// If the field exists in the map, populate the placeholder
		if val, ok := fieldMap[field.DbLabel]; ok {
			field.Value = val
		} else {
			return fmt.Errorf("field: %s not found in map", field.DbLabel)
		}
	}
	return nil
}

// Used to populate form field placeholders with data from database (that has been converted to map[string]string)
func populatePlaceholdersWithDBData(form *[]FormField, fieldMap map[string]string) error {
	// Loop through fields and populate placeholders
	for i := range *form {
		// Get pointer to field
		field := &(*form)[i]
		if field.Type == "select" {
			// Update selectors with current value selected
			setDefaultSelected(field.Selectors, fieldMap[field.DbLabel])
			// Else treat as ordinary input
		} else {
			// If the field exists in the map, populate the placeholder
			if val, ok := fieldMap[field.DbLabel]; ok {
				// Populate placeholder as value from field map
				field.Placeholder = val
			} else {
				return fmt.Errorf("field: %s not found in map", field.DbLabel)
			}
		}
	}

	return nil
}
