package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gilbertom/go-temperatura-cep/internal/config"
	"github.com/gilbertom/go-temperatura-cep/internal/entity"
	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

// CepRepository represents a repository for handling CEP data.
type CepRepository struct{}

// NewCepRepository creates a new instance of CepRepository.
func NewCepRepository() *CepRepository {
    return &CepRepository{}
}

// GetLocalityByCep retrieves the locality information for a given CEP.
func (r *CepRepository) GetLocalityByCep(ctx context.Context, cep string, tracer trace.Tracer) (*entity.Cep, error) {
    var locality entity.Cep
    url := fmt.Sprintf("%s/%s/json/", config.AppConfig.URLCep, cep)

    client := resty.New()
    client.SetTransport(otelhttp.NewTransport(http.DefaultTransport))
    resp, err := client.R().
        SetContext(ctx).
        SetHeader("Content-Type", "application/json").
        Get(url)
    if err != nil {
        return &locality, err
    }
    defer resp.RawBody().Close()

    err = json.Unmarshal(resp.Body(), &locality)
    if err != nil {
        return &locality, err
    }

    if resp.StatusCode() == http.StatusNotFound {
        return &locality, errors.New("can not find zipcode")
    }

    if locality.Erro == "true" {
        return &locality, errors.New("can not find zipcode")
    }

    if locality.Localidade == "" {
        return &locality, errors.New("invalid zipcode")
    }

    if resp.StatusCode() != http.StatusOK {
        return &locality, fmt.Errorf("failed to call URL: %s, status code: %d", url, resp.StatusCode())
    }

    return &locality, nil
}
