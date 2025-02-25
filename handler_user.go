package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/peeta98/blog-aggregator/internal/database"
	"strings"
	"time"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}

	username := cmd.Args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      username,
	})
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return fmt.Errorf("user %s already exists: %v", username, err)
		}
		return fmt.Errorf("couldn't create user: %v", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User created successfully:")
	printUser(user)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}

	username := cmd.Args[0]
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("username doesn't exist: %v", err)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User '%s' is now logged in!\n", username)
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't list users: %v\n", err)
	}

	for _, user := range users {
		fmt.Print(formatUser(user.Name, s.cfg.CurrentUserName))
	}

	return nil
}

func formatUser(username, currentUser string) string {
	if username == currentUser {
		return fmt.Sprintf("* %s (current)\n", username)
	}
	return fmt.Sprintf("* %s\n", username)
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
