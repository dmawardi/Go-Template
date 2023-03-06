package services

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
)

// Creates a user in the database
func CreateUser(user *models.CreateUser) (*db.User, error) {
	// Build hashed password from user password input
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	// Create a new user of type db User
	userToCreate := db.User{
		Username: user.Username,
		Password: string(hashedPassword),
		Name:     user.Name,
		Email:    user.Email,
	}

	// Create above user in database
	result := app.DbClient.Create(&userToCreate)
	if result.Error != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}

	return &userToCreate, nil
}

// Find a list of users in the database
func FindAllUsers(limit int, offset int, order string) (*[]db.User, error) {
	// Query all users based on the received parameters
	users, err := QueryAllUsersBasedOnParams(limit, offset, order, app.DbClient)
	if err != nil {
		fmt.Printf("Error querying db for list of users: %s", err)
		return nil, err
	}

	fmt.Printf("Found user list: %v users\n", len(users))
	return &users, nil
}

// Takes limit, offset, and order parameters, builds a query and executes returning a list of users
func QueryAllUsersBasedOnParams(limit int, offset int, order string, dbClient *gorm.DB) ([]db.User, error) {
	// fmt.Printf("Params in Build query:\nlimit: %v\noffset: %v\norder: %v\n", limit, offset, order)
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

// Find user in database by ID
func FindUserById(userId int) (*db.User, error) {
	fmt.Printf("Finding user with id: %v\n", userId)
	// Create an empty ref object of type user
	user := db.User{}
	// Check if user exists in db
	result := app.DbClient.Select("ID", "name", "username", "email", "role").First(&user, userId)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in finding user: ", result.Error)
		return nil, result.Error
	}
	// else
	return &user, nil
}

// Find user in database by email
func FindUserByEmail(email string) (*db.User, error) {
	// Create an empty ref object of type user
	user := db.User{}
	// Check if user exists in db
	result := app.DbClient.Where("email = ?", email).First(&user)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in finding user: ", result.Error)
		return nil, result.Error
	}
	// else
	return &user, nil
}

// Delete user in database
func DeleteUser(id int) error {
	// Create an empty ref object of type user
	user := db.User{}
	// Check if user exists in db
	result := app.DbClient.Delete(&user, id)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in deleting user: ", result.Error)
		return result.Error
	}
	// else
	return nil
}

// Updates user in database
func UpdateUser(id *int, user *models.UpdateUser) (*db.User, error) {
	// Init
	var err error
	// Find user by id
	foundUser, err := FindUserById(*id)
	if err != nil {
		fmt.Println("User to update not found: ", err)
		return nil, err
	}
	// fmt.Printf("user to update has been found: %v", foundUser)

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
	updateResult := app.DbClient.Model(&foundUser).Updates(user)
	if updateResult.Error != nil {
		fmt.Println("User update failed: ", updateResult.Error)
		return nil, updateResult.Error
	}

	fmt.Printf("updated Username: %v\nPassword: %v", foundUser.Username, foundUser.Password)

	// Find user by id
	updatedUser, err := FindUserById(*id)
	if err != nil {
		fmt.Println("User to update not found: ", err)
		return nil, err
	}
	return updatedUser, nil
}
