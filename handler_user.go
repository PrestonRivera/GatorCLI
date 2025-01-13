package main

import (
	"fmt"
)


func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>",cmd.Name)
	}
	name := cmd.Args[0]

	err := s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("Current user couldn't be set: %s", err)
	}
	
	fmt.Println("User has been switched successfully")
	return nil
}