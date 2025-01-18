package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"github.com/PrestonRivera/GatorCLI/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)


func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf(" * Usage: %s <URL>", cmd.Name)
	}
	url := cmd.Args[0]
	userID := user.ID

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf(" * Failed to get feed: %w", err)
	}
	feedID := feed.ID

	err = s.db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		UserID: userID,
		FeedID: feedID,
	})
	if err != nil {
		return fmt.Errorf(" * Failed to Unfollow feed: %w", err)
	}
	fmt.Println("=========================================================")
	fmt.Printf(" * Successfully unfollowed feed: %s\n", feed.Name)
	fmt.Println("---------------------------------------------------------")
	return nil
}


func handlerFollowing(s *state, cmd command, user database.User) error {
    feeds, err := s.db.GetFeedFollowForUsers(context.Background(), user.ID)
    if err != nil {
        return fmt.Errorf(" * Failed to get users feeds: %w", err)
    }

    fmt.Println("=========================================================")
    if len(feeds) == 0 {
        fmt.Println(" * Not following any feeds")
    } else {
        for _, feed := range feeds {
            fmt.Printf(" * Following: %v\n", feed.FeedName)
        }
    }
    fmt.Println("---------------------------------------------------------")
    return nil
}


func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
	return fmt.Errorf(" * Usage: %s <URL>", cmd.Name)
	}
	url := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(),  url)
	if err != nil {
		return fmt.Errorf(" * Failed to get feed URL: %w", err)
	}

	feedID := feed.ID
	userID := user.ID

	params := database.CreateFeedFollowParams{
		UserID: userID,
		FeedID: feedID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf(" * Failed to create feed follow: %w", err)
	}

	fmt.Println("=========================================================")
	fmt.Printf(" * Feed name: %s\n", feedFollow[0].FeedName)
	fmt.Printf(" * Current User: %s\n", feedFollow[0].UserName)
	fmt.Println("---------------------------------------------------------")
	return nil
}


func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf(" * Failed to list feeds: %w", err)
	}
	fmt.Println("=========================================================")
	for _, feed := range feeds {
		printFeedList(feed)
		fmt.Println("---------------------------------------------------------")
	}
	return nil
}


func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf(" * Usage: %s <name> <url>", cmd.Name)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
		Url: url,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf(" * Failed to create feed: %w", err)
	}

	params := database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf(" * Failed to created feed Follow: %w", err)
	}

	fmt.Println("=========================================================")
	fmt.Println(" * Successfully created Feed")
	printFeed(feed)
	fmt.Println("---------------------------------------------------------")
	return nil
}


func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf(" * Usage: %v <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf(" * Invalid duration: %w", err)
	}
	ticker := time.NewTicker(timeBetweenRequests)

	fmt.Println("=========================================================")
	fmt.Printf(" * Collecting feeds every %v....\n", timeBetweenRequests)
	fmt.Println("---------------------------------------------------------")
	
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			return fmt.Errorf(" * Failed to scrape feeds: %w", err)
		}
	}
}


func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf(" * Database reset was unsuccessful: %w", err)
	}
	fmt.Println("=========================================================")
	fmt.Println(" * Database reset was successful")
	fmt.Println("---------------------------------------------------------")
	return nil
}


func handlerListUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf(" * Failed to get list of users: %w", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Println("=========================================================")
			fmt.Printf("* %s (current) \n", user.Name)
		} else {
			fmt.Printf("* %s \n", user.Name)
		}
	}
	fmt.Println("---------------------------------------------------------")
	return nil
}


func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf(" * Usage: %s <name>",cmd.Name)
	}
	name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf(" * Could not find user: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf(" * Failed to set current user: %w", err)
	}
	fmt.Println("=========================================================")
	fmt.Println(" * User has been switched successfully")
	fmt.Println("---------------------------------------------------------")
	return nil
}


func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	params := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
	}

	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				fmt.Printf(" * User already exists: %s\n", name)
				os.Exit(1)
			}
		}
		return fmt.Errorf(" * Could not create user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf(" * Could not set user: %w", err)
	}
	fmt.Println("=========================================================")
	fmt.Println(" * User created successfully:")
	printUser(user)
	fmt.Println("---------------------------------------------------------")
	return nil
}


func printUser(user database.User) {
	fmt.Printf(" * Name:	%v\n", user.Name)
	fmt.Printf(" * ID:	%v\n", user.ID)
}


func printFeed(feed database.Feed) {
	fmt.Printf(" * Name:	%s\n", feed.Name)
	fmt.Printf(" * URL: 	%s\n", feed.Url)
	fmt.Printf(" * User ID: 	%s\n", feed.UserID)
	fmt.Printf(" * Created at:  %v\n", feed.CreatedAt)
	fmt.Printf(" * Updated at:  %v\n", feed.UpdatedAt)
}

func printFeedList(feed database.ListFeedsRow) {
	fmt.Printf(" * Created by:  %s\n", feed.UsersName)
	fmt.Printf(" * Feed Name:  %s\n", feed.FeedName)
	fmt.Printf(" * Feed URL:  %s\n", feed.Url)
}