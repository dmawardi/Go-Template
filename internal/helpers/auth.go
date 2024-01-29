package helpers

import "strings"

// Policy search helpers
// Searches a list of policies for a given resource based on search term
func SearchPoliciesByResource(maps []map[string]interface{}, searchTerm string) []map[string]interface{} {
	var result []map[string]interface{}

	// Iterate through map of policies
	for _, m := range maps {
		// Grab resource
		resource, ok := m["resource"].(string)
		// If success and resource contains search term
		if ok && ContainsString(resource, searchTerm) {
			result = append(result, m)
		}
	}

	return result
}

// Searches a list of maps for a given key based on search term
func SearchMapKeysFor(maps []map[string]string, mapKeysToSearch []string, searchTerm string) []map[string]string {
	var result []map[string]string
	// Init to record if already added to results
	addedToResult := false

	// Iterate through map of policies
	for _, m := range maps {
		// Reset added to result
		addedToResult = false
		// Iterate through list of keys to search for term
		for _, keyToSearch := range mapKeysToSearch {
			// Grab value
			value, ok := m[keyToSearch]
			// If success, and the record hasn't been added already and value contains search term
			if ok && ContainsString(value, searchTerm) && !addedToResult {
				// Append
				result = append(result, m)
				// Set added to true
				addedToResult = true
			}
		}
	}

	return result
}

// Searches a list of policies for a given resource based on search term
func SearchPoliciesForExactResouceMatch(maps []map[string]interface{}, searchTerm string) []map[string]interface{} {
	var result []map[string]interface{}

	// Iterate through map of policies
	for _, m := range maps {
		// Grab resource
		resource, ok := m["resource"].(string)
		// If success and resource contains search term
		if ok && resource == searchTerm {
			result = append(result, m)
		}
	}

	return result
}

// Function to unslugify a resource name
func UnslugifyResourceName(slugifiedResourceName string) string {
	return strings.ReplaceAll(slugifiedResourceName, "-", "/")
}
