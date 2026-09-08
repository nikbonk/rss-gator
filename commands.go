package main

import "fmt"

type command struct {
	name string
	arg  []string
}

type commands struct {
	cmdMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if handler, ok := c.cmdMap[cmd.name]; ok {
		return handler(s, cmd)
	}

	return fmt.Errorf("Unknown command")
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmdMap[name] = f
}
