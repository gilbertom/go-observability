package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gilbertom/go-temperatura-cep/internal/config"
	"github.com/gilbertom/go-temperatura-cep/internal/entity"
	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

// WeatherRepository is a repository for retrieving weather data.
type WeatherRepository struct{}

// NewWeatherRepository creates a new instance of WeatherRepository.
func NewWeatherRepository() *WeatherRepository {
    return &WeatherRepository{}
}

// GetTemperaturesByLocality retrieves the temperatures for a given locality.
func (r *WeatherRepository) GetTemperaturesByLocality(ctx context.Context, locality string, tracer trace.Tracer) (*entity.Weather, error) {
    ctx, span := tracer.Start(ctx, "GetTemperaturesByLocality")
    defer span.End()

    var weather entity.Weather
    url := fmt.Sprintf("%s?q=%s&lang=pt&key=%s", config.AppConfig.URLWeather, url.QueryEscape(locality), config.AppConfig.APIKeyWeather)

    client := resty.New()
    client.SetTransport(otelhttp.NewTransport(http.DefaultTransport))
    resp, err := client.R().
        SetContext(ctx).
        SetHeader("Content-Type", "application/json").
        Get(url)
    if err != nil {
        return &weather, err
    }
    defer resp.RawBody().Close()

    if resp.StatusCode() != http.StatusOK {
        return &weather, errors.New("failed to fetch weather data")
    }

    err = json.Unmarshal(resp.Body(), &weather)
    if err != nil {
        return &weather, err
    }

    return &weather, nil
}
