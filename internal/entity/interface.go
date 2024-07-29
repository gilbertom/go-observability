package entity

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// CepRepository represents a repository for CEP data.
type CepRepository interface {
    GetLocalityByCep(ctx context.Context, cep string, tracer trace.Tracer) (*Cep, error)
}

// WeatherRepository represents a repository for weather data.
type WeatherRepository interface {
    GetTemperaturesByLocality(ctx context.Context, locality string, trace trace.Tracer) (*Weather, error)
}