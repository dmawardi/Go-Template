package models

type Job struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Login
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ValidationError struct {
	Validation_errors map[string][]string `json:"validation_errors"`
}
