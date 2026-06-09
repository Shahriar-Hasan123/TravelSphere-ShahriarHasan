package services

import (
	"TravelSphere/utils/clients"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func mockWeatherAPIServer(payload *clients.RawWeather) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}))
}

func TestGetWeather_NoKey(t *testing.T) {
	svc := &WeatherService{
		client: clients.NewWeatherAPIClientWithURL("http://unused", ""),
	}

	result, err := svc.GetWeather("Dhaka")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result != nil {
		t.Error("expected nil result when key absent")
	}
}

func TestGetWeather_EmptyCity(t *testing.T) {
	svc := &WeatherService{
		client: clients.NewWeatherAPIClientWithURL("http://unused", "test-key"),
	}

	result, err := svc.GetWeather("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Error("expected nil for empty city")
	}
}

func TestGetWeather_Success(t *testing.T) {
	raw := &clients.RawWeather{}
	raw.Location.Name = "Dhaka"
	raw.Current.TempC = 32.0
	raw.Current.TempF = 89.6
	raw.Current.FeelsLikeC = 38.0
	raw.Current.Humidity = 80
	raw.Current.WindKph = 15.0
	raw.Current.Condition.Text = "Partly cloudy"
	raw.Current.Condition.Icon = "//cdn.icon.png"

	server := mockWeatherAPIServer(raw)
	defer server.Close()

	svc := &WeatherService{
		client: clients.NewWeatherAPIClientWithURL(server.URL, "test-key"),
	}

	weather, err := svc.GetWeather("Dhaka")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if weather == nil {
		t.Fatal("expected weather data, got nil")
	}
	if weather.TempC != 32.0 {
		t.Errorf("expected TempC 32.0, got %f", weather.TempC)
	}
	if weather.City != "Dhaka" {
		t.Errorf("expected city Dhaka, got %q", weather.City)
	}
	if weather.TravelAdvice == "" {
		t.Error("expected non-empty travel advice")
	}
}

func TestGetWeather_WithForecast(t *testing.T) {
	raw := &clients.RawWeather{}
	raw.Location.Name = "Paris"
	raw.Current.TempC = 22.0
	raw.Current.Condition.Text = "Sunny"
	raw.Current.Condition.Icon = "//icon.png"

	// Add tomorrow's forecast.
	raw.Forecast.ForecastDay = append(raw.Forecast.ForecastDay,
		struct {
			Date string `json:"date"`
			Day  struct {
				MaxTempC    float64 `json:"maxtemp_c"`
				MinTempC    float64 `json:"mintemp_c"`
				MaxTempF    float64 `json:"maxtemp_f"`
				MinTempF    float64 `json:"mintemp_f"`
				AvgHumidity float64 `json:"avghumidity"`
				Condition   struct {
					Text string `json:"text"`
					Icon string `json:"icon"`
				} `json:"condition"`
			} `json:"day"`
		}{
			Date: "2026-06-09",
		},
	)
	raw.Forecast.ForecastDay = append(raw.Forecast.ForecastDay,
		raw.Forecast.ForecastDay[0],
	)
	raw.Forecast.ForecastDay[1].Date = "2026-06-10"
	raw.Forecast.ForecastDay[1].Day.MaxTempC = 24.0
	raw.Forecast.ForecastDay[1].Day.MinTempC = 16.0
	raw.Forecast.ForecastDay[1].Day.Condition.Text = "Cloudy"

	server := mockWeatherAPIServer(raw)
	defer server.Close()

	svc := &WeatherService{
		client: clients.NewWeatherAPIClientWithURL(server.URL, "test-key"),
	}

	weather, err := svc.GetWeather("Paris")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if weather.Forecast == nil {
		t.Error("expected forecast data")
	}
}

func TestIsConfigured(t *testing.T) {
	svc1 := &WeatherService{client: clients.NewWeatherAPIClientWithURL("", "key")}
	if !svc1.IsConfigured() {
		t.Error("expected IsConfigured true")
	}

	svc2 := &WeatherService{client: clients.NewWeatherAPIClientWithURL("", "")}
	if svc2.IsConfigured() {
		t.Error("expected IsConfigured false")
	}
}

func TestGetWeather_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer server.Close()

	svc := &WeatherService{
		client: clients.NewWeatherAPIClientWithURL(server.URL, "test-key"),
	}

	_, err := svc.GetWeather("Dhaka")
	if err == nil {
		t.Error("expected error from API failure")
	}
}

func TestGetWeather_NilRawResponse(t *testing.T) {
	// 400 from WeatherAPI means city not found — client returns nil,nil.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer server.Close()

	svc := &WeatherService{
		client: clients.NewWeatherAPIClientWithURL(server.URL, "test-key"),
	}

	result, err := svc.GetWeather("UnknownCity")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Error("expected nil result for unknown city")
	}
}

func TestNewWeatherService_NotNil(t *testing.T) {
	svc := NewWeatherService()
	if svc == nil {
		t.Error("expected non-nil WeatherService")
	}
}

// Tests for deriveTravelAdvice edge cases
func TestDeriveTravelAdvice_VeryHot(t *testing.T) {
	raw := &clients.RawWeather{}
	raw.Current.TempC = 35.0
	raw.Current.Condition.Text = "Sunny"
	raw.Current.WindKph = 5.0
	raw.Current.Humidity = 50.0

	advice := deriveTravelAdvice(raw)
	if !contains(advice, "Very hot") {
		t.Errorf("expected 'Very hot' in advice for 35°C, got: %q", advice)
	}
}

func TestDeriveTravelAdvice_ColdFreezingBelowZero(t *testing.T) {
	raw := &clients.RawWeather{}
	raw.Current.TempC = -5.0
	raw.Current.Condition.Text = "Snow"
	raw.Current.WindKph = 10.0
	raw.Current.Humidity = 60.0

	advice := deriveTravelAdvice(raw)
	if !contains(advice, "Freezing") {
		t.Errorf("expected 'Freezing' in advice for -5°C, got: %q", advice)
	}
}

func TestDeriveTravelAdvice_Thunderstorm(t *testing.T) {
	raw := &clients.RawWeather{}
	raw.Current.TempC = 15.0
	raw.Current.Condition.Text = "Thunderstorm"
	raw.Current.WindKph = 10.0
	raw.Current.Humidity = 70.0

	advice := deriveTravelAdvice(raw)
	if !contains(advice, "Thunderstorms") {
		t.Errorf("expected 'Thunderstorms' in advice, got: %q", advice)
	}
}

func TestDeriveTravelAdvice_Snowfall(t *testing.T) {
	raw := &clients.RawWeather{}
	raw.Current.TempC = 0.0
	raw.Current.Condition.Text = "Heavy Snow"
	raw.Current.WindKph = 10.0
	raw.Current.Humidity = 70.0

	advice := deriveTravelAdvice(raw)
	if !contains(advice, "Snowfall") {
		t.Errorf("expected 'Snowfall' in advice, got: %q", advice)
	}
}

func TestDeriveTravelAdvice_Rain(t *testing.T) {
	raw := &clients.RawWeather{}
	raw.Current.TempC = 15.0
	raw.Current.Condition.Text = "Light Rain"
	raw.Current.WindKph = 10.0
	raw.Current.Humidity = 75.0

	advice := deriveTravelAdvice(raw)
	if !contains(advice, "umbrella") {
		t.Errorf("expected 'umbrella' in advice, got: %q", advice)
	}
}

func TestDeriveTravelAdvice_Drizzle(t *testing.T) {
	raw := &clients.RawWeather{}
	raw.Current.TempC = 12.0
	raw.Current.Condition.Text = "Drizzle"
	raw.Current.WindKph = 5.0
	raw.Current.Humidity = 80.0

	advice := deriveTravelAdvice(raw)
	if !contains(advice, "umbrella") {
		t.Errorf("expected 'umbrella' in advice for drizzle, got: %q", advice)
	}
}

func TestDeriveTravelAdvice_Fog(t *testing.T) {
	raw := &clients.RawWeather{}
	raw.Current.TempC = 10.0
	raw.Current.Condition.Text = "Fog"
	raw.Current.WindKph = 5.0
	raw.Current.Humidity = 90.0

	advice := deriveTravelAdvice(raw)
	if !contains(advice, "visibility") {
		t.Errorf("expected 'visibility' in advice for fog, got: %q", advice)
	}
}

func TestDeriveTravelAdvice_HighWind(t *testing.T) {
	raw := &clients.RawWeather{}
	raw.Current.TempC = 15.0
	raw.Current.Condition.Text = "Clear"
	raw.Current.WindKph = 60.0
	raw.Current.Humidity = 60.0

	advice := deriveTravelAdvice(raw)
	if !contains(advice, "Strong winds") {
		t.Errorf("expected 'Strong winds' in advice, got: %q", advice)
	}
}

func TestDeriveTravelAdvice_HighHumidity(t *testing.T) {
	raw := &clients.RawWeather{}
	raw.Current.TempC = 25.0
	raw.Current.Condition.Text = "Sunny"
	raw.Current.WindKph = 5.0
	raw.Current.Humidity = 85.0

	advice := deriveTravelAdvice(raw)
	if !contains(advice, "High humidity") {
		t.Errorf("expected 'High humidity' in advice, got: %q", advice)
	}
}

func TestDeriveTravelAdvice_ClearConditions(t *testing.T) {
	raw := &clients.RawWeather{}
	raw.Current.TempC = 22.0
	raw.Current.Condition.Text = "Sunny"
	raw.Current.WindKph = 5.0
	raw.Current.Humidity = 50.0

	advice := deriveTravelAdvice(raw)
	// For 22°C with sunny conditions, it returns "Warm and pleasant" advice
	if !contains(advice, "Warm and pleasant") {
		t.Errorf("expected 'Warm and pleasant' for 22°C sunny, got: %q", advice)
	}
}

func TestDeriveTravelAdvice_MinimalConditions(t *testing.T) {
	// Test with all minimal conditions - no special weather, moderate wind/humidity
	raw := &clients.RawWeather{}
	raw.Current.TempC = 20.0 // Exactly 20, hits "Warm and pleasant"
	raw.Current.Condition.Text = "Clear"
	raw.Current.WindKph = 30.0  // Below 50, so no wind advice
	raw.Current.Humidity = 60.0 // Below 80, so no humidity advice

	advice := deriveTravelAdvice(raw)
	// At 20°C with clear conditions and moderate wind/humidity, we should get the pleasant message
	if !contains(advice, "Warm and pleasant") {
		t.Errorf("expected pleasant advice, got: %q", advice)
	}
}

func TestDeriveTravelAdvice_ColdRange(t *testing.T) {
	// Test cold temperature (0-10°C) range
	raw := &clients.RawWeather{}
	raw.Current.TempC = 5.0
	raw.Current.Condition.Text = "Cloudy"
	raw.Current.WindKph = 10.0
	raw.Current.Humidity = 70.0

	advice := deriveTravelAdvice(raw)
	if !contains(advice, "Cold") || contains(advice, "Freezing") {
		t.Errorf("expected 'Cold' for 5°C, got: %q", advice)
	}
}

func TestDeriveTravelAdvice_Mist(t *testing.T) {
	raw := &clients.RawWeather{}
	raw.Current.TempC = 12.0
	raw.Current.Condition.Text = "Mist"
	raw.Current.WindKph = 5.0
	raw.Current.Humidity = 85.0

	advice := deriveTravelAdvice(raw)
	if !contains(advice, "visibility") {
		t.Errorf("expected 'visibility' in advice for mist, got: %q", advice)
	}
}

func TestDeriveTravelAdvice_Blizzard(t *testing.T) {
	raw := &clients.RawWeather{}
	raw.Current.TempC = -10.0
	raw.Current.Condition.Text = "Blizzard"
	raw.Current.WindKph = 70.0
	raw.Current.Humidity = 80.0

	advice := deriveTravelAdvice(raw)
	if !contains(advice, "Snowfall") {
		t.Errorf("expected 'Snowfall' in advice for blizzard, got: %q", advice)
	}
}

func TestDeriveTravelAdvice_MildTemperature(t *testing.T) {
	raw := &clients.RawWeather{}
	raw.Current.TempC = 15.0
	raw.Current.Condition.Text = "Partly Cloudy"
	raw.Current.WindKph = 10.0
	raw.Current.Humidity = 65.0

	advice := deriveTravelAdvice(raw)
	if !contains(advice, "Mild") {
		t.Errorf("expected 'Mild' in advice for 15°C, got: %q", advice)
	}
}

func TestDeriveTravelAdvice_WarmTemperature(t *testing.T) {
	raw := &clients.RawWeather{}
	raw.Current.TempC = 25.0
	raw.Current.Condition.Text = "Sunny"
	raw.Current.WindKph = 10.0
	raw.Current.Humidity = 50.0

	advice := deriveTravelAdvice(raw)
	if !contains(advice, "Warm and pleasant") {
		t.Errorf("expected 'Warm and pleasant' for 25°C, got: %q", advice)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
