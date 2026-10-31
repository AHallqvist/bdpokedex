package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{Timeout: timeout},
		baseURL:    "https://pokeapi.co/api/v2",
		cache:      pokecache.NewCache(5 * time.Minute),
	}
}

func (client Client) ListLocationAreas(pageURL *string) (LocationAreasResponse, error) {
	url := client.baseURL + "/location-area/?offset=0&limit=20"
	if pageURL != nil {
		url = *pageURL
	}
	if data, ok := client.cache.Get(url); ok {
		var page LocationAreasResponse
		if err := json.Unmarshal(data, &page); err != nil {
			return LocationAreasResponse{}, fmt.Errorf("decode cached location areas: %w", err)
		}
		return page, nil
	}
	response, err := client.httpClient.Get(url)
	if err != nil {
		return LocationAreasResponse{}, fmt.Errorf("fetch location areas: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return LocationAreasResponse{}, fmt.Errorf("fetch location areas: %s", response.Status)
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return LocationAreasResponse{}, fmt.Errorf("read location areas: %w", err)
	}
	var page LocationAreasResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return LocationAreasResponse{}, fmt.Errorf("decode location areas: %w", err)
	}
	client.cache.Add(url, data)
	return page, nil
}
