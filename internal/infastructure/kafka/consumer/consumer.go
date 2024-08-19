package consumer

import (
	"context"
	"cors/internal/config"
	userentity "cors/internal/entity/user"
	userusecase "cors/internal/usecase/user"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type Consumer struct {
	consumer        *kafka.Consumer
	user            *userusecase.Repo
	topicCreateUser string
	topicVeryFy     string
}

func NewConsumer(cfg *config.Config, user *userusecase.Repo) (*Consumer, error) {

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.KafkaUrl,
	})
	if err != nil {
		return nil, err
	}
	return &Consumer{
		consumer:        consumer,
		user:            user,
		topicCreateUser: cfg.CreateUserTopic,
		topicVeryFy:     cfg.VeryFyTopic,
	}, nil
}

func (c *Consumer) Consume() {
	forever := make(chan bool)
	go func() {
		err := c.consumer.SubscribeTopics([]string{c.topicCreateUser}, nil)
		if err != nil {
			panic(err)
		}

		for {
			msg, err := c.consumer.ReadMessage(-1)
			if err != nil {
				log.Fatal(err)
			}
			var user *userentity.User

			err = json.Unmarshal(msg.Value, &user)
			if err != nil {
				log.Fatal(err)
			}
			err = c.user.SaveUserToRedis(context.Background(), user)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	go func() {
		err := c.consumer.SubscribeTopics([]string{c.topicVeryFy}, nil)
		if err != nil {
			panic(err)
		}

		for {
			msg, err := c.consumer.ReadMessage(-1)
			if err != nil {
				log.Fatal(err)
			}
			var user *userentity.User

			err = json.Unmarshal(msg.Value, &user)
			if err != nil {
				log.Fatal(err)
			}
			err = c.user.SaveUserToMongoDB(context.Background(), user)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
	<-forever
}
