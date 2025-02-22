package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/peeta98/blog-aggregator/internal/commands"
	"github.com/peeta98/blog-aggregator/internal/config"
	"github.com/peeta98/blog-aggregator/internal/database"
	"log"
	"os"
)

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

	state := config.State{
		Config: &cfg,
		Db:     dbQueries,
	}

	cli := commands.NewCommands()
	cli.Register("login", commands.HandlerLogin)
	cli.Register("register", commands.HandlerRegister)
	cli.Register("reset", commands.HandlerReset)
	cli.Register("users", commands.HandlerListUsers)
	cli.Register("agg", commands.HandlerAggregate)
	cli.Register("addfeed", commands.MiddlewareLoggedIn(commands.HandlerAddFeed))
	cli.Register("feeds", commands.HandlerListFeeds)
	cli.Register("follow", commands.MiddlewareLoggedIn(commands.HandlerFollowFeed))
	cli.Register("following", commands.MiddlewareLoggedIn(commands.HandlerListFeedFollows))
	cli.Register("unfollow", commands.MiddlewareLoggedIn(commands.HandlerUnfollowFeed))

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	cmd := &commands.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err := cli.Run(&state, cmd); err != nil {
		log.Fatal(err)
	}
}
