package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/trustwallet/go-libs/logging"
)

type RequiredAcks int

const (
	// NoResponse doesn't send any response, the TCP ACK is all you get.
	NoResponse RequiredAcks = 0
	// WaitForLocal waits for only the local commit to succeed before responding.
	WaitForLocal RequiredAcks = 1
	// WaitForAll waits for all in-sync replicas to commit before responding.
	WaitForAll RequiredAcks = -1
)

type ProducerConfig struct {
	Brokers []string
	// Topic must be specified either in ProducerConfig or in WriteMessage method. It can't be specified in both places.
	Topic        string
	MaxAttempts  int
	BatchSize    int
	BatchBytes   int
	BatchTimeout time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	RequiredAcks RequiredAcks
	Async        bool
	Log          bool
}

func (c *ProducerConfig) ToKafkaWriterConfig() kafka.WriterConfig {
	config := kafka.WriterConfig{
		Brokers:      c.Brokers,
		Topic:        c.Topic,
		MaxAttempts:  c.MaxAttempts,
		BatchSize:    c.BatchSize,
		BatchBytes:   c.BatchBytes,
		BatchTimeout: c.BatchTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		RequiredAcks: int(c.RequiredAcks),
		Async:        c.Async,
	}

	if c.Log {
		config.Logger = logging.GetLogger()
	}

	return config
}
