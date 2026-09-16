package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
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
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("Error while trying to add feed: %v", err)
	}

	fmt.Printf("%+v\n", feed)

	followCmd := command{name: "follow", arg: []string{url}}
	err = handlerFollowFeed(s, followCmd)
	if err != nil {
		return fmt.Errorf("Error while trying to add feed to follow table: %v", err)
	}

	return nil

}

func handlerGetUsersFeeds(s *state, cmd command) error {
	if len(cmd.arg) != 0 {
		return fmt.Errorf("Expected no arguments")
	}

	feeds, err := s.db.GetUsersFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error while trying to retrieve feeds: %v", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 1, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "ID\tNAME\tURL\tUSER\tCREATED")
	fmt.Fprintln(w, "--\t----\t---\t-------\t-------")

	for _, feed := range feeds {
		fmt.Fprintf(
			w,
			"%s\t%s\t%s\t%s\t%s\n",
			feed.ID,
			feed.Name_2,
			feed.Url,
			feed.Name,
			feed.CreatedAt.Format("2006-01-02 15:04"),
		)
	}

	return nil
}

func handlerFollowFeed(s *state, cmd command) error {
	if len(cmd.arg) == 0 {
		return fmt.Errorf("Expected a URL")
	}

	url := cmd.arg[0]

	currentUser, err := s.db.GetUser(context.Background(), s.configPtr.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error while trying to retrieve gatorconfig user: %v", err)
	}

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Error while trying to retrieve feed: %v", err)
	}

	user, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error while trying to follow a feed: %v", err)
	}

	fmt.Printf("%v followed feed %v", user.UserName, feed.Name)

	return nil
}

func handlerGetFeedFollows(s *state, cmd command) error {
	if len(cmd.arg) != 0 {
		return fmt.Errorf("Expected no arguments")
	}

	currentUser, err := s.db.GetUser(context.Background(), s.configPtr.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error while trying to retrieve gatorconfig user: %v", err)
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return fmt.Errorf("Error while trying to retrieved followed feeds: %v", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 1, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "ID\tFEED\tUSER\tFOLLOWED AT")
	fmt.Fprintln(w, "--\t----\t----\t----------")

	for _, follow := range follows {
		fmt.Fprintf(
			w,
			"%v\t%v\t%v\t%v\n",

			follow.ID,
			follow.FeedName,
			follow.UserName,
			follow.CreatedAt.Format("2006-01-02 15:04"),
		)
	}

	return nil
}
