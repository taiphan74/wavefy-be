package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Load reads configuration from .env and environment variables.
func Load() Config {
	_ = godotenv.Load(".env")

	return Config{
		Port:   getenv("PORT", "8080"),
		AppEnv: getenv("APP_ENV", "development"),
	}
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}

	return fallback
}
