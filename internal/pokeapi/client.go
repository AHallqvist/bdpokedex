package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	httpClient http.Client
	baseURL    string
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
	}
}

func (client Client) ListLocationAreas(pageURL *string) (LocationAreasResponse, error) {
	url := client.baseURL + "/location-area/?limit=20"
	if pageURL != nil {
		url = *pageURL
	}
	response, err := client.httpClient.Get(url)
	if err != nil {
		return LocationAreasResponse{}, fmt.Errorf("fetch location areas: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return LocationAreasResponse{}, fmt.Errorf("fetch location areas: %s", response.Status)
	}
	var page LocationAreasResponse
	if err := json.NewDecoder(response.Body).Decode(&page); err != nil {
		return LocationAreasResponse{}, fmt.Errorf("decode location areas: %w", err)
	}
	return page, nil
}
