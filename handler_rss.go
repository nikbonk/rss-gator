package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/google/uuid"
	"github.com/nikbonk/rss-gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.arg) == 0 {
		return fmt.Errorf("Requires interval")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.arg[0])
	if err != nil {
		return fmt.Errorf("Invalid interval: %v", err)
	}

	ticker := time.NewTicker(timeBetweenRequests)
	fmt.Printf("Collecting feeds every %v...", timeBetweenRequests)

	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return fmt.Errorf("Error while scraping feeds: %v", err)
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arg) != 2 {
		return fmt.Errorf("Requires name and URL")
	}

	name := cmd.arg[0]
	url := cmd.arg[1]

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
	err = handlerFollowFeed(s, followCmd, user)
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

func handlerFollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arg) == 0 {
		return fmt.Errorf("Expected a URL")
	}

	url := cmd.arg[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Error while trying to retrieve feed: %v", err)
	}

	followedFeed, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error while trying to follow a feed: %v", err)
	}

	fmt.Printf("%v followed feed %v", followedFeed.UserName, followedFeed.FeedName)

	return nil
}

func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arg) == 0 {
		return fmt.Errorf("Expected a feed URL")
	}

	url := cmd.arg[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Error while trying to retrieve feed: %v", err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error while trying to unfollow feed: %v", err)
	}

	fmt.Printf("%v successfully unfollowed %v\n", user.Name, feed.Name)

	return nil

}

func handlerGetFeedFollows(s *state, cmd command, user database.User) error {
	if len(cmd.arg) != 0 {
		return fmt.Errorf("Expected no arguments")
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
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

func handlerBrowseFeeds(s *state, cmd command, user database.User) error {
	var limit string
	if len(cmd.arg) == 0 {
		limit = "2"
	} else {
		limit = cmd.arg[0]
	}

	postLimit, err := strconv.Atoi(limit)
	if err != nil {
		return fmt.Errorf("Error while trying to parse post limit: %v\n", err)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(postLimit),
	})
	if err != nil {
		return fmt.Errorf("Error while trying to retrieve posts: %v", err)
	}

	for _, post := range posts {
		p := PostItem{
			Title:       post.Title,
			PublishedAt: post.PublishedAt.Time,
			URL:         post.Url,
			Description: post.Description.String,
		}
		printPost(p)
	}

	return nil
}

type PostItem struct {
	Title       string
	PublishedAt time.Time
	URL         string
	Description string
}

func printPost(p PostItem) {
	divider := strings.Repeat("-", 60)

	fmt.Println(divider)
	// Bold title: \033[1m enables bold, \033[0m resets
	fmt.Printf("\033[1m%s\033[0m\n", p.Title)

	if !p.PublishedAt.IsZero() {
		fmt.Printf("Published: %s\n", p.PublishedAt.Format("2006-01-02 15:04"))
	}

	fmt.Printf("Link:      %s\n", p.URL)

	if p.Description != "" {
		fmt.Printf("\n%s\n", p.Description)
	}
}
