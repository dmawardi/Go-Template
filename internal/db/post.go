package db

import (
	"fmt"
	"time"
)

// DB Schema interface implementation
// Mapping of field names to values to allow for dynamic access
func (schemaObject Post) ObtainValue(keyValue string, fieldMap map[string]string) string {
	// Map of post fields
	fieldMap2 := map[string]string{
		"ID":        fmt.Sprint(schemaObject.ID),
		"CreatedAt": schemaObject.CreatedAt.Format(time.RFC3339),
		"UpdatedAt": schemaObject.UpdatedAt.Format(time.RFC3339),
		"Title":     schemaObject.Title,
		"Body":      schemaObject.Body,
		"UserID":    fmt.Sprint(schemaObject.UserID),
	}
	// Return value of key
	return fieldMap2[keyValue]
}

// Grabs the ID of the schema object as string
func (schemaObject Post) GetID() string {
	return fmt.Sprint(schemaObject.ID)
}
