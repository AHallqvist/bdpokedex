package main

import (
	"encoding/json"
	"testing"

	"github.com/AHallqvist/bdpokedex/internal/pokeapi"
)

func TestCommandInspectCaughtPokemon(t *testing.T) {
	var pokemon pokeapi.Pokemon
	err := json.Unmarshal([]byte(`{
		"name": "pidgey", "base_experience": 50, "height": 3, "weight": 18,
		"stats": [
			{"base_stat": 40, "stat": {"name": "hp"}},
			{"base_stat": 45, "stat": {"name": "attack"}},
			{"base_stat": 40, "stat": {"name": "defense"}},
			{"base_stat": 35, "stat": {"name": "special-attack"}},
			{"base_stat": 35, "stat": {"name": "special-defense"}},
			{"base_stat": 56, "stat": {"name": "speed"}}
		],
		"types": [{"type": {"name": "normal"}}, {"type": {"name": "flying"}}]
	}`), &pokemon)
	if err != nil {
		t.Fatal(err)
	}
	cfg := &config{}
	if _, err := captureCommandOutput(t, func() error {
		attemptCatch(cfg, pokemon, func(int) int { return 0 })
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	output, err := captureCommandOutput(t, func() error {
		return commandInspect(cfg, "pidgey")
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := "Name: pidgey\nHeight: 3\nWeight: 18\nStats:\n" +
		"  -hp: 40\n  -attack: 45\n  -defense: 40\n" +
		"  -special-attack: 35\n  -special-defense: 35\n  -speed: 56\n" +
		"Types:\n  - normal\n  - flying\n"
	if output != expected {
		t.Errorf("output = %q, want %q", output, expected)
	}
}

func TestCommandInspectUncaughtPokemon(t *testing.T) {
	cases := []struct {
		name    string
		pokedex map[string]pokeapi.Pokemon
	}{
		{name: "nil Pokedex"},
		{name: "empty Pokedex", pokedex: map[string]pokeapi.Pokemon{}},
		{name: "different Pokemon caught", pokedex: map[string]pokeapi.Pokemon{"pikachu": {Name: "pikachu"}}},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			cfg := &config{pokedex: testCase.pokedex}
			output, err := captureCommandOutput(t, func() error { return commandInspect(cfg, "pidgey") })
			if err != nil {
				t.Fatal(err)
			}
			if output != "you have not caught that pokemon\n" {
				t.Errorf("unexpected output: %q", output)
			}
		})
	}
}

func TestCommandInspectArguments(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{name: "missing Pokemon"},
		{name: "blank Pokemon", args: []string{" "}},
		{name: "extra arguments", args: []string{"pidgey", "extra"}},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if err := commandInspect(&config{}, testCase.args...); err == nil {
				t.Fatal("expected an argument error")
			}
		})
	}
}
