package services

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/dmawardi/Go-Template/ent"
	"github.com/dmawardi/Go-Template/ent/user"
	"github.com/dmawardi/Go-Template/internal/models"
)

func CreateUser(ctx context.Context, client *ent.Client, user *models.CreateUser) (*ent.User, error) {
	// Build hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	u, err := client.User.
		Create().
		SetName(user.Name).
		SetUsername(user.Username).
		SetEmail(user.Email).
		SetPassword(string(hashedPassword)).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}

func FindUserById(ctx context.Context, client *ent.Client, userId int) (*ent.User, error) {
	// Check if user exists in db
	foundUser, err := app.DbClient.User.
		Query().
		Where(user.ID(userId)).
		// `Only` fails if no user found,
		// or more than 1 user returned.
		Only(app.Ctx)

	fmt.Println("Found user:", foundUser)

	// If error detected
	if err != nil {
		fmt.Println("error in finding user: ", err)
		return nil, err
	}
	// else
	return foundUser, nil
}

func FindUserByEmail(ctx context.Context, client *ent.Client, email string) (*ent.User, error) {
	// Check if user exists in db
	foundUser, err := app.DbClient.User.
		Query().
		Where(user.Email(email)).
		// `Only` fails if no user found,
		// or more than 1 user returned.
		Only(app.Ctx)

	fmt.Println("Found user:", foundUser)

	// If error detected
	if err != nil {
		fmt.Println("error in finding user: ", err)
		return nil, err
	}
	// else
	return foundUser, nil
}

func DeleteUser(ctx context.Context, client *ent.Client, id int) error {
	// Check if user exists in db
	err := app.DbClient.User.
		DeleteOneID(id).
		Exec(ctx)

	// If error detected
	if err != nil {
		fmt.Println("error in deleting user: ", err)
		return err
	}
	// else
	return nil
}

func UpdateUser(ctx context.Context, client *ent.Client, user *models.UpdateUser) (*ent.User, error) {
	var err error
	updateQuery := client.User.
		UpdateOneID(user.Id)

	fmt.Printf("user: %v", user)

	// Check if empty as optional
	if user.Name != "" {
		updateQuery.SetUsername(user.Name)
	}
	if user.Username != "" {
		updateQuery.SetUsername(user.Username)
	}
	if user.Email != "" {
		updateQuery.SetEmail(user.Email)
	}
	// If password empty, use bcrypt to encrypt
	if user.Password != "" {
		// Build hashed password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			return nil, err
		}
		// Set in query
		updateQuery.SetPassword(string(hashedPassword))
	}

	// Save update
	createdUser, err := updateQuery.Save(ctx)
	if err != nil {

		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", createdUser)
	return createdUser, nil
}
