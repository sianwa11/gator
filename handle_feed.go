package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sianwa11/gator/internal/database"
)

func handleFeed(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"

	feed, err := fetchFeed(context.Background(), url)

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", feed)

	return nil
}

func handleAddFeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("name and url required for this command")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]
	
	ctx := context.Background()
	user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("current user does not exist")
	}

	newFeed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
		Url: url,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("error occured creating feed: %w", err)
	}

	fmt.Printf("%v\n", newFeed)
	return nil
}

func handleGetFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("error fectching feeds: %w", err)
	}

	for _, feed := range feeds {
		fmt.Printf("%s\n", feed.Name)
		fmt.Printf("%s\n", feed.Url)
		name, err := s.db.FindUser(ctx, feed.UserID)
		if err != nil {
			log.Fatal("could not find user with id")
		}
		fmt.Printf("%s\n", name)
	}

	return nil
}