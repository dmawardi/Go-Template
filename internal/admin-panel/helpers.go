package adminpanel

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Display
//
// Format time.Time fields for forms
func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339) // or another suitable format
}

// CONTAINS FUNCTIONS
//
// Checks if a string contains another string (Used to search for resource)
func containsString(s, searchTerm string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(searchTerm))
}

// Checks if array contains a particular string value
func arrayContainsString(arr []string, target string) bool {
	for _, s := range arr {
		if s == target {
			return true
		}
	}
	return false
}

// SORT FUNCTIONS
//
// Function to sort permissions data from enforcer
func sortMapStringInterfaceAlphabetically(a, b map[string]interface{}, key string) bool {
	resourceA, okA := a[key].(string)
	resourceB, okB := b[key].(string)

	// If either of the elements doesn't have a valid "resource" string, consider it greater (move it to the end)
	if !okA || !okB {
		return false
	}

	// Compare the "resource" strings alphabetically
	return resourceA < resourceB
}

// Function to sort a map[string]string by a given key
func sortMapStringStringAlphabetically(a, b map[string]string, key string) bool {
	valueA, okA := a[key]
	valueB, okB := b[key]

	// If either of the elements doesn't have a valid string for the given key, consider it greater (move it to the end)
	if !okA || !okB {
		return false
	}

	// Compare the strings alphabetically
	return valueA < valueB
}

// CONVERT FUNCTIONS
//
// Converts list of strings to list of ints
func convertStringSliceToIntSlice(stringSlice []string) ([]int, error) {
	intSlice := make([]int, 0, len(stringSlice)) // Create a slice of ints with the same length

	for _, str := range stringSlice {
		num, err := strconv.Atoi(str) // Convert string to int
		if err != nil {
			return nil, err
		}
		intSlice = append(intSlice, num) // Append the converted int to the slice
	}
	return intSlice, nil
}

// Convert json string to map[string]string
func jsonToMap(jsonStr string) (map[string]string, error) {
	var result map[string]string
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Convert struct string to map. Struct string will be in format: "[key1]:value1|[key2]:value2"
func stringToMap(input string) (map[string]string, error) {
	result := make(map[string]string)

	// Split the input string by "|"
	parts := strings.Split(input, "|")

	for _, part := range parts {
		// Check if the part has "[]" to identify a key name
		keyValueSlice := strings.Split(part, ":")

		// If a key value pair is found
		if len(keyValueSlice) == 2 {
			// Grab the first item in slice as key, and remove the "[" and "]" characters
			key := keyValueSlice[0]
			// // Grab the second item in slice as value
			value := keyValueSlice[1]
			// Add key value pair to result map
			result[key[1:len(key)-1]] = value
		}

	}

	return result, nil
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

// Auth Cookie functions
//
// Create and set jwt token for SSR authentication
func createAndSetHeaderCookie(w http.ResponseWriter, tokenString string) {
	// Create the cookie
	expire := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name: "jwt_token",
		// Token string contians user info
		Value:    tokenString,
		Expires:  expire,
		HttpOnly: true,
		Secure:   true, // Set to false if not using HTTPS
		Path:     "/",
	}

	// Set the cookie in the response header
	http.SetCookie(w, &cookie)
}
