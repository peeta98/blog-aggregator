package main

import (
	"errors"
	"github.com/peeta98/blog-aggregator/internal/database"
)

type command struct {
	Name string
	Args []string
}

type commandHandler func(*state, command) error

type authenticatedCommandHandler func(*state, command, database.User) error

type commands struct {
	registeredCommands map[string]commandHandler
}

func (c *commands) register(name string, handler commandHandler) {
	c.registeredCommands[name] = handler
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return handler(s, cmd)
}

func newCommands() *commands {
	return &commands{
		registeredCommands: make(map[string]commandHandler),
	}
}
