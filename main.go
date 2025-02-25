package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/peeta98/blog-aggregator/internal/config"
	"github.com/peeta98/blog-aggregator/internal/database"
	"log"
	"os"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	programState := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	cli := newCommands()
	cli.register("login", handlerLogin)
	cli.register("register", handlerRegister)
	cli.register("reset", handlerReset)
	cli.register("users", handlerListUsers)
	cli.register("agg", handlerAggregate)
	cli.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cli.register("feeds", handlerListFeeds)
	cli.register("follow", middlewareLoggedIn(commands.HandlerFollowFeed))
	cli.register("following", middlewareLoggedIn(commands.HandlerListFeedFollows))
	cli.register("unfollow", middlewareLoggedIn(commands.HandlerUnfollowFeed))
	cli.register("browse", middlewareLoggedIn(commands.HandlerBrowsePosts))

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	cmd := command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err := cli.run(programState, cmd); err != nil {
		log.Fatal(err)
	}
}

func middlewareLoggedIn(handler authenticatedCommandHandler) commandHandler {
	return func(state *state, cmd command) error {
		user, err := state.db.GetUser(context.Background(), state.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to get authenticated user: %v", err)
		}

		return handler(state, cmd, user)
	}
}
