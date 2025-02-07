package main

import (
	"github.com/peeta98/blog-aggregator/internal/commands"
	"github.com/peeta98/blog-aggregator/internal/config"
	"log"
	"os"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	state := config.State{
		Config: &cfg,
	}

	cli := commands.NewCommands()
	cli.Register("login", commands.HandlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("not enough arguments")
	}

	cmd := &commands.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err := cli.Run(&state, cmd); err != nil {
		log.Fatal(err)
	}
}
