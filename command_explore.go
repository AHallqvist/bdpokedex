package main

import (
	"fmt"
	"strings"
)

func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 || strings.TrimSpace(args[0]) == "" {
		return fmt.Errorf("usage: explore <area-name>")
	}
	area, err := cfg.pokeapiClient.GetLocationArea(args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\nFound Pokemon:\n", args[0])
	if len(area.PokemonEncounters) == 0 {
		fmt.Println("No Pokemon found.")
	}
	for _, encounter := range area.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}
