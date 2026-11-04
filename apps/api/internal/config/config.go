package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds runtime configuration for the API.
type Config struct {
	Port        int
	DatabaseURL string
	Env         string
}

// Load reads configuration from environment variables, applying sensible defaults.
func Load() (Config, error) {
	cfg := Config{
		Port:        8080,
		Env:         getenv("APP_ENV", "development"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
	if v := os.Getenv("API_PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid API_PORT: %w", err)
		}
		cfg.Port = p
	}
	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
