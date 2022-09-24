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

// Created user
type CreatedUser struct {
	Id         string      `json:"id"`
	Username   string      `json:"username"`
	Password   string      `json:"password"`
	Name       string      `json:"name"`
	Email      string      `json:"email"`
	Role       string      `json:"role"`
	Created_at time.Time   `json:"created_at"`
	Updated_at time.Time   `json:"updated_at"`
	Edges      interface{} `json:"edges"`
}

// Update User structure for Data transfer.
type UpdateUser struct {
	Id       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
}
type UpdatedUser struct {
	Id         int       `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
