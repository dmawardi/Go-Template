package service

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/repository"
)

type PostService interface {
	FindAll(limit int, offset int, order string, conditions []string) (*models.PaginatedPosts, error)
	FindById(int) (*db.Post, error)
	Create(post *models.CreatePost) (*db.Post, error)
	Update(int, *models.UpdatePost) (*db.Post, error)
	Delete(int) error
}

type postService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) PostService {
	return &postService{repo: repo}
}

// Creates a post in the database
func (s *postService) Create(post *models.CreatePost) (*db.Post, error) {
	// Create a new user of type db User
	postToCreate := db.Post{
		Title: post.Title,
		Body:  post.Body,
		User:  post.User,
	}

	// Create above post in database
	createdPost, err := s.repo.Create(&postToCreate)
	if err != nil {
		return nil, fmt.Errorf("failed creating post: %w", err)
	}

	return createdPost, nil
}

// Find a list of posts in the database
func (s *postService) FindAll(limit int, offset int, order string, conditions []string) (*models.PaginatedPosts, error) {

	posts, err := s.repo.FindAll(limit, offset, order, conditions)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// Find post in database by ID
func (s *postService) FindById(id int) (*db.Post, error) {
	// Find post by id
	post, err := s.repo.FindById(id)
	// If error detected
	if err != nil {
		return nil, err
	}
	// else
	return post, nil
}

// Delete post in database
func (s *postService) Delete(id int) error {
	err := s.repo.Delete(id)
	// If error detected
	if err != nil {
		fmt.Println("error in deleting post: ", err)
		return err
	}
	// else
	return nil
}

// Updates post in database
func (s *postService) Update(id int, post *models.UpdatePost) (*db.Post, error) {
	// Create db Post type from incoming DTO
	postToUpdate := &db.Post{
		Title: post.Title,
		Body:  post.Body,
		User:  post.User,
	}

	// Update using repo
	updatedPost, err := s.repo.Update(id, postToUpdate)
	if err != nil {
		return nil, err
	}

	return updatedPost, nil
}
