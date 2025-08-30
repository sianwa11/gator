package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sianwa11/gator/internal/database"
)

func handleFeedFollows(s *state, cmd command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("expects a url argument")
	}
	url := cmd.Args[0]

	ctx := context.Background()



	feed, err := s.db.FindFeed(ctx, url)
	if err != nil {
		return fmt.Errorf("error finding feed %v: %w", url, err)
	}

	feedFollow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feed follow: %w", err)
	}

	fmt.Printf("%v\n", feedFollow.FeedName)
	fmt.Printf("%v\n", feedFollow.UserName)

	
	return nil
}

func handleFeedFollowing(s *state, cmd command, user database.User) error {
	usersFeedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not get users feed follows: %w", err)
	}

	for _, feedFollows := range usersFeedFollows {
		fmt.Printf("%v\n", feedFollows.FeedName)
	}

	return nil
}