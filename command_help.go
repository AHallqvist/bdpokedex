package main

import (
	"fmt"
	"sort"
)

func commandHelp(cfg *config, args ...string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	names := make([]string, 0, len(cfg.commands))
	for name := range cfg.commands {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		command := cfg.commands[name]
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}
