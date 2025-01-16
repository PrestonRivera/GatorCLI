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


func handlerAddFeed(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
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
		return fmt.Errorf("Failed to create feed: %w", err)
	}
	fmt.Println("Succesfully created Feed")
	fmt.Println("=========================================================")
	printFeed(feed)
	fmt.Println()
	return nil
}


func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("Failed to fetch feed:  %w", err)
	}
	fmt.Printf("Feed: %+v\n", feed)
	return nil
}


func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Database reset was unsuccessful: %w", err)
	}	
	fmt.Println("Database reset was successful")
	return nil
}


func handlerListUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get list of users: %w", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current) \n", user.Name)
		} else {
			fmt.Printf("* %s \n", user.Name)
		}
	}
	return nil
}


func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>",cmd.Name)
	}
	name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("Could not find user: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("Could not set current user: %w", err)
	}

	fmt.Println("User has been switched successfully")
	return nil
}


func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
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
				fmt.Printf("User already exists: %s\n", name)
				os.Exit(1)
			}
		}
		return fmt.Errorf("Could not create user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Could not set user: %w", err)
	}

	fmt.Println("User created successfully:")
	printUser(user)
	return nil
}


func printUser(user database.User) {
	fmt.Printf(" * ID:		%v\n", user.ID)
	fmt.Printf(" * Name:	%v\n", user.Name)
}


func printFeed(feed database.Feed) {
	fmt.Printf(" * Name:	%s\n", feed.Name)
	fmt.Printf(" * URL: 	%s\n", feed.Url)
	fmt.Printf(" * User ID: 	%s\n", feed.UserID)
	fmt.Printf(" * Created at:  %v\n", feed.CreatedAt)
	fmt.Printf(" * Updated at:  %v\n", feed.UpdatedAt)
}