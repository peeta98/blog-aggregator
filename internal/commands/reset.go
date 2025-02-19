package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/peeta98/blog-aggregator/internal/config"
)

func HandlerReset(s *config.State, cmd *Command) error {
	if len(cmd.Args) > 0 {
		return errors.New("command <reset> doesn't accept args")
	}

	if err := s.Db.DeleteUsers(context.Background()); err != nil {
		return fmt.Errorf("couldn't delete users from DB: %v\n", err)
	}

	fmt.Println("All users have been deleted successfully!")
	return nil
}
