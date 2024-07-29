package repository

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gilbertom/go-temperatura-cep/internal/config"
)

func mockHTTPServer(response string, statusCode int) *httptest.Server {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(statusCode)
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(response))
    })
    return httptest.NewServer(handler)
}

func TestGetLocalityByCep(t *testing.T) {
    originalURL := config.AppConfig.URLCep
    defer func() { config.AppConfig.URLCep = originalURL }()
    
    tests := []struct {
        cep         string
        response    string
        statusCode  int
        expectedErr bool
        expectedLoc string
    }{
        {"01001000", `{"cep": "01001000", "localidade": "São Paulo"}`, http.StatusOK, false, "São Paulo"},
        {"28951620", `{"cep": "28951620", "localidade": "Cabo Frio"}`, http.StatusOK, false, "Cabo Frio"},
        {"88888888", `{"erro": "true"}`, http.StatusOK, true, ""},
        {"12345678", `{"erro": "true"}`, http.StatusOK, true, ""},
        {"1234567A", "", http.StatusBadRequest, true, ""},
        {"1234567", "", http.StatusBadRequest, true, ""},
        {"123456789", "", http.StatusBadRequest, true, ""},
    }

    for _, tt := range tests {
        t.Run(tt.cep, func(t *testing.T) {
            server := mockHTTPServer(tt.response, tt.statusCode)
            defer server.Close()
            
            config.AppConfig.URLCep = server.URL
            
            repo := NewCepRepository()
            locality, err := repo.GetLocalityByCep(nil, tt.cep, nil)
            if (err != nil) != tt.expectedErr {
                t.Fatalf("expected error: %v, got: %v", tt.expectedErr, err)
            }
            if locality.Localidade != tt.expectedLoc {
                t.Fatalf("expected locality: %v, got: %v", tt.expectedLoc, locality.Localidade)
            }
        })
    }
}
