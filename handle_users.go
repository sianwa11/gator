package main

import (
	"context"
	"fmt"
	"log"
)

func handleUsers(s *state, cmd command) error{
	users, err := s.db.GetUsers(context.Background())

	if err != nil {
		log.Fatal("failed to get users")
	}
	

	for _, user := range users {
		currentUser := ""
		if user.Name == s.cfg.CurrentUserName {
			currentUser = "(current)"
		}
		fmt.Printf("* %v %v\n", user.Name, currentUser)
	}

	return nil
}