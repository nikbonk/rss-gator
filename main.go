package main

import (
	"fmt"
	"log"

	"github.com/nikbonk/rss-gator/internal/config"
)

func readConfig() config.Config {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	return cfg

}

func main() {
	cfg := readConfig()

	fmt.Printf("Readon config: %+v\n", cfg)

	err := cfg.SetUser("niklas")
	if err != nil {
		log.Fatalf("error setting username in config: %v", err)
	}

	cfg = readConfig()
	fmt.Printf("Read config again: %+v\n", cfg)

}
