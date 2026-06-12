// WeatherDTO holds current conditions, tomorrow's forecast, and a derived travel advice string - all displayed on the destination page.
package models

// ForecastDay holds weather data for a single forecast day.
type ForecastDay struct {
	Date      string
	MaxTempC  float64
	MinTempC  float64
	MaxTempF  float64
	MinTempF  float64
	Condition string
	Icon      string
}

// Weather holds the full weather block shown on the destination page.
type Weather struct {
	// Current conditions
	TempC      float64
	TempF      float64
	FeelsLikeC float64
	Humidity   int
	WindKph    float64
	Condition  string
	Icon       string
	City       string

	// Tomorrow's forecast
	Forecast *ForecastDay

	// Human-readable travel advice derived from current conditions
	TravelAdvice string
}
