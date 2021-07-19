package mq

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type pool struct {
	mq *MQ

	queues    []Queue
	queuesMux sync.RWMutex

	exchanges    []Exchange
	exchangesMux sync.RWMutex

	retriesNumber int
	timeout       time.Duration
}

type Pool interface {
	AddQueue(queue Queue)
	AddExchange(exchange Exchange)
	Start(ctx context.Context)
}

type PoolOption func(p *pool)

func initPool(mq *MQ, options ...PoolOption) *pool {
	pool := &pool{
		mq:      mq,
		timeout: time.Second * 10, // default timeout value
	}

	for _, opt := range options {
		opt(pool)
	}

	return pool
}

func (p *pool) AddQueue(queue Queue) {
	p.queuesMux.Lock()
	defer p.queuesMux.Unlock()

	p.queues = append(p.queues, queue)
}

func (p *pool) AddExchange(exchange Exchange) {
	p.exchangesMux.Lock()
	defer p.exchangesMux.Unlock()

	p.exchanges = append(p.exchanges, exchange)
}

func (p *pool) Start(ctx context.Context) {
	for _, q := range p.queues {
		q.Consume(ctx)
	}

	for {
		if p.mq.conn.IsClosed() {
			log.Warn("MQ connection lost")

			if p.retriesNumber == 0 {
				log.Fatal("MQ connection closed")
			}

			for i := 0; i < p.retriesNumber; i++ {
				log.Info("Connecting to MQ... Attempt ", i+1)
				conn, amqpChan, err := connect(p.mq.url)
				if err == nil {
					p.reconnect(ctx, conn, amqpChan)
					log.Info("MQ Connection established")
					break
				}
				log.Error("Failed to establish MQ connection: ", err)

				if p.retriesNumber-i == 1 {
					log.Fatal("MQ is not available now")
				}

				time.Sleep(p.timeout)
			}

		}

		time.Sleep(p.timeout)
	}
}

func (p *pool) reconnect(ctx context.Context, conn *amqp.Connection, amqpChan *amqp.Channel) {
	p.mq.conn = conn
	p.mq.amqpChan = amqpChan

	for _, q := range p.queues {
		q.Reconnect(ctx, amqpChan)
	}

	for _, e := range p.exchanges {
		e.Reconnect(ctx, amqpChan)
	}
}
