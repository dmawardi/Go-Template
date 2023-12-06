package repository

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	// Find a list of all users in the Database
	FindAll(limit int, offset int, order string, conditions []string) (*models.PaginatedPosts, error)
	FindById(int) (*db.Post, error)
	Create(post *db.Post) (*db.Post, error)
	Update(int, *db.Post) (*db.Post, error)
	Delete(int) error
}

type postRepository struct {
	DB *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db}
}

// Creates a post in the database
func (r *postRepository) Create(post *db.Post) (*db.Post, error) {
	// Create above post in database
	result := r.DB.Create(&post)
	if result.Error != nil {
		return nil, fmt.Errorf("failed creating post: %w", result.Error)
	}

	return post, nil
}

// Find a list of posts in the database
func (r *postRepository) FindAll(limit int, offset int, order string, conditions []string) (*models.PaginatedPosts, error) {
	// Fetch metadata from database
	var totalCount *int64

	// Count the total number of records
	totalCount, err := db.CountBasedOnConditions(db.Post{}, conditions, r.DB)
	if err != nil {
		return nil, err
	}
	// Calculate next page
	nextPage := offset + limit
	// If next page is greater than total count, set to 0
	if nextPage > int(*totalCount) {
		nextPage = 0
	}
	prevPage := offset - limit
	// If prev page is less than 0, set to 0
	if prevPage < 0 {
		prevPage = 0
	}

	// Build metadata object
	metaData := models.NewSchemaMetaData(*totalCount, limit, offset, &nextPage, &prevPage)
	// Query all post based on the received parameters
	posts, err := QueryAllPostsBasedOnParams(limit, offset, order, conditions, r.DB)
	if err != nil {
		fmt.Printf("Error querying db for list of posts: %s", err)
		return nil, err
	}

	return &models.PaginatedPosts{
		Data: &posts,
		Meta: metaData,
	}, nil
}

// Find post in database by ID
func (r *postRepository) FindById(id int) (*db.Post, error) {
	// Create an empty ref object of type post
	post := db.Post{}
	// Check if post exists in db
	result := r.DB.First(&post, id)

	// If error detected
	if result.Error != nil {
		return nil, result.Error
	}
	// else
	return &post, nil
}

// Delete post in database
func (r *postRepository) Delete(id int) error {
	// Create an empty ref object of type post
	post := db.Post{}
	// Check if post exists in db
	result := r.DB.Delete(&post, id)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in deleting post: ", result.Error)
		return result.Error
	}
	// else
	return nil
}

// Updates post in database
func (r *postRepository) Update(id int, post *db.Post) (*db.Post, error) {
	// Init
	var err error
	// Find post by id
	foundPost, err := r.FindById(id)
	if err != nil {
		fmt.Println("Post to update not found: ", err)
		return nil, err
	}

	// Update post using found post
	updateResult := r.DB.Model(&foundPost).Updates(post)
	if updateResult.Error != nil {
		fmt.Println("Post update failed: ", updateResult.Error)
		return nil, updateResult.Error
	}

	// Retrieve changed post by id
	updatedPost, err := r.FindById(id)
	if err != nil {
		fmt.Println("Post to update not found: ", err)
		return nil, err
	}
	return updatedPost, nil
}

// Takes limit, offset, and order parameters, builds a query and executes returning a list of posts
func QueryAllPostsBasedOnParams(limit int, offset int, order string, conditions []string, dbClient *gorm.DB) ([]db.Post, error) {
	// Build model to query database
	posts := []db.Post{}
	// Build base query for posts table
	query := dbClient.Model(&posts)

	// Add parameters into query as needed
	if limit != 0 {
		query.Limit(limit)
	}
	if offset != 0 {
		query.Offset(offset)
	}
	// order format should be "column_name ASC/DESC" eg. "created_at ASC"
	if order != "" {
		query.Order(order)
	}
	// Add conditions to query
	if len(conditions) > 0 {
		for _, condition := range conditions {
			// Add condition to query
			query.Where(condition)
		}
	}
	// Query database
	result := query.Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	// Return if no errors with result
	return posts, nil
}
