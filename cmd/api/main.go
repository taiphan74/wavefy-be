package main

import (
	"wavefy-be/config"
	"wavefy-be/internal/app"
)

func main() {
	cfg := config.Load()

	if err := app.RunHTTP(":" + cfg.Port); err != nil {
		panic(err)
	}
}
