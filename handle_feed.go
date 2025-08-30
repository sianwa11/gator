package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
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

func handleAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("name and url required for this command")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]
	
	ctx := context.Background()

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

	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: newFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed follow: %w", err)
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

func handleBrowse(s *state, cmd command) error {
	limit := "2"
	if len(cmd.Args) > 0 {
		limit = cmd.Args[0]
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return fmt.Errorf("invalid limit: %w", err)
	}

	ctx := context.Background()

	currUser, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not get current user")
	}
	

	posts, err := s.db.GetPostsForUser(ctx, database.GetPostsForUserParams{
		UserID: currUser.ID,
		Limit: int32(limitInt),
	})
	if err != nil {
		return fmt.Errorf("error getting posts for user: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("%v\n", post.Name)
		fmt.Printf("%v\n", post.Title)
		fmt.Printf("%v\n", post.Description)
		fmt.Printf("%v\n", post.FeedUrl)
		fmt.Printf("%v\n", post.PublishedAt)
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("could not return next feed: %w", err)
	}

	err = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: time.Now(),
		ID: nextFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to mark feed as fetched: %w", err)
	}

	feed, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return fmt.Errorf("could not find feed: %w", err)
	}

	for _, item := range feed.Channel.Item {
		pubTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("failed to parse pubDate '%s': %v", item.PubDate, err)
			continue
		}
		_, err = s.db.CreatePost(ctx, database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			Title: item.Title,
			Url: item.Link,
			PublishedAt: pubTime,	
			Description: item.Description,
			FeedID: nextFeed.ID,		
		})

		if err != nil {
			return fmt.Errorf("error creating post: %w", err)
		}
    log.Printf("Saved post: %s", item.Title)
	}

	log.Printf("Feed %s collected, %v posts found", nextFeed.Name, len(feed.Channel.Item))


	return nil
}