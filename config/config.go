package config

import "time"

type Config struct {
	Port   string
	AppEnv string
	DB     DBConfig
	Auth   AuthConfig
}

type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type AuthConfig struct {
	JWTSecret       string
	AccessTokenTTL  time.Duration
	AccessTokenIss  string
}
