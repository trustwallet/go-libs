package mq

import (
	"context"
	"errors"
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
	url  string
	conn *amqp.Connection

	// This channel should only be used for management related operations, like declaring queues & exchanges
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
	if c.conn != nil && !c.conn.IsClosed() {
		err := c.conn.Close()
		if err != nil {
			return fmt.Errorf("close connection: %v", err)
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

func (c *Client) InitConsumer(queueName QueueName, options *ConsumerOptions, processor MessageProcessor) Consumer {
	return &consumer{
		client:           c,
		queue:            c.InitQueue(queueName),
		messageProcessor: processor,
		options:          options,
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

func (c *Client) initNotifyCloseListeners() (<-chan *amqp.Error, <-chan *amqp.Error) {
	return c.conn.NotifyClose(make(chan *amqp.Error)),
		c.amqpChan.NotifyClose(make(chan *amqp.Error))
}

func (c *Client) ListenConnection(ctx context.Context) error {
	log.Info("start listen connection")

	connErrCh, chanErrCh := c.initNotifyCloseListeners()

	for {
		select {

		case <-ctx.Done():
			err := c.Close()
			if err != nil {
				return fmt.Errorf("close mq: %v", err)
			}
			return nil

		case err, ok := <-chanErrCh:
			if !ok {
				// stop receiving from this channel to avoid multiple reads from closed channel before reconnected
				chanErrCh = nil
			}

			log.Info("received amqp channel close notification")
			if err != nil {
				log.Errorf("amqp channel closed with error: %v", err)
			}

			if c.conn.IsClosed() {
				break
			}

			// close connection to trigger reconnect logic
			// it will send notification to connErrCh
			if err := c.conn.Close(); err != nil {
				return fmt.Errorf("close connection: %v", err)
			}

		case err := <-connErrCh:
			log.Info("received connection close notification")
			if err != nil {
				log.Errorf("connection closed with error: %v", err)
			}

			if err := c.reconnectWithRetry(ctx); err != nil {
				return fmt.Errorf("check mq connection: %v", err)
			}

			// reassign listeners to new connection and channel
			connErrCh, chanErrCh = c.initNotifyCloseListeners()
		}
	}
}

func (c *Client) reconnectWithRetry(ctx context.Context) error {
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
	return publishWithConfig(amqpChan, exchange, key, body, PublishConfig{})
}

func publishWithConfig(amqpChan *amqp.Channel, exchange ExchangeName, key ExchangeKey, body []byte, cfg PublishConfig) error {
	headers := map[string]interface{}{}

	if cfg.MaxRetries != nil {
		headers[headerRemainingRetries] = *cfg.MaxRetries
	}

	var deliveryMode uint8
	if cfg.DeliveryMode == DeliveryModeTransient {
		deliveryMode = amqp.Transient
	} else {
		deliveryMode = amqp.Persistent
	}

	return amqpChan.Publish(string(exchange), string(key), false, false, amqp.Publishing{
		DeliveryMode: deliveryMode,
		ContentType:  "text/plain",
		Body:         body,
		Headers:      headers,
	})
}

type ConnectionClient interface {
	Reconnect(ctx context.Context) error
}

func (c *Client) HealthCheck() error {
	if c.conn.IsClosed() {
		return errors.New("connection is closed")
	}

	return nil
}
