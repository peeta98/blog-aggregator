package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/peeta98/blog-aggregator/internal/config"
	"github.com/peeta98/blog-aggregator/internal/database"
	"github.com/peeta98/blog-aggregator/internal/rss"
	"log"
	"net/url"
	"strconv"
	"time"
)

func HandlerBrowsePosts(s *config.State, cmd *Command, user database.User) error {
	var limitPosts int32
	if len(cmd.Args) != 1 {
		// If optional "limit" parameter is not provided, default the limit to 2
		limitPosts = 2
	} else {
		limit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %v", err)
		}
		if limit <= 0 {
			return fmt.Errorf("limit must be a positive number")
		}
		limitPosts = int32(limit)
	}

	posts, err := s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limitPosts,
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts: %v", err)
	}

	printPosts(posts)

	return nil
}

func printPosts(posts []database.GetPostsForUserRow) {
	if len(posts) == 0 {
		fmt.Println("You have no posts to browse. Try following some feeds first!")
		return
	}

	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}
}

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
	timeBetweenRequests, err := time.ParseDuration("1m0s")
	if err != nil {
		return fmt.Errorf("invalid duration: %v", err)
	}
	fmt.Println("Collecting feeds every 1m0s")

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func scrapeFeeds(s *config.State) error {
	nextFeed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't fetch next feed to scrape: %v", err)
	}

	err = s.Db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return fmt.Errorf("couldn't mark feed as fetched: %v", err)
	}

	rssFeed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}

	err = savePosts(rssFeed.Channel.Item, &nextFeed, s)
	log.Printf("Feed %s collected, %v posts found", nextFeed.Name, len(rssFeed.Channel.Item))

	return nil
}

func savePosts(items []rss.Item, feed *database.Feed, s *config.State) error {
	for _, item := range items {
		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			// Try alternative format if first parse fails
			publishedAt, err = time.Parse(time.RFC822, item.PubDate)
			if err != nil {
				return fmt.Errorf("could not parse date %s: %v", item.PubDate, err)
			}
		}

		_, err = s.Db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			return fmt.Errorf("could not create post: %v", err)
		}
	}
	fmt.Println("Posts have been successfully saved!")
	return nil
}
