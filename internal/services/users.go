package services

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/dmawardi/Go-Template/ent"
)

func (repo *Repository) CreateUser(ctx context.Context, client *ent.Client, r *http.Request) (*ent.User, error) {
	u, err := client.User.
		Create().
		SetName("a8m").
		SetUsername("gonad").
		SetEmail("dopey@gmail.com").
		SetPassword("goose").
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}
