package mq

import (
	"time"

	"github.com/streadway/amqp"
)

type ConsumerOptions struct {
	Workers       int
	PrefetchLimit int
	RetryOnError  bool
	RetryDelay    time.Duration
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

func InitDefaultConsumerOptions(workers int) ConsumerOptions {
	return ConsumerOptions{
		Workers:       workers,
		PrefetchLimit: 10,
		RetryOnError:  true,
		RetryDelay:    time.Second * 1,
	}
}

func PoolOptionRetriesNumber(number int) PoolOption {
	return func(p *pool) {
		p.retriesNumber = number
	}
}
