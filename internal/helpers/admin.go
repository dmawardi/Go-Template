package helpers

import (
	"bytes"
	"html/template"
	"net/http"
)

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

// parseFormToMap parses the form data and converts it into a map[string]string
func ParseFormToMap(r *http.Request) (map[string]string, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	formMap := make(map[string]string)
	for key, values := range r.Form { // range over map
		// In form data, key can have multiple values,
		// we'll take the first one only
		formMap[key] = values[0]
	}

	return formMap, nil
}
