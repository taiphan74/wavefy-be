package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Load reads configuration from .env and environment variables.
func Load() Config {
	_ = godotenv.Load(".env")

	return Config{
		Port:   getenv("PORT", "8080"),
		AppEnv: getenv("APP_ENV", "development"),
		DB: DBConfig{
			Host:            getenv("DB_HOST", "localhost"),
			Port:            getenv("DB_PORT", "5432"),
			User:            getenv("DB_USER", "wavefy"),
			Password:        getenv("DB_PASSWORD", "wavefy"),
			Name:            getenv("DB_NAME", "wavefy"),
			SSLMode:         getenv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getenvInt("DB_MAX_OPEN_CONNS", 20),
			MaxIdleConns:    getenvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getenvDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute),
			ConnMaxIdleTime: getenvDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
		},
	}
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}

	return fallback
}

func getenvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}

	return fallback
}

func getenvDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}

	return fallback
}
