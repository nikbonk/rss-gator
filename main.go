package main

import (
	"log"
	"os"

	"github.com/nikbonk/rss-gator/internal/config"
)

type state struct {
	configPtr *config.Config
}

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	appState := state{configPtr: &cfg}
	cmds := commands{cmdMap: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmd := command{name: os.Args[1], arg: os.Args[2:]}
	if err := cmds.run(&appState, cmd); err != nil {
		log.Fatal(err)
	}
}
