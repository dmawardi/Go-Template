package helpers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// GenerateChangeLog compares two objects of the same type and returns a JSON string representing the differences.
func GenerateChangeLog(changeLog ChangeLogInput) (string, error) {
	// Use reflection to get the values of the old and new objects
	oldVal := reflect.ValueOf(changeLog.OldObj)
	newVal := reflect.ValueOf(changeLog.NewObj)
	
	// Check if the types of the two objects are the same
	if oldVal.Type() != newVal.Type() {
		return "", fmt.Errorf("type mismatch: %v vs %v", oldVal.Type(), newVal.Type())
	}

	// Initialize a map to hold the change log
	changeLogMap := make(map[string]map[string]interface{})
	// Get the type of the old object (since both are of the same type, this is sufficient)
	oldType := oldVal.Type()

	// Check if the old object is empty (indicating a creation)
	isOldObjEmpty := true
	for i := 0; i < oldVal.NumField(); i++ {
		if !reflect.DeepEqual(oldVal.Field(i).Interface(), reflect.Zero(oldVal.Field(i).Type()).Interface()) {
			isOldObjEmpty = false
			break
		}
	}

	// Iterate over the fields of the objects
	for i := 0; i < oldVal.NumField(); i++ {
		// Get the field name
		fieldName := oldType.Field(i).Name
		// Get the values of the field for both objects
		oldFieldValue := oldVal.Field(i).Interface()
		newFieldValue := newVal.Field(i).Interface()

		// Compare the field values; if they differ, or if it's a creation, add them to the change log
		if isOldObjEmpty || !reflect.DeepEqual(oldFieldValue, newFieldValue) {
			changeLogMap[fieldName] = map[string]interface{}{
				"old": oldFieldValue,
				"new": newFieldValue,
			}
		}
	}

	// Marshal the change log map to a JSON string
	changeLogJSON, err := json.Marshal(changeLogMap)
	if err != nil {
		return "", err
	}

	return string(changeLogJSON), nil
}
// GenerateChangeDescription generates a human-readable description of the changes
func GenerateChangeDescription(changeLogJSON string, entityType string, actionType string) (string, error) {
	// If the entity was deleted, return a deletion description
	if actionType == "delete" {
		return fmt.Sprintf("Deleted the %s.", entityType), nil
	}

	// Initialize a map to hold the change log data
	var changeLogMap map[string]map[string]interface{}
	
	// Unmarshal the JSON change log into the map
	err := json.Unmarshal([]byte(changeLogJSON), &changeLogMap)
	if err != nil {
		return "", err
	}

	// If the change log is empty, return a message indicating no changes were detected
	if len(changeLogMap) == 0 {
		return fmt.Sprintf("No changes detected for the %s.", entityType), nil
	}

	// Prepare a slice to hold the field names that were changed
	var changes []string
	for field := range changeLogMap {
		changes = append(changes, fmt.Sprintf("%s", strings.ToLower(field)))
	}

	// If there is only one change, specify that field in the description
	if len(changes) == 1 {
		return fmt.Sprintf("Updated the %s %s.", entityType, changes[0]), nil
	}

	// If there are multiple changes, join them with "and" and specify in the description
	return fmt.Sprintf("Updated the %s %s.", entityType, strings.Join(changes, " and ")), nil
}
type ChangeLogInput struct {
	OldObj interface{}
	NewObj interface{}
}