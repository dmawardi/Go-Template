package models

// Interface for all schemas that makes it compatible with admin panel (Add receiver functions for every schema)
type AdminPanelSchema interface {
	// Returns ID of record
	GetID() string
	// Returns value of schema field
	ObtainValue(keyValue string) string
}