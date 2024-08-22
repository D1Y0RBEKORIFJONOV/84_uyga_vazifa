package app

import (
	httpapp "cors/internal/app/http"
	"cors/internal/config"
	"cors/internal/infastructure/kafka/producer"
	usermongodb "cors/internal/infastructure/repository/mongodb/user"
	userredis "cors/internal/infastructure/repository/redis/user"
	userservice "cors/internal/services/user"
	userusecase "cors/internal/usecase/user"
	"log/slog"
)

type App struct {
	HTTPApp *httpapp.HttpApp
}

func NewApp(logger *slog.Logger, cfg *config.Config) *App {
	mongoDb, err := usermongodb.NewMongoDB(cfg, logger)
	if err != nil {
		panic(err)
	}

	broker, err := producer.NewProducer(cfg)
	if err != nil {
		panic(err)
	}
	redis := userredis.NewRedis(cfg)
	rep := userusecase.NewRepo(mongoDb, broker, mongoDb, mongoDb, redis)
	userService := userservice.NewUser(logger, rep, cfg)
	user := userusecase.NewUserUseCase(userService)
	server := httpapp.NewApp(logger, cfg.HTTPUrl, user)
	return &App{
		HTTPApp: server,
	}
}
