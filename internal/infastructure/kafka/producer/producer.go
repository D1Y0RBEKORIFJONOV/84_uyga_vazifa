package producer

import (
	"cors/internal/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(cfg *config.Config) (*Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.KafkaUrl,
	})
	if err != nil {
		return nil, err
	}
	return &Producer{
		producer: producer,
	}, nil
}

func (p *Producer) Close() {
	p.producer.Close()
}

func (p *Producer) Publish(value []byte, topicKey string) error {
	if err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicKey, Partition: kafka.PartitionAny},
		Value:          value,
	}, nil); err != nil {
		return err
	}
	return nil
}
