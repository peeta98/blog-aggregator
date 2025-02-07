package commands

import (
	"errors"
	"fmt"
	"github.com/peeta98/blog-aggregator/internal/config"
)

func HandlerLogin(s *config.State, cmd *Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("login command requires single argument <username>")
	}

	if len(cmd.Args) != 1 {
		return errors.New("login command only uses one username")
	}

	username := cmd.Args[0]
	err := s.Config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("User '%s' has been set in the config file!\n", username)

	return nil
}
