package config

import "time"

type Config struct {
	Port   string
	AppEnv string
	DB     DBConfig
	Auth   AuthConfig
	Redis  RedisConfig
	Mail   MailConfig
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
	JWTSecret           string
	AccessTokenTTL      time.Duration
	AccessTokenIss      string
	RefreshTokenTTL     time.Duration
	RefreshTokenSecret  string
	PasswordResetTTL    time.Duration
	PasswordResetSecret string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type MailConfig struct {
	Host string
	Port int
	User string
	Pass string
	From string
}
