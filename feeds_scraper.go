package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PrestonRivera/GatorCLI/internal/database"
	"github.com/google/uuid"
)


func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf(" * Failed to get next feed to fetch: %w", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), feed.FeedID)
	if err != nil {
		return fmt.Errorf(" * Failed to mark feed %s as fetched: %w", feed.FeedName, err)
	}
	feedData, err := fetchFeed(context.Background(), feed.Url)
	
	if err != nil {
		return fmt.Errorf(" * Failed to fetch feed: %w", err)
	}

	newPosts := 0
	skipPosts := 0

	for _, item := range feedData.Channel.Item {
		var description sql.NullString
    	if item.Description != "" {
     		description = sql.NullString{
        		String: item.Description,
            	Valid: true,
        	}
    	}

	publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
	if err != nil {
    	log.Printf(" * Failed to parse date %q: %v", item.PubDate, err)
    	publishedAt = time.Now()
	}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: item.Title,
			Url: item.Link,
			Description: description,
			PublishedAt: publishedAt,
			FeedID: feed.FeedID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "unique constraint") {
				continue
			}
			log.Printf(" * Failed to create post %q: %v", item.Title, err)
			continue
		}
		newPosts++
	}
	fmt.Println("=========================================================")
	fmt.Printf(" * Scrape completed for: %s:\n", feed.FeedName)
    fmt.Printf(" * New posts added: %d\n", newPosts)
    fmt.Printf(" * Duplicate posts skipped: %d\n", skipPosts)
	fmt.Println("---------------------------------------------------------")
	return nil
}