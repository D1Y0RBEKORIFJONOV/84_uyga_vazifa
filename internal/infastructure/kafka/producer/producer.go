package producer

import (
	"context"
	"cors/internal/config"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	producer *kgo.Client
}

func NewProducer(cfg *config.Config) (*Producer, error) {
	producer, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.KafkaUrl),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	err = producer.Ping(ctx)
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

	record := &kgo.Record{
		Key:   []byte(topicKey),
		Topic: topicKey,
		Value: value,
	}
	err := p.producer.ProduceSync(context.Background(), record).FirstErr()
	if err != nil {
		return err
	}
	return nil
}
