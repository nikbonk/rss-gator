package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nikbonk/rss-gator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Error while getting next feed to fetch: %v", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return fmt.Errorf("Error while marking feed fetched: %v", err)
	}

	feeds, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("Error while fetching feed: %v", err)
	}

	timeLayouts := []string{time.RFC1123Z, time.RFC1123}
	for _, feed := range feeds.Channel.Item {
		var publishDate time.Time
		for _, layout := range timeLayouts {
			if parsedDate, err := time.Parse(layout, feed.PubDate); err == nil {
				publishDate = parsedDate
				break
			}
		}
		publishedAt := sql.NullTime{
			Time:  publishDate,
			Valid: !publishDate.IsZero(),
		}

		description := sql.NullString{
			String: feed.Description,
			Valid:  feed.Description != "",
		}

		err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       feed.Title,
			Url:         feed.Link,
			Description: description,
			PublishedAt: publishedAt,
			FeedID:      nextFeed.ID,
		})
		if err != nil {
			return fmt.Errorf("Error while creating post: %v", err)
		}
	}

	return nil
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	userAgent := "rss-gator"
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	feed := &RSSFeed{}
	err = xml.Unmarshal(buf, feed)
	if err != nil {
		return nil, err
	}

	unescapeStrings(feed)
	return feed, nil
}

func unescapeStrings(feed *RSSFeed) {
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
}
