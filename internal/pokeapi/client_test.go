package pokeapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func TestListLocationAreas(t *testing.T) {
	var requests atomic.Int64
	var baseURL string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests.Add(1)
		if request.URL.Path != "/location-area/" || request.URL.Query().Get("limit") != "20" {
			t.Errorf("unexpected request: %s", request.URL)
		}
		firstURL := baseURL + "/location-area/?offset=0&limit=20"
		nextURL := baseURL + "/location-area/?offset=20&limit=20"
		page := LocationAreasResponse{Count: 40}
		offset := 0
		if request.URL.Query().Get("offset") == "20" {
			offset = 20
			page.Previous = &firstURL
		} else {
			page.Next = &nextURL
		}
		for index := 0; index < 20; index++ {
			page.Results = append(page.Results, LocationArea{Name: "area-" + strconv.Itoa(offset+index)})
		}
		if err := json.NewEncoder(writer).Encode(page); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	defer server.Close()
	baseURL = server.URL
	client := NewClient(time.Second)
	client.baseURL = baseURL

	first, err := client.ListLocationAreas(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Results) != 20 || first.Results[0].Name != "area-0" || first.Next == nil || first.Previous != nil {
		t.Fatalf("unexpected first page: %+v", first)
	}
	second, err := client.ListLocationAreas(first.Next)
	if err != nil {
		t.Fatal(err)
	}
	if len(second.Results) != 20 || second.Results[0].Name != "area-20" || second.Next != nil || second.Previous == nil {
		t.Fatalf("unexpected second page: %+v", second)
	}
	previous, err := client.ListLocationAreas(second.Previous)
	if err != nil {
		t.Fatal(err)
	}
	if len(previous.Results) != 20 || previous.Results[0].Name != "area-0" || previous.Previous != nil {
		t.Fatalf("unexpected previous page: %+v", previous)
	}
	if count := requests.Load(); count != 2 {
		t.Errorf("made %d requests, want 2; previous page should be cached", count)
	}
}

func TestListLocationAreasErrors(t *testing.T) {
	cases := []struct {
		name   string
		status int
		body   string
	}{
		{name: "HTTP error", status: http.StatusInternalServerError, body: "unavailable"},
		{name: "invalid JSON", status: http.StatusOK, body: "not JSON"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			var requests atomic.Int64
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				requests.Add(1)
				writer.WriteHeader(testCase.status)
				if _, err := writer.Write([]byte(testCase.body)); err != nil {
					t.Errorf("write response: %v", err)
				}
			}))
			defer server.Close()
			client := NewClient(time.Second)
			client.baseURL = server.URL
			for attempt := 0; attempt < 2; attempt++ {
				if _, err := client.ListLocationAreas(nil); err == nil {
					t.Fatal("expected an error")
				}
			}
			if count := requests.Load(); count != 2 {
				t.Errorf("made %d requests, want 2; failed responses must not be cached", count)
			}
		})
	}
}
