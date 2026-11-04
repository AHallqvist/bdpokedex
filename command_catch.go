package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/AHallqvist/bdpokedex/internal/pokeapi"
)

func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 || strings.TrimSpace(args[0]) == "" {
		return fmt.Errorf("usage: catch <pokemon-name>")
	}
	pokemon, err := cfg.pokeapiClient.GetPokemon(args[0])
	if err != nil {
		return err
	}
	attemptCatch(cfg, pokemon, rand.Intn)
	return nil
}

func attemptCatch(cfg *config, pokemon pokeapi.Pokemon, roll func(int) int) {
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	if roll(100+max(0, pokemon.BaseExperience)) >= 100 {
		fmt.Printf("%s escaped!\n", pokemon.Name)
		return
	}
	if cfg.pokedex == nil {
		cfg.pokedex = make(map[string]pokeapi.Pokemon)
	}
	cfg.pokedex[pokemon.Name] = pokemon
	fmt.Printf("%s was caught!\n", pokemon.Name)
}
