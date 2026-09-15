package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nikbonk/rss-gator/internal/database"
)

func handlerFetch(s *state, cmd command) error {
	// if len(cmd.arg) == 0 {
	// 	return fmt.Errorf("Expected a URL")
	// }

	// url := cmd.arg[0]
	// quick static site check
	url := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Error while trying to fetch rss feed: %v", err)
	}

	fmt.Printf("%+v\n", feed)

	return nil

}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.arg) != 2 {
		return fmt.Errorf("Requires name and URL")
	}

	name := cmd.arg[0]
	url := cmd.arg[1]
	user, err := s.db.GetUser(context.Background(), s.configPtr.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error while trying to retrieve gatorconfig user: %v", err)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("Error while trying to add feed: %v", err)
	}

	fmt.Printf("%+v\n", feed)

	return nil

}
