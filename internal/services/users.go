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

func FindUserByEmail(ctx context.Context, client *ent.Client, email string) (*ent.User, error) {
	fmt.Println("Finding user by username...")
	// Check if user exists in db
	// fmt.Println("repo ctx:", repo.App.Ctx)

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
