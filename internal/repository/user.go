package repository

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	// Find a list of all users in the Database
	FindAll(limit int, offset int, order string, conditions []interface{}) (*models.PaginatedUsers, error)
	Create(user *db.User) (*db.User, error)
	Update(int, *db.User) (*db.User, error)
	Delete(int) error
	BulkDelete([]int) error
	// Find
	FindById(int) (*db.User, error)
	FindByEmail(string) (*db.User, error)
	// Verification
	FindByVerificationCode(string) (*db.User, error)
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
func (r *userRepository) FindAll(limit int, offset int, order string, conditions []interface{}) (*models.PaginatedUsers, error) {
	// Build meta data for posts
	metaData, err := models.BuildMetaData(r.DB, db.Post{}, limit, offset, order, conditions)
	if err != nil {
		fmt.Printf("Error building meta data: %s", err)
		return nil, err
	}

	// Query all users based on the received parameters
	var users []db.User
	err = db.QueryAll(r.DB, &users, limit, offset, order, conditions)
	if err != nil {
		fmt.Printf("Error querying db for list of users: %s", err)
		return nil, err
	}

	return &models.PaginatedUsers{
		Data: &users,
		Meta: *metaData,
	}, nil
}

// Find user in database by ID
func (r *userRepository) FindById(userId int) (*db.User, error) {
	// Create an empty ref object of type user
	user := db.User{}
	// Check if user exists in db
	result := r.DB.Select("ID", "name", "username", "email", "role", "verified").First(&user, userId)

	// If error detected
	if result.Error != nil {
		return nil, result.Error
	}

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

// Bulk delete users in database
func (r *userRepository) BulkDelete(ids []int) error {
	// Delete users with specified IDs
	err := db.BulkDeleteByIds(db.User{}, ids, r.DB)
	if err != nil {
		fmt.Println("error in deleting users: ", err)
		return err
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

	fmt.Printf("Updating this user: %v with user: %v\n", *foundUser.Verified, *user.Verified)

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
	fmt.Printf("Updated user: %v\n", *updatedUser.Verified)
	return updatedUser, nil
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

// Find a user by the verification code associated with the user
func (r *userRepository) FindByVerificationCode(token string) (*db.User, error) {
	// Create an empty ref object of type user
	user := db.User{}
	// Check if user exists in db
	result := r.DB.Where("verification_code = ?", token).First(&user)

	// If error detected
	if result.Error != nil {
		return nil, result.Error
	}
	// else
	return &user, nil
}
