package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arg) == 0 {
		return fmt.Errorf("Expected one username")
	}

	if err := s.configPtr.SetUser(cmd.arg[0]); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("%v has been set.", cmd.arg[0])

	return nil
}
