package usecase

import (
	"context"
	"testing"

	"github.com/gilbertom/go-temperatura-cep/internal/entity"
	"go.opentelemetry.io/otel/trace"
)

// Mock Weather Repository
type MockWeatherRepository struct{}

func (m *MockWeatherRepository) GetTemperaturesByLocality(ctx context.Context, locality string, tracer trace.Tracer) (*entity.Weather, error) {
	return &entity.Weather{
		Current: struct {
			TempC float64 `json:"temp_c"`
			TempF float64 `json:"temp_f"`
		}{TempC: 25.0, TempF: 77.0},
	}, nil
}

func TestGetWeatherByLocality(t *testing.T) {
	mockRepo := &MockWeatherRepository{}
	weatherUsecase := NewWeatherUsecase(mockRepo)

	weather, err := weatherUsecase.GetTemperaturesByLocality(nil, "São Paulo", nil)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if weather.Current.TempC != 25.0 {
		t.Fatalf("expected temperature in Celsius: 25.0, got: %v", weather.Current.TempC)
	}

	if weather.Current.TempF != 77.0 {
		t.Fatalf("expected temperature in Fahrenheit: 77.0, got: %v", weather.Current.TempF)
	}
}

func TestConvertCelsiusToFahrenheit(t *testing.T) {
	weatherUsecase := NewWeatherUsecase(nil)

	celsius := 25.0
	expectedFahrenheit := 77.0
	result := weatherUsecase.ConvertCelsiusToFahrenheit(celsius)

	if result != expectedFahrenheit {
		t.Fatalf("expected %v, got %v", expectedFahrenheit, result)
	}
}

func TestConvertCelsiusToKelvin(t *testing.T) {
	weatherUsecase := NewWeatherUsecase(nil)

	celsius := 25.0
	expectedKelvin := 298.0
	result := weatherUsecase.ConvertCelsiusToKelvin(celsius)

	if result != expectedKelvin {
		t.Fatalf("expected %v, got %v", expectedKelvin, result)
	}
}
