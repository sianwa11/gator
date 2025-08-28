package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sianwa11/gator/internal/database"
)

func handleRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("name is required in the command")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.Args[0],
	})

	if err != nil {
		log.Fatal("user already exists")
	}

	err = s.cfg.SetUser(user.Name)
		if err != nil {
		return fmt.Errorf("couldn't set user: %w", err)
	}

	fmt.Println("user was created")
	fmt.Printf("%v\n", user)
	
	return nil
}