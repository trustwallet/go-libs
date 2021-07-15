package mq

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Pool struct {
	connURL  string
	conn     *amqp.Connection
	amqpChan *amqp.Channel

	queues    []Queue
	queuesMux sync.RWMutex

	retriesNumber int
	timeout       time.Duration
}

type PoolOption func(p *Pool)

func InitPool(url string, options ...PoolOption) *Pool {
	conn, amqpChan, err := connect(url)
	if err != nil {
		log.Fatal("Failed to connect to AMQP: ", err)
	}

	pool := &Pool{
		connURL:  url,
		amqpChan: amqpChan,
		conn:     conn,

		timeout: time.Second * 10,
	}

	for _, opt := range options {
		opt(pool)
	}

	return pool
}

func (p *Pool) AddQueue(queue Queue) {
	p.queuesMux.Lock()
	defer p.queuesMux.Unlock()

	p.queues = append(p.queues, queue)
}

func (p *Pool) Close() error {
	err := p.amqpChan.Close()
	if err != nil {
		log.Print(err)
	}

	return p.conn.Close()
}

func (p *Pool) Consume(ctx context.Context) {
	for _, q := range p.queues {
		q.Consume(ctx)
	}

	for {
		if p.conn.IsClosed() {
			log.Warn("MQ connection lost")

			if p.retriesNumber == 0 {
				log.Fatal("MQ connection closed")
			}

			for i := 0; i < p.retriesNumber; i++ {
				log.Info("Connecting to MQ... Attempt ", i+1)
				conn, amqpChan, err := connect(p.connURL)
				if err == nil {
					p.reconnect(conn, amqpChan)
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

func (p *Pool) reconnect(conn *amqp.Connection, amqpChan *amqp.Channel) {
	p.conn = conn
	p.amqpChan = amqpChan

	for _, q := range p.queues {
		q.Reconnect(amqpChan)
	}
}
