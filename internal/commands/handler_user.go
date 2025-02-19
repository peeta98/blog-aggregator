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

func HandlerLogin(s *config.State, cmd *Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("login command requires single argument <username>")
	}

	if len(cmd.Args) != 1 {
		return errors.New("login command only uses one username")
	}

	username := cmd.Args[0]

	_, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		return errors.New("username doesn't exist")
	}

	err = s.Config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("User '%s' has been set in the config file!\n", username)

	return nil
}

func HandlerListUsers(s *config.State, cmd *Command) error {
	if len(cmd.Args) > 0 {
		return errors.New("command <users> doesn't accept args")
	}

	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get users from DB: %v\n", err)
	}

	for _, user := range users {
		fmt.Print(formatUser(user.Name, s.Config.CurrentUserName))
	}

	return nil
}

func formatUser(username, currentUser string) string {
	if username == currentUser {
		return fmt.Sprintf("* %s (current)\n", username)
	}
	return fmt.Sprintf("* %s\n", username)
}
