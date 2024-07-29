package dto

// WeatherResponse represents the response for weather data.
type WeatherResponse struct {
	City       string  `json:"city"`
	TempC      float64 `json:"tempC"`
	TempF      float64 `json:"tempF"`
	TempK      float64 `json:"tempK"`
}
