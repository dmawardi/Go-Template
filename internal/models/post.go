package models

import "github.com/dmawardi/Go-Template/internal/db"

type CreatePost struct {
	Title string  `json:"title,omitempty" valid:"length(3|36),required"`
	Body  string  `json:"body,omitempty" valid:"length(10|),required"`
	User  db.User `json:"user,omitempty" valid:"required"`
}

type UpdatePost struct {
	Title string  `json:"title" valid:"length(3|36),required"`
	Body  string  `json:"body" valid:"length(10|),required"`
	User  db.User `json:"user" valid:"required"`
}

type PaginatedPosts struct {
	Data *[]db.Post     `json:"data"`
	Meta SchemaMetaData `json:"meta"`
}
