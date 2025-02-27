package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/peeta98/blog-aggregator/internal/database"
	"strconv"
)

func handlerBrowsePosts(s *state, cmd command, user database.User) error {
	var limitPosts int32
	if len(cmd.Args) != 1 {
		// If optional "limit" argument is not provided, default the limit to 2
		limitPosts = 2
	} else {
		limit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		if limit <= 0 {
			return errors.New("limit must be a positive number")
		}
		limitPosts = int32(limit)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limitPosts,
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	printPosts(posts, user.Name)
	return nil
}

func printPosts(posts []database.GetPostsForUserRow, userName string) {
	if len(posts) == 0 {
		fmt.Printf("No posts found for user %s\n", userName)
		return
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), userName)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}
}
