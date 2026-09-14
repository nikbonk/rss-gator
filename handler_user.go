package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/nikbonk/rss-gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arg) == 0 {
		return fmt.Errorf("Expected one username")
	}

	if _, err := s.db.GetUser(context.Background(), cmd.arg[0]); err != nil {
		return fmt.Errorf("user not found in database: %v", err)
	}

	if err := s.configPtr.SetUser(cmd.arg[0]); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("%v has been set.", cmd.arg[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arg) == 0 {
		return fmt.Errorf("Expected one username")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arg[0],
	})
	if err != nil {
		return fmt.Errorf("Could'nt create user: %v", err)
	}

	s.configPtr.SetUser(user.Name)
	fmt.Printf("%v has been registered at %v with ID %v", user.Name, user.CreatedAt, user.ID)

	return nil

}

func handlerReset(s *state, cmd command) error {
	if len(cmd.arg) != 0 {
		return fmt.Errorf("Expected no arguments")
	}

	if err := s.db.ResetUsers(context.Background()); err != nil {
		os.Exit(1)
		return fmt.Errorf("Error while deleting users: %v", err)
	}

	fmt.Println("Successfully deleted users")
	os.Exit(0)
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	if len(cmd.arg) != 0 {
		return fmt.Errorf("Expected no arguments")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error while trying to retrieve users: %v", err)
	}

	for _, user := range users {
		if user == s.configPtr.CurrentUserName {
			fmt.Printf("%v (current) \n", user)
		} else {
			fmt.Println(user)
		}
	}

	return nil
}
