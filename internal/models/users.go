package models

import (
	"time"

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
}

// Update User structure for Data transfer.
type UpdateUser struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty" valid:"length(6|25)"`
	Password string `json:"password,omitempty" valid:"length(6|30)"`
	Name     string `json:"name,omitempty" valid:"length(6|80)"`
	Email    string `json:"email,omitempty" valid:"email"`
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
