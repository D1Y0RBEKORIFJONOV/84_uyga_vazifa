package main

import (
	"cors/internal/config"
	"cors/internal/infastructure/kafka/consumer"
	usermongodb "cors/internal/infastructure/repository/mongodb/user"
	userredis "cors/internal/infastructure/repository/redis/user"
	userusecase "cors/internal/usecase/user"
	"cors/logger"
	"log"
)

type consume struct {
}

func (c *consume) Publish(value []byte, topicKey string) error {
	panic("implement me")
}

func main() {
	cfg := config.New()
	log1 := logger.SetupLogger("local")
	mongoDb, err := usermongodb.NewMongoDB(cfg, log1)
	if err != nil {
		log.Fatal(err, "AAAAAAAAAAAAAAAAAAAA")
	}

	redis := userredis.NewRedis(cfg)
	rep := userusecase.NewRepo(mongoDb, &consume{}, mongoDb, mongoDb, redis)

	con, err := consumer.NewConsumer(cfg, rep)
	if err != nil {
		panic(err)
	}
	log1.Info("Starting Consume", cfg.KafkaUrl)
	con.Consume()
}
