package producer

import (
	"context"
	"cors/internal/config"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
	"time"
)

type Producer struct {
	producer *kgo.Client
}

func NewProducer(cfg *config.Config) (*Producer, error) {
	var (
		err      error
		producer *kgo.Client
	)
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Millisecond)
		producer, err = kgo.NewClient(
			kgo.SeedBrokers(cfg.KafkaUrl),
			kgo.AllowAutoTopicCreation(),
		)
		if err != nil {
			log.Printf("Error creating Kafka producer, attempt %d: %v\n", i+1, err)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
	}, nil
}

func (p *Producer) Close() {
	log.Println("Closing producer")
	p.producer.Close()
}

func (p *Producer) Publish(value []byte, topicKey string) error {
	log.Println("Starting Publish")
	defer log.Println("End Publish")

	record := &kgo.Record{
		Key:   []byte(topicKey),
		Topic: topicKey,
		Value: value,
	}
	log.Println("Publishing record:", record)

	// Используем канал для ожидания завершения отправки
	done := make(chan struct{})
	p.producer.Produce(context.Background(), record, func(record *kgo.Record, err error) {
		if err != nil {
			log.Println("Failed to publish record:", err)
		} else {
			log.Println("Record published successfully")
		}
		close(done) // Закрываем канал после завершения отправки
	})

	<-done // Ждем завершения отправки сообщения

	return nil
}
