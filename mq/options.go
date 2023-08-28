package mq

import (
	"time"

	"github.com/trustwallet/go-libs/metrics"
)

type ConsumerOptions struct {
	Workers           int
	RetryOnError      bool
	RetryDelay        time.Duration
	PerformanceMetric metrics.PerformanceMetric

	// MaxRetries specifies the default number of retries for consuming a message.
	// A negative value is equal to infinite retries.
	MaxRetries int
}

func DefaultConsumerOptions(workers int) *ConsumerOptions {
	return &ConsumerOptions{
		Workers:           workers,
		RetryOnError:      true,
		RetryDelay:        time.Second,
		MaxRetries:        -1,
		PerformanceMetric: &metrics.NullablePerformanceMetric{},
	}
}

// Deprecated: We should not put prefetch limit at channel level. We need to set limit at consumer level
// This option no longer works to limit QoS globally.
//
// From rabbitMQ doc https://www.rabbitmq.com/consumer-prefetch.html
// Unfortunately the channel is not the ideal scope for this - since a single channel may consume from multiple queues,
// the channel and the queue(s) need to coordinate with each other for every message sent to ensure they don't go over
// the limit. This is slow on a single machine, and very slow when consuming across a cluster.
func OptionPrefetchLimit(limit int) Option {
	return func(m *Client) error {
		err := m.amqpChan.Qos(
			limit,
			0,
			true,
		)
		if err != nil {
			return err
		}

		return nil
	}
}

func OptionConnCheckTimeout(timeout time.Duration) Option {
	return func(m *Client) error {
		m.connCheckTimeout = timeout
		return nil
	}
}
