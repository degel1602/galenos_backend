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
	JWTSecret         string
	JWTExpiration     time.Duration
	APIUsername       string
	APIPassword       string
}

// Load lee la configuración desde el entorno aplicando valores por defecto
// razonables para desarrollo.
func Load() (*Config, error) {
	dsn := os.Getenv("SQLSERVER_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("SQLSERVER_DSN environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	apiUsername := os.Getenv("API_USERNAME")
	if apiUsername == "" {
		return nil, fmt.Errorf("API_USERNAME environment variable is required")
	}

	apiPassword := os.Getenv("API_PASSWORD")
	if apiPassword == "" {
		return nil, fmt.Errorf("API_PASSWORD environment variable is required")
	}

	return &Config{
		ServerPort:        envOrDefault("SERVER_PORT", "8080"),
		SQLServerDSN:      dsn,
		AllowedOrigins:    strings.Split(envOrDefault("ALLOWED_ORIGINS", "http://localhost:4200"), ","),
		DBMaxOpenConns:    envIntOrDefault("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:    envIntOrDefault("DB_MAX_IDLE_CONNS", 10),
		DBConnMaxLifetime: envDurationOrDefault("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		JWTSecret:         jwtSecret,
		JWTExpiration:     envDurationOrDefault("JWT_EXPIRATION", 24*time.Hour),
		APIUsername:       apiUsername,
		APIPassword:       apiPassword,
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
