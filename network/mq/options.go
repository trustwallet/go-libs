package mq

import (
	"time"
)

type ConsumerOptions struct {
	Workers      int
	RetryOnError bool
	RetryDelay   time.Duration
}

func DefaultConsumerOptions(workers int) ConsumerOptions {
	return ConsumerOptions{
		Workers:      workers,
		RetryOnError: true,
		RetryDelay:   time.Second,
	}
}

func OptionPrefetchLimit(limit int) Option {
	return func(m *Manager) error {
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
	return func(m *Manager) error {
		m.connCheckTimeout = timeout
		return nil
	}
}
