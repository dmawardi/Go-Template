package helpers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
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
