package helpers

import (
	"crypto/rand"
	"net/http"
	"strings"
	"unicode"

	"github.com/asaskevich/govalidator"
)

// Authentication helper functions
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

// Capitalizes the first letter of a string
func CapitalizeFirstLetter(str string) string {
	for _, r := range str {
		u := string(unicode.ToUpper(r))
		return u + str[len(u):]
	}
	return ""
}
