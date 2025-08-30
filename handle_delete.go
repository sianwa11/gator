package main

import (
	"context"
	"fmt"

	"github.com/sianwa11/gator/internal/database"
)

func handleDelete(s *state, cmd command) error {
	err := s.db.DeleteUser(context.Background())
	if err != nil {}
	fmt.Println("users successfully deleted")
	return nil
}

func handleUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("command required a url")
	}

	url := cmd.Args[0]

	ctx := context.Background()

	feed, err := s.db.FindFeed(ctx, url)
	if err != nil {
		return fmt.Errorf("could not find feed with url: %w", err)
	}

	err = s.db.DeleteFeedFollowsForUser(ctx, database.DeleteFeedFollowsForUserParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error unfollowing feed: %w", err)
	}

	return nil
}