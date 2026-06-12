package clients

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockWeatherServer(statusCode int, payload interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if payload != nil {
			json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestFetchForecast_NoAPIKey(t *testing.T) {
	client := &WeatherAPIClient{
		baseURL:    "http://example.com",
		apiKey:     "",
		httpClient: &http.Client{},
	}

	result, err := client.FetchForecast("Dhaka")
	if err != nil {
		t.Fatalf("expected nil error when key absent, got %v", err)
	}
	if result != nil {
		t.Error("expected nil result when key absent")
	}
}

func TestFetchForecast_Success(t *testing.T) {
	mock := RawWeather{}
	mock.Location.Name = "Dhaka"
	mock.Current.TempC = 32.5
	mock.Current.TempF = 90.5
	mock.Current.Condition.Text = "Sunny"
	mock.Current.Condition.Icon = "//cdn.weatherapi.com/sunny.png"

	server := mockWeatherServer(200, mock)
	defer server.Close()

	client := &WeatherAPIClient{
		baseURL:    server.URL,
		apiKey:     "test-key",
		httpClient: &http.Client{},
	}

	result, err := client.FetchForecast("Dhaka")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Location.Name != "Dhaka" {
		t.Errorf("expected 'Dhaka', got %q", result.Location.Name)
	}
	if result.Current.TempC != 32.5 {
		t.Errorf("expected 32.5, got %f", result.Current.TempC)
	}
}

func TestFetchForecast_CityNotFound(t *testing.T) {
	server := mockWeatherServer(400, nil)
	defer server.Close()

	client := &WeatherAPIClient{
		baseURL:    server.URL,
		apiKey:     "test-key",
		httpClient: &http.Client{},
	}

	result, err := client.FetchForecast("InvalidCity")
	if err != nil {
		t.Fatalf("expected nil error for 400, got %v", err)
	}
	if result != nil {
		t.Error("expected nil result for 400")
	}
}

func TestFetchForecast_ServerError(t *testing.T) {
	server := mockWeatherServer(500, nil)
	defer server.Close()

	client := &WeatherAPIClient{
		baseURL:    server.URL,
		apiKey:     "test-key",
		httpClient: &http.Client{},
	}

	_, err := client.FetchForecast("Dhaka")
	if err == nil {
		t.Error("expected error for 500 response")
	}
}

func TestFetchForecast_NotFoundStatus(t *testing.T) {
	server := mockWeatherServer(404, nil)
	defer server.Close()

	client := &WeatherAPIClient{
		baseURL:    server.URL,
		apiKey:     "test-key",
		httpClient: &http.Client{},
	}

	result, err := client.FetchForecast("Unknown")
	if err != nil {
		t.Fatalf("expected nil error for 404, got %v", err)
	}
	if result != nil {
		t.Error("expected nil result for 404")
	}
}

func TestFetchForecast_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("bad json"))
	}))
	defer server.Close()

	client := &WeatherAPIClient{
		baseURL:    server.URL,
		apiKey:     "test-key",
		httpClient: &http.Client{},
	}

	_, err := client.FetchForecast("Dhaka")
	if err == nil {
		t.Error("expected decode error")
	}
}

func TestIsConfigured(t *testing.T) {
	c1 := &WeatherAPIClient{apiKey: "abc"}
	if !c1.IsConfigured() {
		t.Error("expected IsConfigured true")
	}

	c2 := &WeatherAPIClient{apiKey: ""}
	if c2.IsConfigured() {
		t.Error("expected IsConfigured false")
	}
}

func TestNewWeatherAPIClient_DefaultURL(t *testing.T) {
	t.Setenv("WEATHERAPI_BASE_URL", "")
	t.Setenv("WEATHERAPI_KEY", "")
	client := NewWeatherAPIClient()
	if client.baseURL != weatherAPIDefaultURL {
		t.Errorf("expected default URL, got %q", client.baseURL)
	}
}

func TestFetchForecast_RequestFailed(t *testing.T) {
	client := &WeatherAPIClient{
		baseURL:    "http://127.0.0.1:1",
		apiKey:     "test-key",
		httpClient: &http.Client{},
	}

	_, err := client.FetchForecast("Dhaka")
	if err == nil {
		t.Error("expected error for unreachable server")
	}
}

func TestNewWeatherAPIClient_CustomURL(t *testing.T) {
	t.Setenv("WEATHERAPI_BASE_URL", "http://weather.example.com")
	t.Setenv("WEATHERAPI_KEY", "myweatherkey")
	client := NewWeatherAPIClient()
	if client.baseURL != "http://weather.example.com" {
		t.Errorf("expected custom URL, got %q", client.baseURL)
	}
}
