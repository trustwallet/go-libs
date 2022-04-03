package kafka

import (
	"context"
	"fmt"

	kafka "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducerFromConfig(config ProducerConfig) *Producer {
	return &Producer{writer: kafka.NewWriter(config.ToKafkaWriterConfig())}
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func (p *Producer) SendMessage(ctx context.Context, key, value []byte, topic string) (int, int64, error) {
	message := kafka.Message{
		Value: value,
		Topic: topic,
	}

	err := p.writer.WriteMessages(ctx, message)
	if err != nil {
		return message.Partition, message.Offset, fmt.Errorf("failed to send message to kafka topic (%s): %w", topic, err)
	}

	return message.Partition, message.Offset, nil
}
