package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumerFromConfig(config ConsumerConfig) *Consumer {
	return &Consumer{reader: kafka.NewReader(config.ToKafkaReaderConfig())}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func (c *Consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}

func (c *Consumer) FetchMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.FetchMessage(ctx)
}

func (c *Consumer) CommitMessage(ctx context.Context, message kafka.Message) error {
	return c.reader.CommitMessages(ctx, message)
}
