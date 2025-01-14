package main

import (
	"fmt"
	"log"
	"os"
	"database/sql"
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

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	
	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments were provided")
		return 
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
