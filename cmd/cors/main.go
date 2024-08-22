package main

import (
	"cors/internal/app"
	"cors/internal/config"
	"cors/logger"
)

func main() {
	cfg := config.New()
	log := logger.SetupLogger(cfg.LogLevel)
	application := app.NewApp(log, cfg)
	application.HTTPApp.Start()
}
