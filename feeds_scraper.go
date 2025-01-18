package main

import (
	"context"
	"fmt"
	"time"
)


func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf(" * Failed to get next feed to fetch: %w", err)
	}
	fmt.Println("=========================================================")
	fmt.Println()
	fmt.Println(" * Found a feed to fetch!!!!")
	fmt.Println()
	fmt.Println("---------------------------------------------------------")

	err = s.db.MarkFeedFetched(context.Background(), feed.FeedID)
	if err != nil {
		return fmt.Errorf(" * Failed to mark feed %s as fetched: %w", feed.FeedName, err)
	}
	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf(" * Failed to fetch feed: %w", err)
	}

	for _, item := range feedData.Channel.Item {
    	if item.Title != "" {
        	fmt.Println("=========================================================")
			fmt.Println()  
        	fmt.Printf("Found post: %s, URL: %s\n", item.Title,  item.Link)
			fmt.Println()
    	}
	}
	fmt.Println("=========================================================")
	fmt.Printf("Feed %s collected, %v posts found\n", feed.FeedName, len(feedData.Channel.Item))
	fmt.Printf("Feed %s collected at %v\n", 
    feed.FeedName, 
    time.Now().Format("15:04:05"))
	fmt.Println("---------------------------------------------------------")
	fmt.Println()
	return nil
}