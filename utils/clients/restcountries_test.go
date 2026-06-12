package clients

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockRestCountriesServer(statusCode int, payload interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(payload)
	}))
}

func TestFetchAll_Success(t *testing.T) {
	mock := []RawCountry{
		{Population: 2837743, Region: "Europe"},
	}
	mock[0].Name.Common = "Albania"
	mock[0].Name.Official = "Republic of Albania"

	server := mockRestCountriesServer(200, mock)
	defer server.Close()

	client := &RestCountriesClient{
		baseURL:    server.URL,
		httpClient: &http.Client{},
	}

	countries, err := client.FetchAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(countries) != 1 {
		t.Errorf("expected 1 country, got %d", len(countries))
	}
	if countries[0].Name.Common != "Albania" {
		t.Errorf("expected Albania, got %q", countries[0].Name.Common)
	}
}

func TestFetchAll_ServerError(t *testing.T) {
	server := mockRestCountriesServer(500, nil)
	defer server.Close()

	client := &RestCountriesClient{
		baseURL:    server.URL,
		httpClient: &http.Client{},
	}

	_, err := client.FetchAll()
	if err == nil {
		t.Error("expected error for 500 response, got nil")
	}
}

func TestFetchAll_NotFound(t *testing.T) {
	server := mockRestCountriesServer(404, nil)
	defer server.Close()

	client := &RestCountriesClient{
		baseURL:    server.URL,
		httpClient: &http.Client{},
	}

	countries, err := client.FetchAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if countries != nil {
		t.Errorf("expected nil for 404, got %v", countries)
	}
}

func TestFetchAll_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("not valid json"))
	}))
	defer server.Close()

	client := &RestCountriesClient{
		baseURL:    server.URL,
		httpClient: &http.Client{},
	}

	_, err := client.FetchAll()
	if err == nil {
		t.Error("expected decode error, got nil")
	}
}

func TestFetchByName_Success(t *testing.T) {
	mock := []RawCountry{{Population: 67000000, Region: "Europe"}}
	mock[0].Name.Common = "France"

	server := mockRestCountriesServer(200, mock)
	defer server.Close()

	client := &RestCountriesClient{
		baseURL:    server.URL,
		httpClient: &http.Client{},
	}

	countries, err := client.FetchByName("france")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(countries) != 1 || countries[0].Name.Common != "France" {
		t.Error("expected France in results")
	}
}

func TestNewRestCountriesClient_DefaultURL(t *testing.T) {
	t.Setenv("RESTCOUNTRIES_BASE_URL", "")
	client := NewRestCountriesClient()
	if client.baseURL != restCountriesDefaultURL {
		t.Errorf("expected default URL, got %q", client.baseURL)
	}
}

func TestNewRestCountriesClient_CustomURL(t *testing.T) {
	t.Setenv("RESTCOUNTRIES_BASE_URL", "http://custom.example.com")
	client := NewRestCountriesClient()
	if client.baseURL != "http://custom.example.com" {
		t.Errorf("expected custom URL, got %q", client.baseURL)
	}
}

func TestFetchAll_RequestFailed(t *testing.T) {
	// Point client at an invalid address to force a connection error.
	client := &RestCountriesClient{
		baseURL:    "http://127.0.0.1:1",
		httpClient: &http.Client{},
	}

	_, err := client.FetchAll()
	if err == nil {
		t.Error("expected error for unreachable server")
	}
}

func TestFetchByName_NotFound(t *testing.T) {
	server := mockRestCountriesServer(404, nil)
	defer server.Close()

	client := &RestCountriesClient{
		baseURL:    server.URL,
		httpClient: &http.Client{},
	}

	countries, err := client.FetchByName("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if countries != nil {
		t.Error("expected nil for 404 response")
	}
}

func TestFetchByName_ServerError(t *testing.T) {
	server := mockRestCountriesServer(500, nil)
	defer server.Close()

	client := &RestCountriesClient{
		baseURL:    server.URL,
		httpClient: &http.Client{},
	}

	_, err := client.FetchByName("france")
	if err == nil {
		t.Error("expected error for 500 response")
	}
}
