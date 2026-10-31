package main

import "fmt"

func commandMap(cfg *config) error {
	if cfg.locationAreasLoaded && cfg.Next == nil {
		fmt.Println("You're on the last page")
		return nil
	}
	return displayLocationAreas(cfg, cfg.Next)
}

func commandMapb(cfg *config) error {
	if cfg.Previous == nil {
		fmt.Println("You're on the first page")
		return nil
	}
	return displayLocationAreas(cfg, cfg.Previous)
}

func displayLocationAreas(cfg *config, pageURL *string) error {
	page, err := cfg.pokeapiClient.ListLocationAreas(pageURL)
	if err != nil {
		return err
	}
	cfg.Next = page.Next
	cfg.Previous = page.Previous
	cfg.locationAreasLoaded = true
	for _, area := range page.Results {
		fmt.Println(area.Name)
	}
	return nil
}
