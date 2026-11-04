package main

import (
	"testing"

	"github.com/AHallqvist/bdpokedex/internal/pokeapi"
)

func TestCommandPokedex(t *testing.T) {
	cases := []struct {
		name     string
		pokedex  map[string]pokeapi.Pokemon
		expected string
	}{
		{name: "nil Pokedex", expected: "Your Pokedex:\n"},
		{name: "empty Pokedex", pokedex: map[string]pokeapi.Pokemon{}, expected: "Your Pokedex:\n"},
		{
			name:     "one Pokemon",
			pokedex:  map[string]pokeapi.Pokemon{"pidgey": {Name: "pidgey"}},
			expected: "Your Pokedex:\n - pidgey\n",
		},
		{
			name: "multiple Pokemon sorted",
			pokedex: map[string]pokeapi.Pokemon{
				"pidgey":   {Name: "pidgey"},
				"caterpie": {Name: "caterpie"},
			},
			expected: "Your Pokedex:\n - caterpie\n - pidgey\n",
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			cfg := &config{pokedex: testCase.pokedex}
			output, err := captureCommandOutput(t, func() error { return commandPokedex(cfg) })
			if err != nil {
				t.Fatal(err)
			}
			if output != testCase.expected {
				t.Errorf("output = %q, want %q", output, testCase.expected)
			}
		})
	}
}

func TestCommandPokedexRejectsArguments(t *testing.T) {
	output, err := captureCommandOutput(t, func() error { return commandPokedex(&config{}, "pidgey") })
	if err == nil {
		t.Fatal("expected an argument error")
	}
	if output != "" {
		t.Errorf("unexpected output for invalid arguments: %q", output)
	}
}
