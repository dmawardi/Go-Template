package models

import (
	"time"
)

type Job struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Login
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Users

// Create User structure for Data transfer.
type CreateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

// Update User structure for Data transfer.
type UpdateUser struct {
	Id         int       `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
