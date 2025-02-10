package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/PrestonRivera/GatorCLI/internal/config"
	"github.com/PrestonRivera/GatorCLI/internal/database"
	_ "github.com/lib/pq"
)


type state struct {
	db  *database.Queries
	cfg *config.Config
}

const dbURL = "postgres://postgres:postgres@localhost:5432/gator"

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config: ", err)
		return 
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Set up for db connection and queries failed")
	}
	defer db.Close()
	dbQueries := database.New(db)

	programState := &state{
		db: dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		cliCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("help", handlerHelp)
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollowFeed))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}


func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf(" * Failed to get usser: %w", err)
		}
		return handler(s, cmd, user)
	}
}