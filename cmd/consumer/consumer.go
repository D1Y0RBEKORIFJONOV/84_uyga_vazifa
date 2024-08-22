package main

import (
	"cors/internal/config"
	"cors/internal/infastructure/kafka/consumer"
	"cors/internal/infastructure/kafka/producer"
	usermongodb "cors/internal/infastructure/repository/mongodb/user"
	userredis "cors/internal/infastructure/repository/redis/user"
	userusecase "cors/internal/usecase/user"
	"cors/logger"
)

func main() {
	cfg := config.New()
	log := logger.SetupLogger("local")
	mongoDb, err := usermongodb.NewMongoDB(cfg, log)
	if err != nil {
		panic(err)
	}

	broker, err := producer.NewProducer(cfg)
	if err != nil {
		panic(err)
	}
	redis := userredis.NewRedis(cfg)
	rep := userusecase.NewRepo(mongoDb, broker, mongoDb, mongoDb, redis)

	con, err := consumer.NewConsumer(cfg, rep)
	if err != nil {
		panic(err)
	}
	log.Info("Starting Consume")
	con.Consume()
}
