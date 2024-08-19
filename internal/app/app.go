package app

import (
	httpapp "cors/internal/app/http"
	"cors/internal/config"
	userservice "cors/internal/services/user"
	userusecase "cors/internal/usecase/user"
	"log/slog"
)

type App struct {
	HTTPApp *httpapp.HttpApp
}

func NewApp(logger *slog.Logger, cfg *config.Config) *App {
	userService := userservice.NewUser()
	user := userusecase.NewUserUseCase(userService)
	server := httpapp.NewApp(logger, cfg.HTTPUrl, user)
	return &App{
		HTTPApp: server,
	}
}
