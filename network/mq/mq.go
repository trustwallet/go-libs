package mq

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type (
	QueueName    string
	ExchangeName string
	ExchangeKey  string
	Message      []byte
)

const (
	reconnectionAttemptsNum = 5
	reconnectionTimeout     = time.Second * 30
)

type Client struct {
	url      string
	conn     *amqp.Connection
	amqpChan *amqp.Channel

	connClients []ConnectionClient

	connCheckTimeout time.Duration
}

type Option func(c *Client) error

func Connect(url string, options ...Option) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	amqpChan, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	c := &Client{
		url:      url,
		conn:     conn,
		amqpChan: amqpChan,

		connCheckTimeout: time.Second * 10, // default value
	}

	for _, opt := range options {
		err = opt(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Client) Close() error {
	if c.amqpChan != nil {
		err := c.amqpChan.Close()
		if err != nil {
			log.Errorf("Close amqp channel: %v", err)
		}
	}

	if c.conn != nil && !c.conn.IsClosed() {
		err := c.conn.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) InitQueue(name QueueName) Queue {
	return &queue{
		name:   name,
		client: c,
	}
}

func (c *Client) InitExchange(name ExchangeName) Exchange {
	return &exchange{
		name:   name,
		client: c,
	}
}

func (c *Client) InitConsumer(queueName QueueName, options ConsumerOptions, fn func(message Message) error) Consumer {
	return &consumer{
		client:  c,
		queue:   c.InitQueue(queueName),
		fn:      fn,
		options: options,
	}
}

func (c *Client) StartConsumers(ctx context.Context, consumers ...Consumer) error {
	for _, consumer := range consumers {
		err := consumer.Start(ctx)
		if err != nil {
			return err
		}

		c.AddConnectionClient(consumer)
	}

	return nil
}

func (c *Client) AddConnectionClient(connClient ConnectionClient) {
	c.connClients = append(c.connClients, connClient)
}

func (c *Client) ListenConnectionAsync(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		err := c.ListenConnection(ctx)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()
}

func (c *Client) ListenConnection(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			err := c.Close()
			if err != nil {
				return fmt.Errorf("close mq: %v", err)
			}
			return nil
		default:
			err := c.checkConnection(ctx)
			if err != nil {
				return fmt.Errorf("check mq connection: %v", err)
			}

			time.Sleep(time.Second * 10)
		}
	}
}

func (c *Client) checkConnection(ctx context.Context) error {
	if c.conn.IsClosed() {
		log.Warn("MQ connection lost")

		for i := 0; i < reconnectionAttemptsNum; i++ {
			time.Sleep(reconnectionTimeout)

			log.Info("Connecting to MQ... Attempt ", i+1)

			err := c.reconnect()
			if err != nil {
				log.Errorf("Reconnect: %v", err)
				continue
			}

			for _, connClient := range c.connClients {
				err = connClient.Reconnect(ctx)
				if err != nil {
					log.Errorf("Reconnect for %+v: %v", connClient, err)
					continue
				}
			}

			log.Info("MQ connection established")
			return nil
		}

		return fmt.Errorf("failed to establish MQ connection")
	}

	return nil
}

func (c *Client) reconnect() error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return err
	}
	amqpChan, err := conn.Channel()
	if err != nil {
		return err
	}

	c.conn = conn
	c.amqpChan = amqpChan

	return nil
}

func publish(amqpChan *amqp.Channel, exchange ExchangeName, key ExchangeKey, body []byte) error {
	return amqpChan.Publish(string(exchange), string(key), false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
}

type ConnectionClient interface {
	Reconnect(ctx context.Context) error
}
