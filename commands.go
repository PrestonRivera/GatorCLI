package main

import (
	"errors"
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