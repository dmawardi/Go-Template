package service

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/email"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	FindAll(limit int, offset int, order string, conditions []models.QueryConditionParameters) (*models.PaginatedUsersWithRole, error)
	FindById(int) (*models.UserWithRole, error)
	FindByEmail(string) (*models.UserWithRole, error)
	Create(user *models.CreateUser) (*models.UserWithRole, error)
	Update(int, *models.UpdateUser) (*models.UserWithRole, error)
	Delete(int) error
	BulkDelete([]int) error
	CheckPasswordMatch(id int, password []byte) bool
	// Login
	LoginUser(login *models.Login) (string, error)
	// Takes an email and if the email is found in the database, will reset the password and send an email to the user with the new password
	ResetPasswordAndSendEmail(email string) error
	// Verifies user email in database
	VerifyEmailCode(token string) error
	// Sends verification email for user
	ResendVerificationEmail(id int) error
}

type userService struct {
	repo repository.UserRepository
	auth repository.AuthPolicyRepository
	mail email.Email
}

// Builds a new service with injected repository. Includes email service
func NewUserService(repo repository.UserRepository, auth repository.AuthPolicyRepository, mail email.Email) UserService {
	return &userService{repo: repo, auth: auth, mail: mail}
}

// Creates a user in the database
func (s *userService) Create(user *models.CreateUser) (*models.UserWithRole, error) {
	// Build hashed password from user password input
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt password: %w", err)
	}
	// Create a new user of type db User
	toCreate := db.User{
		Username: user.Username,
		Password: string(hashedPassword),
		Name:     user.Name,
		Email:    user.Email,
		Verified: &user.Verified,
	}

	// Create above user in database
	created, err := s.repo.Create(&toCreate)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}

	// If no role found, assign user role
	if user.Role == "" {
		user.Role = "role:user"
	}
	// Assign user role
	success, err := s.auth.AssignUserRole(fmt.Sprint(created.ID), user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed assigning user role: %w", err)
	}
	if !*success {
		return nil, fmt.Errorf("failed assigning user role: %w", err)
	}

	// Combine user and role data
	userToReturn := BuildUserWithRole(created, user.Role)

	return userToReturn, nil
}

// Find a list of users in the database
func (s *userService) FindAll(limit int, offset int, order string, conditions []models.QueryConditionParameters) (*models.PaginatedUsersWithRole, error) {
	// Query all users based on the received parameters
	users, err := s.repo.FindAll(limit, offset, order, conditions)
	if err != nil {
		return nil, err
	}

	// Init full users slice
	var fullUsers []models.UserWithRole
	// Iterate through users and attach role to complete the data
	for _, user := range *users.Data {
		// Get user role and attach to user
		fullUser, err := findRoleAndAttach(&user, s.auth)
		if err != nil {
			return nil, err
		}
		// Append to full users slice
		fullUsers = append(fullUsers, *fullUser)
	}

	return &models.PaginatedUsersWithRole{Data: &fullUsers, Meta: users.Meta}, nil
}

// Find user in database by ID
func (s *userService) FindById(userId int) (*models.UserWithRole, error) {
	// Find user by id
	user, err := s.repo.FindById(userId)
	// If error detected
	if err != nil {
		return nil, err
	}

	// Get user role and attach to user
	fullUser, err := findRoleAndAttach(user, s.auth)
	if err != nil {
		return nil, err
	}
	// else
	return fullUser, nil
}

// Find user in database by email
func (s *userService) FindByEmail(email string) (*models.UserWithRole, error) {
	user, err := s.repo.FindByEmail(email)
	// If error detected
	if err != nil {
		fmt.Printf("error found in Find by email: %v", err)
		return nil, err
	}
	// Get user role and attach to user
	fullUser, err := findRoleAndAttach(user, s.auth)
	if err != nil {
		return nil, err
	}
	// else
	return fullUser, nil
}

// Delete user in database
func (s *userService) Delete(id int) error {
	err := s.repo.Delete(id)
	// If error detected
	if err != nil {
		fmt.Println("error in deleting user: ", err)
		return err
	}

	// Delete all user roles
	success, err := s.auth.DeleteRolesForUser(fmt.Sprint(id))
	if err != nil {
		fmt.Printf("error in deleting user roles: %v\n", err)
		return err
	}
	if !*success {
		fmt.Printf("error in deleting user roles (no roles assigned?): %v", err)
		return err
	}

	// else
	return nil
}

// Deletes multiple users in database
func (s *userService) BulkDelete(ids []int) error {
	err := s.repo.BulkDelete(ids)
	// If error detected
	if err != nil {
		fmt.Println("error in bulk deleting users: ", err)
		return err
	}
	// Iterate through ids and delete all user roles
	for _, id := range ids {
		// Delete all user roles
		success, err := s.auth.DeleteRolesForUser(fmt.Sprint(id))
		if err != nil {
			fmt.Printf("error in deleting user roles: %v\n", err)
			return err
		}
		// If not successful in deletion
		if !*success {
			fmt.Printf("error in deleting user roles (no roles assigned?): %v\n", err)
			return err
		}
	}
	// else
	return nil
}

// Updates user in database
func (s *userService) Update(id int, user *models.UpdateUser) (*models.UserWithRole, error) {
	// Create db User type from incoming DTO
	toUpdate := &db.User{Name: user.Name, Username: user.Username, Email: user.Email, Password: user.Password, Verified: &user.Verified}

	// Update using repo
	updated, err := s.repo.Update(id, toUpdate)
	if err != nil {
		return nil, err
	}

	// Assign user role (if role update found)
	if user.Role != "" {
		// Update user role in policy table
		success, err := s.auth.AssignUserRole(fmt.Sprint(updated.ID), user.Role)
		if err != nil {
			return nil, fmt.Errorf("failed assigning user role: %w", err)
		}
		if !*success {
			return nil, fmt.Errorf("failed assigning user role: %w", err)
		}
	}

	// Get user role and attach to user
	fullUser, err := findRoleAndAttach(updated, s.auth)
	if err != nil {
		return nil, err
	}

	return fullUser, nil
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
	emailString, err := helpers.LoadTemplate(helpers.BuildPathFromWorkingDirectory("internal/email/templates/password-reset.tmpl"), data)
	if err != nil {
		fmt.Printf("error in loading template: %v", err)
		return err
	}

	// Send email with new password async (non-blocking)
	go s.mail.SendEmail(userEmail, "Password Reset Request", emailString)

	// Return no error found
	return nil
}

func (s *userService) LoginUser(login *models.Login) (token string, err error) {
	// Init token string
	tokenString := ""

	// Find user by email
	found, err := s.FindByEmail(login.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// If user is found
	// Compare stored (hashed) password with input password
	err = bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(login.Password))
	if err != nil {
		return "", errors.New("incorrect username/password")
	}

	// If match found provide the token string for the user
	if err == nil {
		fmt.Println("User logging in: ", found.Email)
		// Set login status to true
		tokenString, err = auth.GenerateJWT(int(found.ID), found.Email, found.Role)
		if err != nil {
			fmt.Println("Failed to create JWT")
		}
	}
	return tokenString, nil
}

func (s *userService) CheckPasswordMatch(id int, password []byte) bool {
	// Find user by id
	user, err := s.repo.FindById(id)
	if err != nil {
		fmt.Println("error in finding user: ", err)
		return false
	}

	// Compare stored (hashed) password with input password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), password)
	if err != nil {
		fmt.Println("error in comparing passwords: ", err)
		return false
	}
	// else
	return true
}

// Verifies user email in database
func (s *userService) VerifyEmailCode(token string) error {
	// Find user by token
	user, err := s.repo.FindByVerificationCode(token)
	if err != nil {
		return err
	}

	// Is the verification code's expiry before now?
	if user.VerificationCodeExpiry.Before(time.Now()) {
		// If so, it's expired
		return fmt.Errorf("verification code expired")
	}

	// else
	// Update user
	// Set verified to true
	trueVerified := true
	user.Verified = &trueVerified
	// Set verification code to empty string
	user.VerificationCode = ""

	// Update user in database
	_, err = s.repo.Update(int(user.ID), user)
	if err != nil {
		return err
	}

	// Return no error found
	return nil
}

// Sends verification email for user
func (s *userService) ResendVerificationEmail(id int) error {
	// Find user by id
	user, err := s.repo.FindById(id)
	if err != nil {
		return err
	}
	// Generate verification code and set expiry and store in user update
	userUpdate, err := helpers.GenerateVerificationCodeAndSetExpiry()
	if err != nil {
		return err
	}

	// Update user in database
	_, err = s.repo.Update(int(user.ID), userUpdate)
	if err != nil {
		return err
	}
	// Build data for email template (SERVER_PORT prefixed with :)
	baseUrl := fmt.Sprintf("%s%s", os.Getenv("SERVER_BASE_URL"), os.Getenv("SERVER_PORT"))
	// Build URL for verification
	tokenUrl := template.URL("http://" + baseUrl + "/api/users/verify-email/" + user.VerificationCode)
	data := struct {
		Name     string
		TokenUrl template.URL
	}{
		Name:     user.Name,
		TokenUrl: tokenUrl,
	}

	// Build HTML email template from file using injected data
	emailString, err := helpers.LoadTemplate(helpers.BuildPathFromWorkingDirectory("/internal/email/templates/email-verification.tmpl"), data)
	if err != nil {
		fmt.Printf("error in loading template: %v", err)
		return err
	}

	// Send email with verification async (non-blocking)
	go s.mail.SendEmail(user.Email, "Please Verify your Email", emailString)

	// Return no error found
	return nil
}

// Helper function to find user role and attach to user
func findRoleAndAttach(user *db.User, auth repository.AuthPolicyRepository) (*models.UserWithRole, error) {
	fullUser := &models.UserWithRole{}
	// Get user role
	role, err := auth.FindRoleByUserId(fmt.Sprint(user.ID))
	// If no role found, return user without role
	if err != nil {
		// Give empty value for role
		fullUser = BuildUserWithRole(user, "")
		// Ignore error (no role found)
		return fullUser, nil
	}
	// Else
	fullUser = BuildUserWithRole(user, role)

	// else
	return fullUser, nil
}

// Builds new models.UserWithRole object from db user and
func BuildUserWithRole(user *db.User, role string) *models.UserWithRole {
	return &models.UserWithRole{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Name:     user.Name,
		Email:    user.Email,
		// Authorization
		Role: role,
		// Verification
		Verified:               user.Verified,
		VerificationCode:       user.VerificationCode,
		VerificationCodeExpiry: user.VerificationCodeExpiry,
		// Timestamps
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}
