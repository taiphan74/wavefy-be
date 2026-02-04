package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
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
			JWTSecret:           getenvRequired("AUTH_JWT_SECRET"),
			AccessTokenTTL:      getenvDuration("AUTH_ACCESS_TOKEN_TTL", 24*time.Hour),
			AccessTokenIss:      getenvRequired("AUTH_ACCESS_TOKEN_ISSUER"),
			RefreshTokenTTL:     getenvDuration("AUTH_REFRESH_TOKEN_TTL", 7*24*time.Hour),
			RefreshTokenSecret:  getenvRequired("AUTH_REFRESH_TOKEN_SECRET"),
			PasswordResetTTL:    getenvDuration("AUTH_PASSWORD_RESET_TTL", 5*time.Minute),
			PasswordResetSecret: getenvRequired("AUTH_PASSWORD_RESET_SECRET"),
			VerifyEmailTTL:      getenvDuration("AUTH_VERIFY_EMAIL_TTL", 24*time.Hour),
			VerifyEmailSecret:   getenvRequired("AUTH_VERIFY_EMAIL_SECRET"),
		},
		Redis: RedisConfig{
			Addr:     getenvRequired("REDIS_ADDR"),
			Password: getenv("REDIS_PASSWORD", ""),
			DB:       getenvInt("REDIS_DB", 0),
		},
		Mail: MailConfig{
			Host: getenv("SMTP_HOST", ""),
			Port: getenvInt("SMTP_PORT", 0),
			User: getenv("SMTP_USER", ""),
			Pass: getenv("SMTP_PASS", ""),
			From: getenv("SMTP_FROM", ""),
		},
		Google: GoogleOAuthConfig{
			ClientID:      getenv("GOOGLE_CLIENT_ID", ""),
			ClientSecret:  getenv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:   getenv("GOOGLE_REDIRECT_URL", ""),
			OAuthStateTTL: getenvDuration("GOOGLE_OAUTH_STATE_TTL", 10*time.Minute),
		},
		R2: R2Config{
			AccountID:       getenvRequired("R2_ACCOUNT_ID"),
			Bucket:          getenvRequired("R2_BUCKET"),
			Region:          getenvRequired("R2_REGION"),
			AccessKeyID:     getenvRequired("R2_ACCESS_KEY_ID"),
			SecretAccessKey: getenvRequired("R2_SECRET_ACCESS_KEY"),
			Endpoint:        getenvRequired("R2_ENDPOINT"),
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
		if parsed, err := parseHumanDuration(value); err == nil {
			return parsed
		}
	}

	return fallback
}

func parseHumanDuration(value string) (time.Duration, error) {
	normalized := strings.TrimSpace(strings.ToLower(value))
	if normalized == "" {
		return 0, errors.New("invalid duration")
	}
	// Support days suffix (e.g. "7d") by converting to hours.
	if strings.HasSuffix(normalized, "d") {
		days := strings.TrimSuffix(normalized, "d")
		n, err := strconv.Atoi(days)
		if err != nil {
			return 0, err
		}
		return time.Duration(n) * 24 * time.Hour, nil
	}
	return time.ParseDuration(normalized)
}

func getenvRequired(key string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	panic("missing required env: " + key)
}
