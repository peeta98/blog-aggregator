package commands

import (
	"fmt"
	"github.com/peeta98/blog-aggregator/internal/config"
	"github.com/peeta98/blog-aggregator/internal/database"
)

type CommandHandler func(*config.State, Command) error

type AuthenticatedCommandHandler func(*config.State, *Command, database.User) error

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]CommandHandler
}

func (c *Commands) Register(name string, handler CommandHandler) {
	c.Handlers[name] = handler
}

func (c *Commands) Run(s *config.State, cmd *Command) error {
	handler, exists := c.Handlers[cmd.Name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}

	return handler(s, cmd)
}

func NewCommands() *Commands {
	return &Commands{
		Handlers: make(map[string]CommandHandler),
	}
}
