package models

import "github.com/dmawardi/Go-Template/internal/db"

type CreatePost struct {
	Title string  `json:"title,omitempty" valid:"required"`
	Body  string  `json:"body,omitempty" valid:"required"`
	User  db.User `json:"user,omitempty" valid:"required"`
}

type UpdatePost struct {
	Title string  `json:"title"`
	Body  string  `json:"body"`
	User  db.User `json:"user" valid:"required"`
}

type PaginatedPosts struct {
	Data *[]db.Post     `json:"data"`
	Meta SchemaMetaData `json:"meta"`
}
