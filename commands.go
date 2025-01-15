package main

import (
	"context"
	"errors"
	"fmt"
)


type command struct {
	Name string
	Args []string
}


type commands struct {
	cliCommands map[string]func(*state, command) error
}


func (c *commands) register(name string, f func(*state, command) error) {
	c.cliCommands[name] = f
}


func (c *commands) run(s *state, cmd command) error {
	f, ok := c.cliCommands[cmd.Name]
	if !ok {
		return errors.New("run command not found")
	}
	return f(s, cmd)
}


func (c *commands) reset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Database reset was unsuccessful: %w", err)
	}	
	fmt.Println("Database reset was successful")
	return nil
}


func (c *commands) users(s *state, cmd command) error {
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