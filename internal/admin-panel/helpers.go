package adminpanel

import (
	"net/http"
	"time"
)

// Table/Form Display
//
// Format time.Time fields for forms
func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339) // or another suitable format
}

// parseFormToMap parses the form data and converts it into a map[string]string
func parseFormToMap(r *http.Request) (map[string]string, error) {
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
