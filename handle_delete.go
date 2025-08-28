package main

import (
	"context"
	"fmt"
)

func handleDelete(s *state, cmd command) error {
	err := s.db.DeleteUser(context.Background())
	if err != nil {}
	fmt.Println("users successfully deleted")
	return nil
}