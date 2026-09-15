package main

import (
	"context"
	"fmt"
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
