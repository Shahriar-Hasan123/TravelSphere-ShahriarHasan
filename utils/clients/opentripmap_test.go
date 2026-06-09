package clients

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockOpenTripMapServer(statusCode int, payload interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if payload != nil {
			json.NewEncoder(w).Encode(payload)
		}
	}))
}

func buildMockPlacesResponse(places []RawPlace) rawPlacesResponse {
	resp := rawPlacesResponse{}
	for _, p := range places {
		resp.Features = append(resp.Features, struct {
			Properties RawPlace `json:"properties"`
		}{Properties: p})
	}
	return resp
}

func TestFetchAttractionsByCoords_Success(t *testing.T) {
	mock := buildMockPlacesResponse([]RawPlace{
		{Name: "Eiffel Tower", Kinds: "architecture,historic", Xid: "W5013364"},
		{Name: "Louvre Museum", Kinds: "museums,historic", Xid: "W123456"},
	})

	server := mockOpenTripMapServer(200, mock)
	defer server.Close()

	client := &OpenTripMapClient{
		baseURL:    server.URL,
		apiKey:     "test-key",
		httpClient: &http.Client{},
	}

	places, err := client.FetchAttractionsByCoords(48.85, 2.29, 100000, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(places) != 2 {
		t.Errorf("expected 2 places, got %d", len(places))
	}
	if places[0].Name != "Eiffel Tower" {
		t.Errorf("expected 'Eiffel Tower', got %q", places[0].Name)
	}
}

func TestFetchAttractionsByCoords_NoAPIKey(t *testing.T) {
	client := &OpenTripMapClient{
		baseURL:    "http://example.com",
		apiKey:     "",
		httpClient: &http.Client{},
	}

	_, err := client.FetchAttractionsByCoords(48.85, 2.29, 100000, 10)
	if err == nil {
		t.Error("expected error when API key is missing")
	}
}

func TestFetchAttractionsByCoords_ServerError(t *testing.T) {
	server := mockOpenTripMapServer(500, nil)
	defer server.Close()

	client := &OpenTripMapClient{
		baseURL:    server.URL,
		apiKey:     "test-key",
		httpClient: &http.Client{},
	}

	_, err := client.FetchAttractionsByCoords(48.85, 2.29, 100000, 10)
	if err == nil {
		t.Error("expected error for 500 response")
	}
}

func TestFetchAttractionsByCoords_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := &OpenTripMapClient{
		baseURL:    server.URL,
		apiKey:     "test-key",
		httpClient: &http.Client{},
	}

	_, err := client.FetchAttractionsByCoords(48.85, 2.29, 100000, 10)
	if err == nil {
		t.Error("expected decode error")
	}
}

func TestNewOpenTripMapClient_DefaultURL(t *testing.T) {
	t.Setenv("OPENTRIPMAP_BASE_URL", "")
	t.Setenv("OPENTRIPMAP_API_KEY", "")
	client := NewOpenTripMapClient()
	if client.baseURL != openTripMapDefaultURL {
		t.Errorf("expected default URL, got %q", client.baseURL)
	}
}

func TestFetchAttractionsByCoords_RequestFailed(t *testing.T) {
	client := &OpenTripMapClient{
		baseURL:    "http://127.0.0.1:1",
		apiKey:     "test-key",
		httpClient: &http.Client{},
	}

	_, err := client.FetchAttractionsByCoords(48.85, 2.29, 100000, 10)
	if err == nil {
		t.Error("expected error for unreachable server")
	}
}

func TestNewOpenTripMapClient_CustomURL(t *testing.T) {
	t.Setenv("OPENTRIPMAP_BASE_URL", "http://custom.example.com")
	t.Setenv("OPENTRIPMAP_API_KEY", "mykey")
	client := NewOpenTripMapClient()
	if client.baseURL != "http://custom.example.com" {
		t.Errorf("expected custom URL, got %q", client.baseURL)
	}
	if client.apiKey != "mykey" {
		t.Errorf("expected apiKey 'mykey', got %q", client.apiKey)
	}
}