package main

import (
	"context"

	"wavefy-be/config"
	"wavefy-be/internal/app"
	"wavefy-be/internal/db"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	conn, err := db.Open(ctx, cfg.DB)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = conn.Close()
	}()

	server := app.NewHTTP(conn)
	if err := server.Run(":" + cfg.Port); err != nil {
		panic(err)
	}
}
