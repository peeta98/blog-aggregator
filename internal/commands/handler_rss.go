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

func HandlerListFeeds(s *config.State, cmd *Command) error {
	if len(cmd.Args) != 0 {
		return errors.New("command <feeds> doesn't accept args")
	}

	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't fetch feeds from DB: %v", err)
	}

	printFeeds(feeds)

	return nil
}

func printFeeds(feeds []database.GetFeedsRow) {
	for i, feed := range feeds {
		fmt.Printf("Feed %d\n", i+1)
		fmt.Printf("Name of the feed: %s\n", feed.Name)
		fmt.Printf("URL of the feed: %s\n", feed.Url)
		fmt.Printf("User that created the feed: %s\n", feed.UserName)
	}
}

func HandlerAddFeed(s *config.State, cmd *Command, user database.User) error {
	if len(cmd.Args) != 2 {
		return errors.New("command <addfeed> requires two arguments <name> <url>")
	}

	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]
	err := validateFeedUrl(feedUrl)
	if err != nil {
		return err
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

	_, err = s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %v", err)
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

func HandlerListFeedFollows(s *config.State, cmd *Command, user database.User) error {
	if len(cmd.Args) != 0 {
		return errors.New("command <following> doesn't accept any args")
	}

	feedFollows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("couldn't get feeds that user follows: %v", err)
	}

	printFollowedFeedNames(feedFollows)
	return nil
}

func printFollowedFeedNames(feedFollows []database.GetFeedFollowsForUserRow) {
	if len(feedFollows) == 0 {
		fmt.Println("User is currently not following any feed.")
		return
	}

	fmt.Println("List of Feeds that current user follows:")
	for i, feedFollow := range feedFollows {
		fmt.Printf("Feed %d: %s\n", i+1, feedFollow.FeedName)
	}
}

func HandlerFollowFeed(s *config.State, cmd *Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return errors.New("command <follow> requires one argument <feedUrl>")
	}

	feedUrl := cmd.Args[0]
	if err := validateFeedUrl(feedUrl); err != nil {
		return err
	}

	feed, err := s.Db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("couldn't get feed based on the current URL: %v", err)
	}

	feedFollow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %v", err)
	}

	fmt.Printf("User %s now follows %s!\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func HandlerUnfollowFeed(s *config.State, cmd *Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return errors.New("command <unfollow> requires one argument <feedUrl>")
	}
	feedUrl := cmd.Args[0]
	if err := validateFeedUrl(feedUrl); err != nil {
		return err
	}

	err := s.Db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    feedUrl,
	})
	if err != nil {
		return fmt.Errorf("couldn't unfollow feed: %v", err)
	}

	fmt.Printf("%s has successfully unfollowed %s!\n", user.Name, feedUrl)

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
