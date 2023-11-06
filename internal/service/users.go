package service

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/email"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	FindAll(limit int, offset int, order string, conditions []string) (*models.PaginatedUsers, error)
	FindById(int) (*db.User, error)
	FindByEmail(string) (*db.User, error)
	Create(user *models.CreateUser) (*db.User, error)
	Update(int, *models.UpdateUser) (*db.User, error)
	Delete(int) error
	// Takes an email and if the email is found in the database, will reset the password and send an email to the user with the new password
	ResetPasswordAndSendEmail(email string) error
}

type userService struct {
	repo repository.UserRepository
	mail email.Email
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo, mail: email.NewSMTPEmail()}
}

// Creates a user in the database
func (s *userService) Create(user *models.CreateUser) (*db.User, error) {
	// Build hashed password from user password input
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt password: %w", err)
	}
	// Create a new user of type db User
	userToCreate := db.User{
		Username: user.Username,
		Password: string(hashedPassword),
		Name:     user.Name,
		Email:    user.Email,
	}

	// Create above user in database
	createdUser, err := s.repo.Create(&userToCreate)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}

	return createdUser, nil
}

// Find a list of users in the database
func (s *userService) FindAll(limit int, offset int, order string, conditions []string) (*models.PaginatedUsers, error) {

	users, err := s.repo.FindAll(limit, offset, order, conditions)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Find user in database by ID
func (s *userService) FindById(userId int) (*db.User, error) {
	fmt.Printf("Finding user with id: %v\n", userId)
	// Find user by id
	user, err := s.repo.FindById(userId)
	// If error detected
	if err != nil {
		return nil, err
	}
	// else
	return user, nil
}

// Find user in database by email
func (s *userService) FindByEmail(email string) (*db.User, error) {
	user, err := s.repo.FindByEmail(email)
	// If error detected
	if err != nil {
		fmt.Printf("error found in Find by email: %v", err)
		return nil, err
	}
	// else
	return user, nil
}

// Delete user in database
func (s *userService) Delete(id int) error {
	err := s.repo.Delete(id)
	// If error detected
	if err != nil {
		fmt.Println("error in deleting user: ", err)
		return err
	}
	// else
	return nil
}

// Takes an email and if the email is found in the database, will reset the password and send an email to the user with the new password
func (s *userService) ResetPasswordAndSendEmail(userEmail string) error {
	// Check if user exists in db
	foundUser, err := s.repo.FindByEmail(userEmail)
	if err != nil {
		fmt.Println("error in resetting password. User not found: ", userEmail)
		return err
	}
	// Else
	// Generate random password
	randomPassword, err := helpers.GenerateRandomString(10)
	if err != nil {
		return err
	}
	// Update found user's password
	s.repo.Update(int(foundUser.ID), &db.User{Password: randomPassword})

	// Build data for template
	data := struct {
		Name        string
		NewPassword string
	}{
		Name:        foundUser.Name,
		NewPassword: randomPassword,
	}

	// Build HTML email template from file using injected data
	emailString, err := helpers.LoadTemplate("internal/email/templates/password-reset.tmpl", data)
	if err != nil {
		fmt.Printf("error in loading template: %v", err)
		return err
	}

	// Send email with new password async (non-blocking)
	go s.mail.SendEmail(userEmail, "Password Reset Request", emailString)

	// Return no error found
	return nil
}

// Updates user in database
func (s *userService) Update(id int, user *models.UpdateUser) (*db.User, error) {
	// Create db User type of incoming DTO
	dbUser := &db.User{Name: user.Name, Username: user.Username, Email: user.Email, Password: user.Password}

	// Update using repo
	updatedUser, err := s.repo.Update(id, dbUser)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
