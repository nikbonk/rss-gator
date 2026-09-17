package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/nikbonk/rss-gator/internal/config"
	"github.com/nikbonk/rss-gator/internal/database"
)

type state struct {
	configPtr *config.Config
	db        *database.Queries
}

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("Error trying to open sql connection: %v", err)
	}

	dbQueries := database.New(db)
	appState := state{configPtr: &cfg, db: dbQueries}

	cmds := commands{cmdMap: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerFetch)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerGetUsersFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollowFeed))
	cmds.register("following", middlewareLoggedIn(handlerGetFeedFollows))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollowFeed))

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmd := command{name: os.Args[1], arg: os.Args[2:]}
	if err := cmds.run(&appState, cmd); err != nil {
		log.Fatal(err)
	}
}
