package repository

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	// Find a list of all users in the Database
	FindAll(limit int, offset int, order string, conditions []string) (*models.PaginatedUsers, error)
	FindById(int) (*db.User, error)
	FindByEmail(string) (*db.User, error)
	Create(user *db.User) (*db.User, error)
	Update(int, *db.User) (*db.User, error)
	Delete(int) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

// Creates a user in the database
func (r *userRepository) Create(user *db.User) (*db.User, error) {
	// Create above user in database
	result := r.DB.Create(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed creating user: %w", result.Error)
	}

	return user, nil
}

// Find a list of users in the database
func (r *userRepository) FindAll(limit int, offset int, order string, conditions []string) (*models.PaginatedUsers, error) {
	// Fetch metadata from database
	var totalCount *int64

	// Count the total number of records
	totalCount, err := db.CountBasedOnConditions(db.User{}, conditions, r.DB)
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
	metaData := models.SchemaMetaData{
		Total_Records:    *totalCount,
		Records_Per_Page: limit,
		Current_Page:     offset,
		Next_Page:        helpers.ReturnNilIfZero(nextPage),
		Prev_Page:        helpers.ReturnNilIfZero(prevPage),
	}
	// Query all users based on the received parameters
	users, err := QueryAllUsersBasedOnParams(limit, offset, order, r.DB)
	if err != nil {
		fmt.Printf("Error querying db for list of users: %s", err)
		return nil, err
	}

	return &models.PaginatedUsers{
		Data: &users,
		Meta: metaData,
	}, nil
}

// Find user in database by ID
func (r *userRepository) FindById(userId int) (*db.User, error) {
	// Create an empty ref object of type user
	user := db.User{}
	// Check if user exists in db
	result := r.DB.Select("ID", "name", "username", "email", "role").First(&user, userId)

	// If error detected
	if result.Error != nil {
		return nil, result.Error
	}
	// else
	return &user, nil
}

// Find user in database by email
func (r *userRepository) FindByEmail(email string) (*db.User, error) {
	// Create an empty ref object of type user
	user := db.User{}
	// Check if user exists in db
	result := r.DB.Where("email = ?", email).First(&user)

	// If error detected
	if result.Error != nil {
		return nil, result.Error
	}
	// else
	return &user, nil
}

// Delete user in database
func (r *userRepository) Delete(id int) error {
	// Create an empty ref object of type user
	user := db.User{}
	// Check if user exists in db
	result := r.DB.Delete(&user, id)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in deleting user: ", result.Error)
		return result.Error
	}
	// else
	return nil
}

// Updates user in database
func (r *userRepository) Update(id int, user *db.User) (*db.User, error) {
	// Init
	var err error
	// Find user by id
	foundUser, err := r.FindById(id)
	if err != nil {
		fmt.Println("User to update not found: ", err)
		return nil, err
	}

	// If password from update object is not empty, use bcrypt to encrypt
	if user.Password != "" {
		// Build hashed password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			return nil, err
		}
		// Save in user update object
		user.Password = string(hashedPassword)
	}

	// Update user using found user
	updateResult := r.DB.Model(&foundUser).Updates(user)
	if updateResult.Error != nil {
		fmt.Println("User update failed: ", updateResult.Error)
		return nil, updateResult.Error
	}

	// Retrieve changed user by id
	updatedUser, err := r.FindById(id)
	if err != nil {
		fmt.Println("User to update not found: ", err)
		return nil, err
	}
	return updatedUser, nil
}

// Takes limit, offset, and order parameters, builds a query and executes returning a list of users
func QueryAllUsersBasedOnParams(limit int, offset int, order string, dbClient *gorm.DB) ([]db.User, error) {
	// Build model to query database
	users := []db.User{}
	// Build base query for users table
	query := dbClient.Model(&users)

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
	// Query database
	result := query.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	// Return if no errors with result
	return users, nil
}

func CountUsersBasedOnParams(conditions []string, dbClient *gorm.DB) (*int64, error) {
	// Fetch metadata from database
	var totalCount int64

	// Count the total number of records
	query := dbClient.Model(&db.User{})

	// Add conditions to query
	if len(conditions) > 0 {
		for _, condition := range conditions {
			// Add condition to query
			query.Where(condition)
		}

	}

	// Execute query
	countResult := query.Count(&totalCount)
	if countResult.Error != nil {
		return nil, countResult.Error
	}
	return &totalCount, nil
}
