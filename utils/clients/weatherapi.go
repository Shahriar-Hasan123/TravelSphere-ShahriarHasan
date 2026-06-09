// HTTP client for WeatherAPI (bonus feature).

package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

const weatherAPIDefaultURL = "http://api.weatherapi.com/v1"

// RawWeather holds all fields we extract from the WeatherAPI forecast response.
// We use the /forecast.json endpoint (free tier supports 1–3 days) so we get both current conditions and tomorrow's forecast in a single request.
type RawWeather struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		TempC      float64 `json:"temp_c"`
		TempF      float64 `json:"temp_f"`
		FeelsLikeC float64 `json:"feelslike_c"`
		Humidity   int     `json:"humidity"`
		WindKph    float64 `json:"wind_kph"`
		Condition  struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		ForecastDay []struct {
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
		} `json:"forecastday"`
	} `json:"forecast"`
}

// WeatherAPIClient makes requests to the WeatherAPI forecast endpoint.
type WeatherAPIClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewWeatherAPIClient creates a client reading credentials from the environment.
func NewWeatherAPIClient() *WeatherAPIClient {
	baseURL := os.Getenv("WEATHERAPI_BASE_URL")
	if baseURL == "" {
		baseURL = weatherAPIDefaultURL
	}
	return &WeatherAPIClient{
		baseURL:    baseURL,
		apiKey:     os.Getenv("WEATHERAPI_KEY"),
		httpClient: &http.Client{Timeout: 8 * time.Second},
	}
}

// NewWeatherAPIClientWithURL creates a client with a custom base URL — used in tests.
func NewWeatherAPIClientWithURL(baseURL, apiKey string) *WeatherAPIClient {
	return &WeatherAPIClient{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// IsConfigured returns true when the API key is present in the environment.
func (c *WeatherAPIClient) IsConfigured() bool {
	return c.apiKey != ""
}

// FetchForecast retrieves current weather AND a 2-day forecast for a city.
func (c *WeatherAPIClient) FetchForecast(city string) (*RawWeather, error) {
	if !c.IsConfigured() {
		return nil, nil
	}

	// days=2 gives us today + tomorrow on the free tier.
	endpoint := fmt.Sprintf(
		"%s/forecast.json?key=%s&q=%s&days=2&aqi=no&alerts=no",
		c.baseURL,
		url.QueryEscape(c.apiKey),
		url.QueryEscape(city),
	)

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("weatherapi request failed: %w", err)
	}
	defer resp.Body.Close()

	// 400 means city not found — degrade gracefully.
	if resp.StatusCode == http.StatusBadRequest ||
		resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weatherapi returned status %d", resp.StatusCode)
	}

	var raw RawWeather
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("weatherapi decode failed: %w", err)
	}
	return &raw, nil
}
