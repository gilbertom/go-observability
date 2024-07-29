package config

import (
	"log"
	"os"

	"github.com/subosito/gotenv"
)

// Config is the configuration struct for the application.
type Config struct {
	PortHTTPServiceA      string
	PortHTTPServiceB      string
	URLServiceB						string
	URLCep                string
	URLWeather            string
	APIKeyWeather         string
}

// AppConfig is the configuration for the application.
var AppConfig Config

// LoadConfig loads the configuration for the application.
func LoadConfig() {
	err := gotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AppConfig = Config{
		PortHTTPServiceA:       os.Getenv("PORT_HTTP_SERVICE_A"),
		PortHTTPServiceB:       os.Getenv("PORT_HTTP_SERVICE_B"),
		URLServiceB:            os.Getenv("URL_SERVICE_B"),
		URLCep:                 os.Getenv("URL_CEP"),
		URLWeather:             os.Getenv("URL_WEATHER"),
		APIKeyWeather:          os.Getenv("API_KEY_WEATHER"),
	}
}
