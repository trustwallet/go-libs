package mq

import (
	"time"
)

type ConsumerOptions struct {
	Workers      int
	RetryOnError bool
	RetryDelay   time.Duration
}

func InitDefaultConsumerOptions(workers int) ConsumerOptions {
	return ConsumerOptions{
		Workers:      workers,
		RetryOnError: true,
		RetryDelay:   time.Second * 1,
	}
}
