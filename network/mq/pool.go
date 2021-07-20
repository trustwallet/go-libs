package mq

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var retriesDelay = []time.Duration{
	time.Second,
	time.Second * 3,
	time.Second * 10,
	time.Second * 15,
	time.Second * 30,
}

type pool struct {
	mq *Manager

	consumers    []Consumer
	consumersMux sync.RWMutex

	retriesNumber int
	timeout       time.Duration
}

type Pool interface {
	AddConsumer(queueName QueueName, tag string, consumerOptions ConsumerOptions, consumerFn func([]byte) error)
	Consume(ctx context.Context) error
}

type PoolOption func(p *pool)

func initPool(mq *Manager, options ...PoolOption) *pool {
	pool := &pool{
		mq:      mq,
		timeout: time.Second * 10, // default timeout value
	}

	for _, opt := range options {
		opt(pool)
	}

	return pool
}

func (p *pool) AddConsumer(queueName QueueName, tag string, consumerOptions ConsumerOptions, consumerFn func([]byte) error) {
	p.consumersMux.Lock()
	defer p.consumersMux.Unlock()

	p.consumers = append(p.consumers, p.mq.InitConsumer(queueName, tag, consumerOptions, consumerFn))
}

func (p *pool) Consume(ctx context.Context) error {
	for _, q := range p.consumers {
		err := q.Consume(ctx)
		if err != nil {
			return err
		}
	}

	for {
		if p.mq.conn.IsClosed() {
			log.Warn("MQ connection lost")

			for i := 0; i < len(retriesDelay); i++ {
				time.Sleep(retriesDelay[i])

				log.Info("Connecting to MQ... Attempt ", i+1)

				err := p.reconnect(ctx)
				if err == nil {
					log.Info("MQ Connection established")
					break
				}
				log.Error("Failed to establish MQ connection: ", err)

				if p.retriesNumber-i == 1 {
					log.Fatal("MQ is not available now")
				}
			}

		}

		time.Sleep(p.timeout)
	}
}

func (p *pool) reconnect(ctx context.Context) error {
	err := p.mq.reconnect()
	if err != nil {
		return err
	}

	for _, c := range p.consumers {
		err = c.Consume(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
