package models

type Job struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Login
type Login struct {
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"required"`
}

type ValidationError struct {
	Validation_errors map[string][]string `json:"validation_errors"`
}
