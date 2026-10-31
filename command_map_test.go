package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/AHallqvist/bdpokedex/internal/pokeapi"
)

func TestMapPagination(t *testing.T) {
	requests := 0
	var baseURL string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests++
		offset, err := strconv.Atoi(request.URL.Query().Get("offset"))
		if err != nil {
			t.Errorf("invalid offset: %v", err)
		}
		page := pokeapi.LocationAreasResponse{Count: 60}
		if offset < 40 {
			next := baseURL + "?offset=" + strconv.Itoa(offset+20)
			page.Next = &next
		}
		if offset > 0 {
			previous := baseURL + "?offset=" + strconv.Itoa(offset-20)
			page.Previous = &previous
		}
		for index := 0; index < 20; index++ {
			page.Results = append(page.Results, pokeapi.LocationArea{Name: "area-" + strconv.Itoa(offset+index)})
		}
		if err := json.NewEncoder(writer).Encode(page); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	defer server.Close()
	baseURL = server.URL
	firstURL := baseURL + "?offset=0"
	cfg := &config{pokeapiClient: pokeapi.NewClient(time.Second), Next: &firstURL}
	cases := []struct {
		name     string
		callback func(*config) error
		offset   int
		boundary string
	}{
		{name: "back before first request", callback: commandMapb, boundary: "You're on the first page\n"},
		{name: "first page", callback: commandMap, offset: 0},
		{name: "second page", callback: commandMap, offset: 20},
		{name: "last page", callback: commandMap, offset: 40},
		{name: "past last page", callback: commandMap, boundary: "You're on the last page\n"},
		{name: "back to second page", callback: commandMapb, offset: 20},
		{name: "back to first page", callback: commandMapb, offset: 0},
		{name: "before first page", callback: commandMapb, boundary: "You're on the first page\n"},
		{name: "forward again", callback: commandMap, offset: 20},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			before := requests
			output, err := captureCommandOutput(t, func() error { return testCase.callback(cfg) })
			if err != nil {
				t.Fatal(err)
			}
			if testCase.boundary != "" {
				if output != testCase.boundary || requests != before {
					t.Fatalf("boundary output = %q, requests = %d; want %q and %d", output, requests, testCase.boundary, before)
				}
				return
			}
			var expected strings.Builder
			for index := 0; index < 20; index++ {
				expected.WriteString("area-" + strconv.Itoa(testCase.offset+index) + "\n")
			}
			if output != expected.String() {
				t.Errorf("output = %q, want %q", output, expected.String())
			}
			if (cfg.Previous == nil) != (testCase.offset == 0) || (cfg.Next == nil) != (testCase.offset == 40) {
				t.Errorf("unexpected pagination state for offset %d", testCase.offset)
			}
		})
	}
}

func TestMapFailurePreservesPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()
	url := server.URL
	cfg := &config{
		pokeapiClient:       pokeapi.NewClient(time.Second),
		Next:                &url,
		Previous:            &url,
		locationAreasLoaded: true,
	}
	for _, callback := range []func(*config) error{commandMap, commandMapb} {
		output, err := captureCommandOutput(t, func() error { return callback(cfg) })
		if err == nil {
			t.Fatal("expected an API error")
		}
		if output != "" || cfg.Next != &url || cfg.Previous != &url || !cfg.locationAreasLoaded {
			t.Fatal("failed request changed pagination state or printed locations")
		}
	}
}

func captureCommandOutput(t *testing.T, callback func() error) (string, error) {
	t.Helper()
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	defer writer.Close()
	stdout := os.Stdout
	os.Stdout = writer
	defer func() { os.Stdout = stdout }()
	commandErr := callback()
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	return string(output), commandErr
}
