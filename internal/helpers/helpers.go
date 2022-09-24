package helpers

import (
	"encoding/json"
	"net/http"
	"strings"
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
	object := r.URL.Path
	// Split path
	fullPathArray := strings.Split(object, "/")
	// Remove final element from slice
	amendedPathArray := fullPathArray[:len(fullPathArray)-1]
	// Join strings in slice for clean URL
	pathWithoutParameters := strings.Join(amendedPathArray, "/")
	return pathWithoutParameters
}
