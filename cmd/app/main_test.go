package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTemperaturesByCep(t *testing.T) {
    go func() {
        main()
    }()
    
    time.Sleep(1 * time.Second)
    
    t.Run("Valid CEP", func(t *testing.T) {
        resp, err := http.Get("http://localhost:8080/cep?cep=28951620")
        if err != nil {
            t.Fatalf("Falha ao fazer a requisição: %v", err)
        }
        defer resp.Body.Close()
        
        assert.Equal(t, http.StatusOK, resp.StatusCode, "Esperado status 200")

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            t.Fatalf("Falha ao ler o corpo da resposta: %v", err)
        }
        
        var response map[string]float64
        if err := json.Unmarshal(body, &response); err != nil {
            t.Fatalf("Falha ao deserializar a resposta JSON: %v, %v", err, string(body))
        }
        fmt.Println("Response Status:", resp.Status)
        fmt.Println("Response Body:", string(body))

        keys := []string{"temp_C", "temp_F", "temp_K"}
        for _, key := range keys {
            if _, ok := response[key]; !ok {
                t.Fatalf("Chave %q não encontrada na resposta", key)
            }
        }
    })

    t.Run("Can not find zipcode", func(t *testing.T) {
        resp, err := http.Get("http://localhost:8080/cep?cep=88888888")
        if err != nil {
            t.Fatalf("Falha ao fazer a requisição: %v", err)
        }
        defer resp.Body.Close()
        
        assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Esperado status 404")

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            t.Fatalf("Falha ao ler o corpo da resposta: %v", err)
        }
        
        expectedErrorMessage := "can not find zipcode"
        assert.Contains(t, string(body), expectedErrorMessage, "Esperado mensagem de erro")
        
        fmt.Println("Response Status:", resp.Status)
        fmt.Println("Response Body:", string(body))
    })

    t.Run("Invalid zipcode format", func(t *testing.T) {
        resp, err := http.Get("http://localhost:8080/cep?cep=2895162A")
        if err != nil {
            t.Fatalf("Falha ao fazer a requisição: %v", err)
        }
        defer resp.Body.Close()
        
        assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode, "Esperado status 422")

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            t.Fatalf("Falha ao ler o corpo da resposta: %v", err)
        }
        
        expectedErrorMessage := "invalid zipcode"
        assert.Contains(t, string(body), expectedErrorMessage, "Esperado mensagem de erro")
        
        fmt.Println("Response Status:", resp.Status)
        fmt.Println("Response Body:", string(body))
    })
}