package mq

import (
	"time"

	"github.com/streadway/amqp"
)

type ConsumerOptions struct {
	Workers      int
	RetryOnError bool
	RetryDelay   time.Duration
}

func OptionPrefetchLimit(limit int) Option {
	return func(amqpChan *amqp.Channel) error {
		err := amqpChan.Qos(
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

func DefaultConsumerOptions(workers int) ConsumerOptions {
	return ConsumerOptions{
		Workers:      workers,
		RetryOnError: true,
		RetryDelay:   time.Second,
	}
}

func PoolOptionRetriesNumber(number int) PoolOption {
	return func(p *pool) {
		p.retriesNumber = number
	}
}
