package main

import (
	"time"

	"github.com/AHallqvist/bdpokedex/internal/pokeapi"
)

func main() {
	cfg := &config{
		pokeapiClient: pokeapi.NewClient(10 * time.Second),
	}
	startRepl(cfg)
}
