package consumer

import (
	"context"
	"cors/internal/config"
	userentity "cors/internal/entity/user"
	userusecase "cors/internal/usecase/user"
	"encoding/json"
	"fmt"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
)

type Consumer struct {
	consumer        *kgo.Client
	user            *userusecase.Repo
	topicCreateUser string
	topicVeryFy     string
}

func NewConsumer(cfg *config.Config, user *userusecase.Repo) (*Consumer, error) {
	consumer, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.KafkaUrl),
		kgo.ConsumeTopics(cfg.CreateUserTopic, cfg.VeryFyTopic),
		kgo.ConsumerGroup("1"),
	)

	err = createTopics(cfg)
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
	log.Println("Consumer start")
	ctx := context.Background()
	for {
		fetches := c.consumer.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				log.Println("Error fetching records:", err)
			}
			continue
		}
		fetches.EachPartition(func(ftp kgo.FetchTopicPartition) {
			for _, record := range ftp.Records {
				fmt.Println("Partition:", record.Partition, "Value:", string(record.Value))
				var user userentity.User
				err := json.Unmarshal(record.Value, &user)
				if err != nil {
					log.Println("Error unmarshaling record:", err)
					continue
				}

				switch record.Topic {
				case c.topicCreateUser:
					err = c.user.SaveUserToRedis(ctx, &user)
				case c.topicVeryFy:
					err = c.user.SaveUserToMongo(ctx, &user)
				default:
					log.Println("Unknown topic:", record.Topic)
					continue
				}

				if err != nil {
					log.Println("Error saving user:", err)
				}
			}
		})
	}
}

func createTopics(cfg *config.Config) error {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.KafkaUrl),
	)
	if err != nil {
		return err
	}
	defer client.Close()

	admin := kadm.NewClient(client)
	ctx := context.Background()

	topics := []string{cfg.CreateUserTopic, cfg.VeryFyTopic}

	createResp, err := admin.CreateTopics(ctx, 1, 1, nil, topics...)
	if err != nil {
		return err
	}

	for _, ctr := range createResp.Sorted() {
		if ctr.Err != nil {
			log.Printf("Error creating topic %s: %v", ctr.Topic, ctr.Err)
		} else {
			log.Printf("Topic %s created successfully", ctr.Topic)
		}
	}

	return nil
}
