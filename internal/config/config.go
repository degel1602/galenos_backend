// Package config centraliza la lectura de configuración desde variables de
// entorno. Es la única capa, junto a main, que conoce de dónde vienen los
// valores externos.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config agrupa toda la configuración externa de la aplicación.
type Config struct {
	ServerPort        string
	SQLServerDSN      string
	AllowedOrigins    []string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
}

// Load lee la configuración desde el entorno aplicando valores por defecto
// razonables para desarrollo.
func Load() (*Config, error) {
	dsn := os.Getenv("SQLSERVER_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("SQLSERVER_DSN environment variable is required")
	}

	return &Config{
		ServerPort:        envOrDefault("SERVER_PORT", "8080"),
		SQLServerDSN:      dsn,
		AllowedOrigins:    strings.Split(envOrDefault("ALLOWED_ORIGINS", "http://localhost:4200"), ","),
		DBMaxOpenConns:    envIntOrDefault("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:    envIntOrDefault("DB_MAX_IDLE_CONNS", 10),
		DBConnMaxLifetime: envDurationOrDefault("DB_CONN_MAX_LIFETIME", 5*time.Minute),
	}, nil
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envIntOrDefault(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return parsed
}

func envDurationOrDefault(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return parsed
}
