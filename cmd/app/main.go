package main

import (
	"cors/internal/app"
	"cors/internal/config"
	"cors/logger"
)

func main() {
	cfg := config.New()
	log := logger.SetupLogger("local")
	application := app.NewApp(log, cfg)
	err := application.HTTPApp.Server.Run(cfg.HTTPUrl)
	if err != nil {
		panic(err)
	}
}
