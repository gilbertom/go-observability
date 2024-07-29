package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"github.com/gilbertom/go-temperatura-cep/internal/config"
	"github.com/gilbertom/go-temperatura-cep/internal/entity"
	"github.com/gilbertom/go-temperatura-cep/internal/repository"
	"github.com/gilbertom/go-temperatura-cep/internal/usecase"
	"github.com/gilbertom/go-temperatura-cep/internal/web/webserver"
	"github.com/go-resty/resty/v2"
)

var tracer trace.Tracer

func main() {
    config.LoadConfig()

    tp, err := initTracer()
    if err != nil {
        log.Fatalf("failed to initialize tracer: %v", err)
    }
    defer func() { _ = tp.Shutdown(context.Background()) }()

    go func() {
        http.HandleFunc("/serviceA", handlePostServiceA)
        log.Println("Service A running at:", config.AppConfig.PortHTTPServiceA)
        log.Fatal(http.ListenAndServe(":"+config.AppConfig.PortHTTPServiceA, otelhttp.NewHandler(http.DefaultServeMux, "Process Service A")))
    }()
    
    cepRepo := repository.NewCepRepository()
    weatherRepo := repository.NewWeatherRepository()

    cepUsecase := usecase.NewCepUsecase(cepRepo)
    weatherUsecase := usecase.NewWeatherUsecase(weatherRepo)

    HTTPHandler := webserver.NewHTTPHandler(tracer, cepUsecase, weatherUsecase)

    log.Println("Service B running at:", config.AppConfig.PortHTTPServiceB)
    http.HandleFunc("/", HTTPHandler.GetTemperaturesByCep)
    log.Fatal(http.ListenAndServe(":"+config.AppConfig.PortHTTPServiceB, otelhttp.NewHandler(http.DefaultServeMux, "Process Service B")))
}

func handlePostServiceA(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    var cepRequest entity.CEPRequest

    err := json.NewDecoder(r.Body).Decode(&cepRequest)
    if err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    if !validateCEP(cepRequest.CEP) {
        http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
        return
    }

    respServiceB, err := callServiceB(ctx, cepRequest.CEP)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(respServiceB)
}

func callServiceB(ctx context.Context, cep string) (entity.ResponseServiceB, error) {
    ctx, span := tracer.Start(ctx, "callServiceB")
    defer span.End()
    
    var respServiceB entity.ResponseServiceB
    
    url := fmt.Sprintf("%s:%s?cep=%s", config.AppConfig.URLServiceB, config.AppConfig.PortHTTPServiceB, cep)
    client := resty.New()
    client.SetTransport(otelhttp.NewTransport(http.DefaultTransport))
    resp, err := client.R().
        SetContext(ctx).
        SetHeader("Content-Type", "application/json").
        Get(url)
    if err != nil {
        return respServiceB, err
    }
    defer resp.RawBody().Close()

    if resp.StatusCode() != http.StatusOK {
        return respServiceB, fmt.Errorf("failed to call URL: %s, status code: %d", url, resp.StatusCode())
    }

    err = json.Unmarshal(resp.Body(), &respServiceB)
    if err != nil {
        return respServiceB, err
    }
    return respServiceB, nil
}

func validateCEP(cep string) bool {
    match, _ := regexp.MatchString(`^\d{8}$`, cep)
    return match
}

func initTracer() (*sdktrace.TracerProvider, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    exporter, err := otlptracegrpc.New(ctx,
        otlptracegrpc.WithEndpoint("otel-collector:4317"),
        otlptracegrpc.WithInsecure())
    if err != nil {
        return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
    }

    tp := sdktrace.NewTracerProvider(
        sdktrace.WithSampler(sdktrace.AlwaysSample()),
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String("Tracing"),
        )),
    )

    otel.SetTracerProvider(tp)
    otel.SetTextMapPropagator(propagation.TraceContext{})

    tracer = otel.Tracer("tracing otel")

    return tp, nil
}
