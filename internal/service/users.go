package service

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/repository"
)

type UserService interface {
	FindAll(int, int, string) (*[]db.User, error)
	FindById(int) (*db.User, error)
	FindByEmail(string) (*db.User, error)
	Create(user *models.CreateUser) (*db.User, error)
	Update(int, *models.UpdateUser) (*db.User, error)
	Delete(int) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

// Creates a user in the database
func (s *userService) Create(user *models.CreateUser) (*db.User, error) {
	// Create above user in database
	createdUser, err := s.repo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}

	return createdUser, nil
}

// Find a list of users in the database
func (s *userService) FindAll(limit int, offset int, order string) (*[]db.User, error) {

	users, err := s.repo.FindAll(limit, offset, order)
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
		fmt.Println("error in finding user: ", err)
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
		fmt.Println("error in finding user: ", err)
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

// Updates user in database
func (s *userService) Update(id int, user *models.UpdateUser) (*db.User, error) {
	updatedUser, err := s.repo.Update(id, user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
