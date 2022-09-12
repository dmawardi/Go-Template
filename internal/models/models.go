package models

type Job struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Login
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
