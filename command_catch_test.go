package main

import (
	"reflect"
	"testing"

	"github.com/AHallqvist/bdpokedex/internal/pokeapi"
)

func TestCommandCatchArguments(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{name: "missing Pokemon"},
		{name: "blank Pokemon", args: []string{" "}},
		{name: "extra arguments", args: []string{"pikachu", "extra"}},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if err := commandCatch(&config{}, testCase.args...); err == nil {
				t.Fatal("expected an argument error before making a request")
			}
		})
	}
}

func TestAttemptCatch(t *testing.T) {
	cases := []struct {
		name       string
		experience int
		roll       int
		caught     bool
	}{
		{name: "caught", experience: 112, roll: 0, caught: true},
		{name: "escaped", experience: 112, roll: 211},
		{name: "catch threshold", experience: 112, roll: 99, caught: true},
		{name: "escape threshold", experience: 112, roll: 100},
		{name: "higher experience", experience: 340, roll: 439},
		{name: "zero experience", experience: 0, roll: 99, caught: true},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			cfg := &config{}
			pokemon := pokeapi.Pokemon{Name: "pikachu", BaseExperience: testCase.experience}
			rolls := 0
			output, err := captureCommandOutput(t, func() error {
				attemptCatch(cfg, pokemon, func(bound int) int {
					rolls++
					if bound != 100+testCase.experience {
						t.Errorf("roll bound = %d, want %d", bound, 100+testCase.experience)
					}
					return testCase.roll
				})
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
			expected := "Throwing a Pokeball at pikachu...\npikachu escaped!\n"
			if testCase.caught {
				expected = "Throwing a Pokeball at pikachu...\npikachu was caught!\n"
			}
			if output != expected || rolls != 1 {
				t.Errorf("output = %q, rolls = %d; want %q and 1", output, rolls, expected)
			}
			stored, caught := cfg.pokedex["pikachu"]
			if caught != testCase.caught || (caught && !reflect.DeepEqual(stored, pokemon)) {
				t.Errorf("unexpected Pokedex state: %+v", cfg.pokedex)
			}
		})
	}
}

func TestEscapePreservesCaughtPokemon(t *testing.T) {
	pokemon := pokeapi.Pokemon{Name: "pikachu", BaseExperience: 112}
	cfg := &config{pokedex: map[string]pokeapi.Pokemon{"pikachu": pokemon}}
	_, err := captureCommandOutput(t, func() error {
		attemptCatch(cfg, pokemon, func(bound int) int { return bound - 1 })
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if stored, ok := cfg.pokedex["pikachu"]; !ok || !reflect.DeepEqual(stored, pokemon) {
		t.Fatal("an escaped retry removed an already caught Pokemon")
	}
}
