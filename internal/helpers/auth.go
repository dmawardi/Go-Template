package helpers

import (
	"strings"

	"github.com/dmawardi/Go-Template/internal/models"
)

// Policy search helpers
// Searches a list of policies for a given resource based on search term
func SearchPoliciesByResource(maps []models.PolicyRuleCombinedActions, searchTerm string) []models.PolicyRuleCombinedActions {
	var result []models.PolicyRuleCombinedActions

	// Iterate through map of policies
	for _, m := range maps {

		// If success and resource contains search term
		if ContainsString(m.Resource, searchTerm) {
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

// SearchG2Records searches through all fields of each G2Record for the searchTerm and adds any match found in the results
func SearchGRecords(records []models.GRecord, searchTerm string) []models.GRecord {
	var result []models.GRecord

	// Iterate through the slice of G2Records
	for _, record := range records {
		// Check if the searchTerm is in any of the record's fields
		if strings.Contains(record.Role, searchTerm) || strings.Contains(record.InheritsFrom, searchTerm) {
			result = append(result, record)
		}
	}

	return result
}

// FilterOnlyRolesToGRecords filters out the role inheritances from a slice of roles and assignments (All G records)
func FilterOnlyInheritanceToGRecords(rolesAndAssignments [][]string) ([]models.GRecord, error) {
	var roleInheritancePolicies []models.GRecord

	for _, policy := range rolesAndAssignments {
		// Assuming policy[0] contains the role/subject and policy[1] contains the inherited role
		// Adjust the indexing based on your actual policy structure
		if strings.HasPrefix(policy[0], "role:") && strings.HasPrefix(policy[1], "role:") {
			roleInheritancePolicies = append(roleInheritancePolicies, models.GRecord{Role: policy[0], InheritsFrom: policy[1]})
		}
	}

	return roleInheritancePolicies, nil
}

func FilterOnlyRolesToList(rolesAndAssignments [][]string) []string {
	var roles []string
	// Filter out the roles that are not user assigned
	for _, policy := range rolesAndAssignments {
		if strings.HasPrefix(policy[0], "role:") {
			// If not already contained within the slice, add it
			if !ArrayContainsString(roles, policy[0]) {
				roles = append(roles, policy[0])
			}
		}
		// If inherits from is a role, add it to the slice
		if strings.HasPrefix(policy[1], "role:") {
			// If not already contained within the slice, add it
			if !ArrayContainsString(roles, policy[1]) {
				roles = append(roles, policy[1])
			}
		}
	}
	return roles
}

// Grabs a slice of GRecords and filters out the roles into a string slice
func ConvertInheritanceGRecordsToRoleList(roles []models.GRecord) []string {
	var roleList []string
	// Iterate through inheritance policies and add to roles slice
	for _, role := range roles {
		if !ArrayContainsString(roleList, role.Role) {
			roleList = append(roleList, role.Role)
		}
		if !ArrayContainsString(roleList, role.InheritsFrom) {
			roleList = append(roleList, role.InheritsFrom)
		}
	}

	return roleList
}

// Searches a list of policies for a given resource based on search term
func SearchPoliciesForExactResouceMatch(maps []models.PolicyRuleCombinedActions, searchTerm string) []models.PolicyRuleCombinedActions {
	var result []models.PolicyRuleCombinedActions

	// Iterate through map of policies
	for _, m := range maps {
		// If success and resource contains search term
		if m.Resource == searchTerm {
			result = append(result, m)
		}
	}

	return result
}

func ApplyNamingConventionToRoleInheritanceRecord(inherit *models.GRecord) {
	inherit.Role = "role:" + inherit.Role
	inherit.InheritsFrom = "role:" + inherit.InheritsFrom
}

// Function to unslugify a resource name
func UnslugifyResourceName(slugifiedResourceName string) string {
	return strings.ReplaceAll(slugifiedResourceName, "-", "/")
}

func SlugifyResourceName(resource string) string {
	return strings.ReplaceAll(resource, "/", "-")
}
