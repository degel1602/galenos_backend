// Package config centraliza la lectura de configuración desde variables de
// entorno. Es la única capa, junto a main, que conoce de dónde vienen los
// valores externos.
package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config agrupa toda la configuración externa de la aplicación.
type Config struct {
	ServerPort        string
	ServerHost        string
	SQLServerDSN      string
	AllowedOrigins    []string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
	ReniecApp         string
	ReniecUsuario     string
	ReniecClave       string
	ReniecURL         string
	ReniecTimeout     time.Duration
	SISUsuario        string
	SISClave          string
	SISURL            string
	SISDNIAutorizado  string
	SISTimeout        time.Duration
	AuthUsername      string
	AuthPassword      string
	AuthSecret        string
	AuthTTL           time.Duration
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
		ServerHost:        envOrDefault("SERVER_HOST", defaultServerHost()),
		SQLServerDSN:      dsn,
		AllowedOrigins:    allowedOriginsWithServer(envOrDefault("ALLOWED_ORIGINS", "http://localhost:4200"), envOrDefault("SERVER_PORT", "8080"), envOrDefault("SERVER_HOST", defaultServerHost())),
		DBMaxOpenConns:    envIntOrDefault("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:    envIntOrDefault("DB_MAX_IDLE_CONNS", 10),
		DBConnMaxLifetime: envDurationOrDefault("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		ReniecApp:         envOrDefault("RENIEC_APP", "HNSEB"),
		ReniecUsuario:     envOrDefault("RENIEC_USUARIO", "44602631"),
		ReniecClave:       os.Getenv("RENIEC_CLAVE"),
		ReniecURL:         envOrDefault("RENIEC_URL", "https://wsvmin.minsa.gob.pe/wsreniecmq/serviciomq.asmx"),
		ReniecTimeout:     envDurationOrDefault("RENIEC_TIMEOUT", 30*time.Second),
		SISUsuario:        envOrDefault("SIS_USUARIO", "HNAL"),
		SISClave:          os.Getenv("SIS_CLAVE"),
		SISURL:            envOrDefault("SIS_URL", "http://app.sis.gob.pe/sisWSAFI/Service.asmx"),
		SISDNIAutorizado:  envOrDefault("SIS_DNI_AUTORIZADO", ""),
		SISTimeout:        envDurationOrDefault("SIS_TIMEOUT", 30*time.Second),
		AuthUsername:      apiUsername,
		AuthPassword:      apiPassword,
		AuthSecret:        jwtSecret,
		AuthTTL:           envDurationOrDefault("JWT_TTL", 24*time.Hour),
	}, nil
}

// defaultServerHost devuelve el IP de la interfaz de red principal para
// que otros usuarios de la red local puedan consumir la API (en lugar de
// localhost). Si no se puede detectar, retorna localhost.
func defaultServerHost() string {
	host := localIPv4()
	if host == "" {
		return "localhost"
	}
	return host
}

// localIPv4 resuelve el IP IPv4 saliente consultando la interfaz por
// defecto (conexión UDP sin tráfico real). Retorna vacío si falla.
func localIPv4() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	if addr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
		return addr.IP.String()
	}
	return ""
}

// allowedOriginsWithServer agrega el origen propio de la API (IP
// configurado y localhost) a la lista de orígenes permitidos, de modo que
// Swagger y clientes que apuntan a la propia máquina no fallen por CORS.
func allowedOriginsWithServer(csv, port, host string) []string {
	origins := strings.Split(csv, ",")
	for _, o := range origins {
		o = strings.TrimSpace(o)
		if o == "" {
			continue
		}
		if o == "http://"+host+":"+port || o == "http://localhost:"+port {
			return origins
		}
	}

	serverOrigin := "http://" + host + ":" + port
	if !strings.Contains(csv, serverOrigin) {
		origins = append(origins, serverOrigin)
	}
	if host != "localhost" && !strings.Contains(csv, "http://localhost:"+port) {
		origins = append(origins, "http://localhost:"+port)
	}
	return origins
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
