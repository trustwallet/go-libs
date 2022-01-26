package mq

import (
	"time"
)

type ConsumerOptions struct {
	Workers      int
	RetryOnError bool
	RetryDelay   time.Duration

	// MaxRetries specifies the default number of retries for consuming a message.
	// A negative value is equal to infinite retries.
	MaxRetries int
}

func DefaultConsumerOptions(workers int) ConsumerOptions {
	return ConsumerOptions{
		Workers:      workers,
		RetryOnError: true,
		RetryDelay:   time.Second,
		MaxRetries:   -1,
	}
}

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
