package httpapp

import (
	"cors/internal/http/router"
	userusecase "cors/internal/usecase/user"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type HttpApp struct {
	Logger  *slog.Logger
	HTTPUrl string
	Server  *gin.Engine
}

func NewApp(logger *slog.Logger, httpUrl string, user *userusecase.User) *HttpApp {
	server := router.NewRouter(user)
	return &HttpApp{
		Logger:  logger,
		HTTPUrl: httpUrl,
		Server:  server,
	}
}
