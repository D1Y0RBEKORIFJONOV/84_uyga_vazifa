package consumer

import (
	"context"
	"cors/internal/config"
	userentity "cors/internal/entity/user"
	userusecase "cors/internal/usecase/user"
	"encoding/json"
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
	"time"
)

type Consumer struct {
	consumer        *kgo.Client
	user            *userusecase.Repo
	topicCreateUser string
	topicVeryFy     string
}

func NewConsumer(cfg *config.Config, user *userusecase.Repo) (*Consumer, error) {
	var (
		err      error
		consumer *kgo.Client
	)

	// Попытка подключения к Kafka до 10 раз
	for i := 0; i < 10; i++ {
		consumer, err = kgo.NewClient(
			kgo.SeedBrokers(cfg.KafkaUrl),
			kgo.ConsumeTopics(cfg.CreateUserTopic, cfg.VeryFyTopic),
			kgo.ConsumerGroup("111"), // Убедитесь, что название группы уникальное
		)
		if err != nil {
			log.Printf("Error creating Kafka consumer, attempt %d: %v\n", i+1, err)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %v", err)
	}

	return &Consumer{
		consumer:        consumer,
		user:            user,
		topicCreateUser: cfg.CreateUserTopic,
		topicVeryFy:     cfg.VeryFyTopic,
	}, nil
}

func (c *Consumer) Consume() {
	log.Println("Consumer start")
	ctx := context.Background()

	for {
		// Получаем сообщения из Kafka
		fetches := c.consumer.PollFetches(ctx)

		// Проверка на ошибки
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				log.Println("Error fetching records:", err)
			}
			continue
		}

		// Обрабатываем каждую партицию
		fetches.EachPartition(func(ftp kgo.FetchTopicPartition) {
			log.Printf("Reading from topic: %s, partition: %d\n", ftp.Topic, ftp.Partition)

			for _, record := range ftp.Records {
				log.Printf("Partition: %d, Offset: %d, Value: %s\n", record.Partition, record.Offset, string(record.Value))

				var user userentity.User
				err := json.Unmarshal(record.Value, &user)
				if err != nil {
					log.Println("Error unmarshaling record:", err)
					continue
				}

				// Определяем, из какого топика пришло сообщение и сохраняем пользователя
				switch record.Topic {
				case c.topicCreateUser:
					log.Println("Processing CreateUser message")
					err = c.user.SaveUserToRedis(ctx, &user)
				case c.topicVeryFy:
					log.Println("Processing VeryFy message")
					err = c.user.SaveUserToMongo(ctx, &user)
				default:
					log.Println("Unknown topic:", record.Topic)
					continue
				}

				// Логирование ошибок при сохранении пользователя
				if err != nil {
					log.Println("Error saving user:", err)
				}
			}
		})
	}
}
