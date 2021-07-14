package mq

import (
	"time"
)

type ConsumerOptions struct {
	Workers       int
	PrefetchLimit int
	RetryOnError  bool
	RetryDelay    time.Duration
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
	return func(p *Pool) {
		p.retriesNumber = number
	}
}
