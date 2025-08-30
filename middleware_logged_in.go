package main

import (
	"context"
	"fmt"

	"github.com/sianwa11/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		if s.cfg.CurrentUserName == "" {
			return fmt.Errorf("you must be logged in to use this command")
		}

		usr, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error getting current user: %w", err)
		}

		return handler(s, cmd, usr)
	}
}