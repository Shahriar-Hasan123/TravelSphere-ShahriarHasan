// WeatherService retrieves current weather and forecast for a capital city, and derives a human-readable travel conditions summary.
package services

import (
	"TravelSphere/models"
	"TravelSphere/utils/clients"
	"fmt"
	"strings"
)

// WeatherService provides weather data and travel condition advice.
type WeatherService struct {
	client *clients.WeatherAPIClient
}

// NewWeatherService creates a WeatherService with its API client.
func NewWeatherService() *WeatherService {
	return &WeatherService{
		client: clients.NewWeatherAPIClient(),
	}
}

// IsConfigured reports whether the WeatherAPI key is set.
func (s *WeatherService) IsConfigured() bool {
	return s.client.IsConfigured()
}

// GetWeather fetches current weather and forecast for the given city.
func (s *WeatherService) GetWeather(city string) (*models.Weather, error) {
	if city == "" {
		return nil, nil
	}

	raw, err := s.client.FetchForecast(city)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	weather := &models.Weather{
		TempC:      raw.Current.TempC,
		TempF:      raw.Current.TempF,
		FeelsLikeC: raw.Current.FeelsLikeC,
		Humidity:   raw.Current.Humidity,
		WindKph:    raw.Current.WindKph,
		Condition:  raw.Current.Condition.Text,
		Icon:       "https:" + raw.Current.Condition.Icon,
		City:       raw.Location.Name,
	}

	// Populate tomorrow's forecast if the API returned it.
	if len(raw.Forecast.ForecastDay) > 1 {
		tomorrow := raw.Forecast.ForecastDay[1]
		weather.Forecast = &models.ForecastDay{
			Date:      tomorrow.Date,
			MaxTempC:  tomorrow.Day.MaxTempC,
			MinTempC:  tomorrow.Day.MinTempC,
			MaxTempF:  tomorrow.Day.MaxTempF,
			MinTempF:  tomorrow.Day.MinTempF,
			Condition: tomorrow.Day.Condition.Text,
			Icon:      "https:" + tomorrow.Day.Condition.Icon,
		}
	}

	// Derive a travel advice string from current conditions.
	weather.TravelAdvice = deriveTravelAdvice(raw)

	return weather, nil
}

// deriveTravelAdvice builds a human-readable travel conditions summary based on temperature, wind, humidity, and the condition text.
func deriveTravelAdvice(raw *clients.RawWeather) string {
	cond := strings.ToLower(raw.Current.Condition.Text)
	temp := raw.Current.TempC
	wind := raw.Current.WindKph
	hum := raw.Current.Humidity

	var advice []string

	// Temperature advice.
	switch {
	case temp >= 30:
		advice = append(advice, "Very hot — stay hydrated and use sun protection.")
	case temp >= 20:
		advice = append(advice, "Warm and pleasant — great for sightseeing.")
	case temp >= 10:
		advice = append(advice, "Mild — a light jacket is recommended.")
	case temp >= 0:
		advice = append(advice, "Cold — dress in warm layers.")
	default:
		advice = append(advice, "Freezing conditions — heavy winter clothing required.")
	}

	// Precipitation advice.
	switch {
	case strings.Contains(cond, "thunder"):
		advice = append(advice, "Thunderstorms expected — avoid outdoor activities.")
	case strings.Contains(cond, "snow") || strings.Contains(cond, "blizzard"):
		advice = append(advice, "Snowfall likely — check road conditions before travelling.")
	case strings.Contains(cond, "rain") || strings.Contains(cond, "drizzle"):
		advice = append(advice, "Carry an umbrella.")
	case strings.Contains(cond, "fog") || strings.Contains(cond, "mist"):
		advice = append(advice, "Reduced visibility — allow extra travel time.")
	}

	// Wind advice.
	if wind > 50 {
		advice = append(advice, fmt.Sprintf("Strong winds at %.0f km/h — secure loose items.", wind))
	}

	// Humidity advice.
	if hum > 80 {
		advice = append(advice, "High humidity — light, breathable clothing advised.")
	}

	if len(advice) == 0 {
		return "Conditions look good for travel."
	}
	return strings.Join(advice, " ")
}
