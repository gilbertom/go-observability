package webserver

import (
	"encoding/json"
	"net/http"

	"github.com/gilbertom/go-temperatura-cep/internal/usecase"
	"github.com/gilbertom/go-temperatura-cep/internal/web/webserver/dto"
	"go.opentelemetry.io/otel/trace"
)

// HTTPHandler handles HTTP requests.
type HTTPHandler struct {
    Tracer trace.Tracer
    cepUsecase *usecase.CepUsecase
    weatherUsecase *usecase.WeatherUsecase
}

// NewHTTPHandler creates a new HTTPHandler instance.
func NewHTTPHandler(tracer trace.Tracer, u *usecase.CepUsecase, w *usecase.WeatherUsecase) *HTTPHandler {
    return &HTTPHandler{
        Tracer: tracer,
        cepUsecase: u,
        weatherUsecase: w,
    }
}

// GetTemperaturesByCep handles the request to get the temperatures by CEP.
func (h *HTTPHandler) GetTemperaturesByCep(w http.ResponseWriter, r *http.Request) {
    ctx, span := h.Tracer.Start(r.Context(), "GetTemperaturesByCep")
    defer span.End()

    cep := r.URL.Query().Get("cep")
    
    if validCep := h.cepUsecase.ValidateCep(cep); !validCep {
        http.Error(w, "CEP is invalid", http.StatusBadRequest)
        return
    }

    locality, err := h.cepUsecase.GetLocalityByCep(ctx, cep, h.Tracer)
    if err != nil {
        if err.Error() == "invalid zipcode" {
            http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
            return
        }

        if err.Error() == "can not find zipcode" {
            http.Error(w, "can not find zipcode", http.StatusNotFound)
            return
        }
    }

    weather, err := h.weatherUsecase.GetTemperaturesByLocality(ctx, locality.Localidade, h.Tracer)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    response := dto.WeatherResponse{
        City:       locality.Localidade,
        TempC:      weather.Current.TempC,
        TempF:      h.weatherUsecase.ConvertCelsiusToFahrenheit(weather.Current.TempC),
        TempK:      h.weatherUsecase.ConvertCelsiusToKelvin(weather.Current.TempC),
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}
