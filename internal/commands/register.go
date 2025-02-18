package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/peeta98/blog-aggregator/internal/config"
	"github.com/peeta98/blog-aggregator/internal/database"
	"strings"
	"time"
)

func HandlerRegister(s *config.State, cmd *Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("register command requires single argument <username>")
	}

	if len(cmd.Args) != 1 {
		return errors.New("register command only uses one username")
	}

	username := cmd.Args[0]
	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return fmt.Errorf("user %s already exists: %v", username, err)
		}
		return fmt.Errorf("unable to create user in DB: %v", err)
	}

	err = s.Config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("User successfully created with the following data: %+v\n", user)

	return nil
}
