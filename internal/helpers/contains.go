package helpers

import "strings"

// Checks if a string contains another string (Used to search for policy resource)
func ContainsString(s, searchTerm string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(searchTerm))
}

// Checks if array contains a particular string value (Used in policy)
func ArrayContainsString(arr []string, target string) bool {
	for _, s := range arr {
		if s == target {
			return true
		}
	}
	return false
}
