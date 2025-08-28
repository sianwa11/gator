package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/sianwa11/gator/internal/config"
	"github.com/sianwa11/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db *database.Queries
	cfg *config.Config
}


func main(){
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error opening db connection: %v", err)
	}

	dbQueries := database.New(db)

	programState := &state{
		cfg: &cfg,
		db: dbQueries,
	}

	cmds := &commands{
	registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handleRegister)
	cmds.register("reset", handleDelete)
	cmds.register("users", handleUsers)
	cmds.register("agg", handleFeed)
	cmds.register("addfeed", handleAddFeed)
	cmds.register("feeds", handleGetFeeds)

	
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]


	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}

}