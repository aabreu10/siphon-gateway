package config

import (
	"fmt"
	"os"
	"strconv"
)

// gateway configuration from env vars
type Config struct {
	ServerPort        string
	PostgresHost      string
	PostgresPort      string
	PostgresUser      string
	PostgresPassword  string
	PostgresDB        string
	RabbitMQURL       string
	TargetURL         string
	WorkerConcurrency int
}

// reads env vars with defaults
func Load() *Config {
	return &Config{
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		PostgresHost:      getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:      getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:      getEnv("POSTGRES_USER", "siphon"),
		PostgresPassword:  getEnv("POSTGRES_PASSWORD", "siphon_secret"),
		PostgresDB:        getEnv("POSTGRES_DB", "siphon_gateway"),
		RabbitMQURL:       getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		TargetURL:         getEnv("TARGET_URL", "http://localhost:8080/api/v1/echo"),
		WorkerConcurrency: getEnvInt("WORKER_CONCURRENCY", 5),
	}
}

// returns the postgres connection string
func (c *Config) PostgresDSN() string {
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
