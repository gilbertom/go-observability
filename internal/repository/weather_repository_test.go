package repository

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gilbertom/go-temperatura-cep/internal/config"
)

// Mock HTTP server
func mockWeatherHTTPServerResponse(response string, statusCode int) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	})
	return httptest.NewServer(handler)
}

// Test function for successful API call
func TestGetTemperaturesByLocalitySuccess(t *testing.T) {
	originalURL := config.AppConfig.URLWeather
	defer func() { config.AppConfig.URLWeather = originalURL }()

	successResponse := `{
		"location": {"name": "São Paulo"},
		"current": {"temp_c": 25.0}
	}`

	server := mockWeatherHTTPServerResponse(successResponse, http.StatusOK)
	defer server.Close()

	config.AppConfig.URLWeather = server.URL
	config.AppConfig.APIKeyWeather = "valid-api-key"

	repo := NewWeatherRepository()
	weather, err := repo.GetTemperaturesByLocality(nil, "São Paulo", nil)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if weather.Current.TempC != 25.0 {
		t.Fatalf("expected temperature: 25.0, got: %v", weather.Current.TempC)
	}
}

// Test function for empty API key
func TestGetTemperaturesByLocalityEmptyAPIKey(t *testing.T) {
	originalURL := config.AppConfig.URLWeather
	defer func() { config.AppConfig.URLWeather = originalURL }()

	server := mockWeatherHTTPServerResponse(``, http.StatusInternalServerError)

	config.AppConfig.URLWeather = server.URL
	config.AppConfig.APIKeyWeather = ""

	repo := NewWeatherRepository()
	_, err := repo.GetTemperaturesByLocality(nil, "São Paulo", nil)
	if err == nil {
		t.Fatalf("expected an error due to empty API key, got nil")
	}

	if err.Error() != "failed to fetch weather data" {
		t.Fatalf("expected error %s, got %s", "failed to fetch weather data", err.Error())
	}
}

// Test function for invalid locality
func TestGetTemperaturesByLocalityInvalidLocality(t *testing.T) {
	originalURL := config.AppConfig.URLWeather
	defer func() { config.AppConfig.URLWeather = originalURL }()

	errorResponse := `{
		"error": {
			"code": 1006,
			"message": "No matching location found."
		}
	}`

	server := mockWeatherHTTPServerResponse(errorResponse, http.StatusBadRequest)
	defer server.Close()

	config.AppConfig.URLWeather = server.URL
	config.AppConfig.APIKeyWeather = "valid-api-key"

	repo := NewWeatherRepository()
	_, err := repo.GetTemperaturesByLocality(nil, "InvalidLocality", nil)
	if err == nil {
		t.Fatalf("expected an error due to invalid locality, got nil")
	}

	if err.Error() != "failed to fetch weather data" {
		t.Fatalf("expected error %s, got %s", "failed to fetch weather data", err.Error())
	}
}

// Test function for network error
func TestGetTemperaturesByLocalityNetworkError(t *testing.T) {
	originalURL := config.AppConfig.URLWeather
	defer func() { config.AppConfig.URLWeather = originalURL }()

	server := mockWeatherHTTPServerResponse(``, http.StatusInternalServerError)
	defer server.Close()

	config.AppConfig.URLWeather = server.URL
	config.AppConfig.APIKeyWeather = "valid-api-key"

	repo := NewWeatherRepository()
	_, err := repo.GetTemperaturesByLocality(nil, "São Paulo", nil)
	if err == nil {
		t.Fatalf("expected a network error, got nil")
	}

	if err.Error() != "failed to fetch weather data" {
		t.Fatalf("expected error %s, got %s", "failed to fetch weather data", err.Error())
	}
}
