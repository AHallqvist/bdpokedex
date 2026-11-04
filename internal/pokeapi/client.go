package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/AHallqvist/bdpokedex/internal/pokecache"
)

type Client struct {
	httpClient http.Client
	baseURL    string
	cache      *pokecache.Cache
}

type LocationAreasResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreaDetails struct {
	Name              string             `json:"name"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon struct {
		Name string `json:"name"`
	} `json:"pokemon"`
}

type Pokemon struct {
	Name           string        `json:"name"`
	BaseExperience int           `json:"base_experience"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	Stats          []PokemonStat `json:"stats"`
	Types          []PokemonType `json:"types"`
}

type PokemonStat struct {
	BaseStat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type PokemonType struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{Timeout: timeout},
		baseURL:    "https://pokeapi.co/api/v2",
		cache:      pokecache.NewCache(5 * time.Minute),
	}
}

func (client Client) ListLocationAreas(pageURL *string) (LocationAreasResponse, error) {
	requestURL := client.baseURL + "/location-area/?offset=0&limit=20"
	if pageURL != nil {
		requestURL = *pageURL
	}
	var page LocationAreasResponse
	if err := client.getJSON(requestURL, &page); err != nil {
		return LocationAreasResponse{}, err
	}
	return page, nil
}

func (client Client) GetLocationArea(name string) (LocationAreaDetails, error) {
	requestURL := client.baseURL + "/location-area/" + url.PathEscape(name) + "/"
	var area LocationAreaDetails
	if err := client.getJSON(requestURL, &area); err != nil {
		return LocationAreaDetails{}, err
	}
	return area, nil
}

func (client Client) GetPokemon(name string) (Pokemon, error) {
	requestURL := client.baseURL + "/pokemon/" + url.PathEscape(name) + "/"
	var pokemon Pokemon
	if err := client.getJSON(requestURL, &pokemon); err != nil {
		return Pokemon{}, err
	}
	return pokemon, nil
}

func (client Client) getJSON(requestURL string, target any) error {
	if data, ok := client.cache.Get(requestURL); ok {
		if err := json.Unmarshal(data, target); err != nil {
			return fmt.Errorf("decode cached PokeAPI response: %w", err)
		}
		return nil
	}
	response, err := client.httpClient.Get(requestURL)
	if err != nil {
		return fmt.Errorf("fetch PokeAPI response: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("fetch PokeAPI response: %s", response.Status)
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read PokeAPI response: %w", err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("decode PokeAPI response: %w", err)
	}
	client.cache.Add(requestURL, data)
	return nil
}
