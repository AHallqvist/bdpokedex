package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/AHallqvist/bdpokedex/internal/pokeapi"
)

type config struct {
	commands            map[string]cliCommand
	pokeapiClient       pokeapi.Client
	pokedex             map[string]pokeapi.Pokemon
	Next                *string
	Previous            *string
	locationAreasLoaded bool
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func startRepl(cfg *config) {
	cfg.commands = map[string]cliCommand{
		"catch": {
			name:        "catch",
			description: "Try to catch a Pokemon: catch <pokemon-name>",
			callback:    commandCatch,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location area: explore <area-name>",
			callback:    commandExplore,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a caught Pokemon: inspect <pokemon-name>",
			callback:    commandInspect,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location areas",
			callback:    commandMapb,
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}
		command, ok := cfg.commands[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		if err := command.callback(cfg, words[1:]...); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read input:", err)
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
