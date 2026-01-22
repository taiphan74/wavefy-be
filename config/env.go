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
		AppEnv: getenvRequired("APP_ENV"),
		DB: DBConfig{
			Host:            getenvRequired("DB_HOST"),
			Port:            getenvRequired("DB_PORT"),
			User:            getenvRequired("DB_USER"),
			Password:        getenvRequired("DB_PASSWORD"),
			Name:            getenvRequired("DB_NAME"),
			SSLMode:         getenvRequired("DB_SSLMODE"),
			MaxOpenConns:    getenvInt("DB_MAX_OPEN_CONNS", 20),
			MaxIdleConns:    getenvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getenvDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute),
			ConnMaxIdleTime: getenvDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
		},
		Auth: AuthConfig{
			JWTSecret:      getenvRequired("AUTH_JWT_SECRET"),
			AccessTokenTTL: getenvDuration("AUTH_ACCESS_TOKEN_TTL", 24*time.Hour),
			AccessTokenIss: getenvRequired("AUTH_ACCESS_TOKEN_ISSUER"),
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

func getenvRequired(key string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	panic("missing required env: " + key)
}
