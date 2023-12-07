package models

import (
	"time"

	"github.com/dmawardi/Go-Template/internal/db"
	"gorm.io/gorm"
)

type LoginResponse struct {
	Token string `json:"token"`
}

// Users
// Create User structure for Data transfer.
type CreateUser struct {
	Username string `json:"username" valid:"length(6|25),required"`
	Password string `json:"password" valid:"length(6|30),required"`
	Name     string `json:"name" valid:"length(6|80),required"`
	Email    string `json:"email" valid:"email,required"`
	Verified bool   `json:"verified,omitempty"`
	Role     string `json:"role,omitempty" valid:"required,in(admin|moderator|user)"`
}

// Created user (for admin use)
type CreatedUser struct {
	ID        uint           `json:"id"`
	Username  string         `json:"username"`
	Password  string         `json:"password"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Role      string         `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

// The user sent to users
type PartialUser struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Verified bool   `json:"verified"`
}

// Update User structure for Data transfer.
type UpdateUser struct {
	Username string `json:"username,omitempty" valid:"length(6|25)"`
	Password string `json:"password,omitempty" valid:"length(6|30)"`
	Name     string `json:"name,omitempty" valid:"length(6|80)"`
	Email    string `json:"email,omitempty" valid:"email"`
	Verified bool   `json:"verified,omitempty"`
	Role     string `json:"role,omitempty" valid:"in(admin|moderator|user)"`
}
type UpdatedUser struct {
	ID        uint           `json:"id"`
	Username  string         `json:"username"`
	Password  string         `json:"password"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type ResetPasswordAndEmailVerification struct {
	Email string `json:"email" valid:"email,required"`
}

type PaginatedUsers struct {
	Data *[]db.User     `json:"data"`
	Meta SchemaMetaData `json:"meta"`
}

// Used to init the query params for easy extraction in controller
// Returns: map[string]string{"age": "int", "name": "string", "active": "bool"}
func UserQueryParams() map[string]string {
	return map[string]string{
		"email":    "string",
		"name":     "string",
		"username": "string",
		"verified": "bool",
		"role":     "string",
	}
}
