package mq

import (
	"context"
	"os"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Pool struct {
	conn     *amqp.Connection
	amqpChan *amqp.Channel

	queues    []Queue
	queuesMux sync.RWMutex

	retriesNumber int
	timeout       time.Duration
}

type PoolOption func(p *Pool)

func InitPool(url string, options ...PoolOption) *Pool {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	amqpChan, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	pool := &Pool{
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

	return conn.Close()
}

func (p *Pool) Consume(ctx context.Context) {
	for _, q := range p.queues {
		q.RunConsumer(ctx)
	}

	for {
		if conn.IsClosed() {
			log.Fatal("MQ connection closed")
		}

		time.Sleep(p.timeout)
	}
}

func QuitWorker(timeout time.Duration, quit chan<- os.Signal) {
	log.Info("Run CancelWorker")
	for {
		if conn.IsClosed() {
			log.Error("MQ is not available now")
			quit <- syscall.SIGTERM
			return
		}
		time.Sleep(timeout)
	}
}

func FatalWorker(timeout time.Duration) {
	log.Info("Run MQ FatalWorker")
	for {
		if conn.IsClosed() {
			log.Fatal("MQ is not available now")
		}
		time.Sleep(timeout)
	}
}

func RetryWorker(timeout time.Duration, retriesAmount int, url string) {
	log.Info("Run MQ RetryWorker")
	for {
		if conn.IsClosed() {
			log.Warn("MQ connection lost")
			for i := 0; i < retriesAmount; i++ {
				log.Info("Connecting to MQ... Attempt ", i+1)
				//err := Init(url)
				//if err == nil {
				//	break
				//}
				//log.Error("Failed to establish MQ connection: ", err)

				if retriesAmount-i == 1 {
					log.Fatal("MQ is not available now")
				}

				time.Sleep(timeout)
			}

			log.Info("MQ connection established")
		}
		time.Sleep(timeout)
	}
}
