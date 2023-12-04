package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Interface for all schemas (used for Admin panel)
type AdminPanelSchema interface {
	GetID() string
	ObtainValue(keyValue string) string
}

// Schemas
type User struct {
	// gorm.Model `json:"-"`
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `swaggertype:"string" json:"created_at,omitempty"`
	UpdatedAt time.Time      `swaggertype:"string" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `json:"name,omitempty"`
	Username  string         `json:"username,omitempty"`
	Email     string         `json:"email,omitempty" gorm:"uniqueIndex"`
	Password  string         `json:"-"`
	Role      string         `json:"role,omitempty" gorm:"default:user"`
	// Verification
	Verified               *bool     `json:"verified,omitempty" gorm:"default:false"`
	VerificationCode       string    `json:"verification_code,omitempty" gorm:"default:null"`
	VerificationCodeExpiry time.Time `json:"verification_code_expiry,omitempty" gorm:"default:null"`
	// Relationships
	Posts []Post `json:"posts,omitempty" gorm:"foreignKey:UserID"`
}

// DB Schema interface implementation
// Mapping of field names to values to allow for dynamic access
func (schemaObject User) ObtainValue(keyValue string) string {
	// Map of user fields
	fieldMap := map[string]string{
		"ID":                     fmt.Sprint(schemaObject.ID),
		"CreatedAt":              schemaObject.CreatedAt.Format(time.RFC3339),
		"UpdatedAt":              schemaObject.UpdatedAt.Format(time.RFC3339),
		"Name":                   schemaObject.Name,
		"Username":               schemaObject.Username,
		"Email":                  schemaObject.Email,
		"Role":                   schemaObject.Role,
		"Verified":               fmt.Sprint(schemaObject.Verified),
		"VerificationCode":       schemaObject.VerificationCode,
		"VerificationCodeExpiry": schemaObject.VerificationCodeExpiry.Format(time.RFC3339),
	}
	// Return value of key
	return fieldMap[keyValue]
}

// Grabs the ID of the schema object as string
func (schemaObject User) GetID() string {
	return fmt.Sprint(schemaObject.ID)
}

type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `swaggertype:"string" json:"created_at,omitempty"`
	UpdatedAt time.Time      `swaggertype:"string" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Title     string         `json:"title,omitempty"`
	Body      string         `json:"body,omitempty"`
	UserID    uint           `json:"user_id,omitempty"`
	User      User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// DB Schema interface implementation
// Mapping of field names to values to allow for dynamic access
func (schemaObject Post) ObtainValue(keyValue string) string {
	// Map of post fields
	fieldMap := map[string]string{
		"ID":        fmt.Sprint(schemaObject.ID),
		"CreatedAt": schemaObject.CreatedAt.Format(time.RFC3339),
		"UpdatedAt": schemaObject.UpdatedAt.Format(time.RFC3339),
		"Title":     schemaObject.Title,
		"Body":      schemaObject.Body,
		"UserID":    fmt.Sprint(schemaObject.UserID),
	}
	// Return value of key
	return fieldMap[keyValue]
}

// Grabs the ID of the schema object as string
func (schemaObject Post) GetID() string {
	return fmt.Sprint(schemaObject.ID)
}
