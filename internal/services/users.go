package services

import (
	"context"
	"fmt"
	"log"

	"github.com/dmawardi/Go-Template/ent"
	"github.com/dmawardi/Go-Template/internal/models"
)

func (repo *Repository) CreateUser(ctx context.Context, client *ent.Client, user *models.CreateUser) (*ent.User, error) {
	u, err := client.User.
		Create().
		SetName(user.Name).
		SetUsername(user.Username).
		SetEmail(user.Email).
		SetPassword(user.Password).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}
