package config

import (
	"fmt"
	"os"
	"strconv"
)

// gateway configuration from env vars
type Config struct {
	ServerPort        string
	DatabaseURL       string
	PostgresHost      string
	PostgresPort      string
	PostgresUser      string
	PostgresPassword  string
	PostgresDB        string
	RabbitMQURL       string
	TargetURL         string
	WorkerConcurrency int
	AdminAPIKey       string
	IngestAPIKey      string
	CorsAllowedOrigin string
}

// reads env vars with defaults
func Load() *Config {
	// render sets PORT, docker-compose uses SERVER_PORT
	port := getEnv("PORT", getEnv("SERVER_PORT", "8080"))

	return &Config{
		ServerPort:        port,
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		PostgresHost:      getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:      getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:      getEnv("POSTGRES_USER", "siphon"),
		PostgresPassword:  getEnv("POSTGRES_PASSWORD", "siphon_secret"),
		PostgresDB:        getEnv("POSTGRES_DB", "siphon_gateway"),
		RabbitMQURL:       getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		TargetURL:         getEnv("TARGET_URL", fmt.Sprintf("http://localhost:%s/api/v1/echo", port)),
		WorkerConcurrency: getEnvInt("WORKER_CONCURRENCY", 5),
		AdminAPIKey:       getEnv("ADMIN_API_KEY", ""),
		IngestAPIKey:      getEnv("INGEST_API_KEY", ""),
		CorsAllowedOrigin: getEnv("CORS_ALLOWED_ORIGIN", "http://localhost:5173"),
	}
}

// returns the postgres connection string
func (c *Config) PostgresDSN() string {
	// Si pasamos la URL completa desde Render, usamos esa directamente
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}

	// Si no (por ejemplo en local con Docker), construimos la URL como antes
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.PostgresUser, c.PostgresPassword, c.PostgresHost, c.PostgresPort, c.PostgresDB,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
