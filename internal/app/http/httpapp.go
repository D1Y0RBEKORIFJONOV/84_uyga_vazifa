package httpapp

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"cors/internal/http/router"
	userusecase "cors/internal/usecase/user"
	"github.com/gin-gonic/gin"
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

func (app *HttpApp) Start() {
	const op = "app.Start"
	log := app.Logger.With(
		slog.String(op, "Starting server"),
		slog.String("port", app.HTTPUrl))
	log.Info("Starting server ")
	err := app.Server.Run(app.HTTPUrl)
	if err != nil {
		log.Info("Failed to start server", "error", err)
	}
}

func (h *HttpApp) Shutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	const op = "HttpApp.Shutdown"
	log := h.Logger.With(
		slog.String(op, "shutting down http app"),
		slog.String("port", h.HTTPUrl),
	)
	log.Info("shutting down http app")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error(op, err)
	}

	log.Info("http app stopped gracefully")
}
