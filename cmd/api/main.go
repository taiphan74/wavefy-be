// @title           Wavefy API
// @version         1.0
// @description     API documentation for Wavefy BE
// @host            localhost:8080
// @BasePath        /api
package main

import (
	"context"

	"wavefy-be/config"
	"wavefy-be/docs"
	"wavefy-be/internal/app"
	"wavefy-be/internal/cache"
	"wavefy-be/internal/db"
	"wavefy-be/internal/mail"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	conn, err := db.Open(ctx, cfg.DB)
	if err != nil {
		panic(err)
	}
	defer func() {
		sqlDB, dbErr := conn.DB()
		if dbErr == nil {
			_ = sqlDB.Close()
		}
	}()

	redisClient, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = redisClient.Close()
	}()

	docs.SwaggerInfo.Host = "localhost:" + cfg.Port
	docs.SwaggerInfo.BasePath = "/api"

	mailer, err := mail.FromConfig(cfg.Mail)
	if err != nil {
		panic(err)
	}

	server := app.NewHTTP(conn, redisClient, cfg.Auth, mailer)
	if err := server.Run(":" + cfg.Port); err != nil {
		panic(err)
	}
}
