package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/peeta98/blog-aggregator/internal/config"
	"github.com/peeta98/blog-aggregator/internal/database"
	"github.com/peeta98/blog-aggregator/internal/rss"
	"net/url"
	"time"
)

func HandlerAddFeed(s *config.State, cmd *Command) error {
	if len(cmd.Args) != 2 {
		return errors.New("addfeed command requires two arguments <name> <url>")
	}

	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]
	err := validateFeedUrl(feedUrl)
	if err != nil {
		return err
	}

	user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get authenticated user: %v", err)
	}

	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %v", err)
	}

	fmt.Printf("Feed successfully created with name '%s' and url '%s'\n", feed.Name, feed.Url)

	return nil
}

func validateFeedUrl(feedUrl string) error {
	parsedUrl, err := url.Parse(feedUrl)
	if err != nil {
		return fmt.Errorf("invalid feed URL: %v", err)
	}
	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		return errors.New("feed URL must use HTTP or HTTPS protocol")
	}
	return nil
}

func HandlerAggregate(s *config.State, cmd *Command) error {
	if len(cmd.Args) > 0 {
		return errors.New("command <agg> doesn't accept args")
	}

	rssFeed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("couldn't fetch rss feed from the URL provided: %v", err)
	}

	printRSSFeed(rssFeed)
	return nil
}

func printRSSFeed(rssFeed *rss.Feed) {
	fmt.Printf("Channel: %s\n", rssFeed.Channel.Title)
	fmt.Printf("Description: %s\n", rssFeed.Channel.Description)
	fmt.Println("\nArticles:")
	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("\nTitle: %s\n", item.Title)
		fmt.Printf("Link: %s\n", item.Link)
		fmt.Printf("Description: %s\n", item.Description)
		fmt.Printf("Published: %s\n", item.PubDate)
	}
}
