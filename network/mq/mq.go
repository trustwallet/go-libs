package mq

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type (
	QueueName    string
	ExchangeName string
	ExchangeKey  string
)

const (
	reconnectionAttemptsNum = 5
	reconnectionTimeout     = time.Second * 30
)

type Manager struct {
	url      string
	conn     *amqp.Connection
	amqpChan *amqp.Channel

	connClients []ConnectionClient

	connCheckTimeout time.Duration
}

type Option func(m *Manager) error

func Open(url string, options ...Option) (*Manager, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	amqpChan, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	m := &Manager{
		url:      url,
		conn:     conn,
		amqpChan: amqpChan,

		connCheckTimeout: time.Second * 10, // default value
	}

	for _, opt := range options {
		err = opt(m)
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (m *Manager) Close() error {
	if m.amqpChan != nil {
		err := m.amqpChan.Close()
		if err != nil {
			log.Errorf("Close amqp channel: %v", err)
		}
	}

	if m.conn != nil && !m.conn.IsClosed() {
		err := m.conn.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) InitQueue(name QueueName) Queue {
	return &queue{
		name:    name,
		manager: m,
	}
}

func (m *Manager) InitExchange(name ExchangeName) Exchange {
	return &exchange{
		name:    name,
		manager: m,
	}
}

func (m *Manager) InitConsumer(queueName QueueName, options ConsumerOptions, fn func([]byte) error) Consumer {
	return &consumer{
		manager: m,
		queue:   m.InitQueue(queueName),
		fn:      fn,
		options: options,
	}
}

func (m *Manager) AddConnectionClient(connClient ConnectionClient) {
	m.connClients = append(m.connClients, connClient)
}

func (m *Manager) ListenConnection(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			err := m.Close()
			if err != nil {
				return fmt.Errorf("close mq: %v", err)
			}
			return nil
		default:
			err := m.checkConnection(ctx)
			if err != nil {
				return fmt.Errorf("check mq connection: %v", err)
			}

			time.Sleep(time.Second * 10)
		}
	}
}

func (m *Manager) checkConnection(ctx context.Context) error {
	if m.conn.IsClosed() {
		log.Warn("MQ connection lost")

		for i := 0; i < reconnectionAttemptsNum; i++ {
			time.Sleep(reconnectionTimeout)

			log.Info("Connecting to MQ... Attempt ", i+1)

			err := m.reconnect()
			if err != nil {
				log.Errorf("Reconnect: %v", err)
				continue
			}

			for _, connClient := range m.connClients {
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

func (m *Manager) reconnect() error {
	conn, err := amqp.Dial(m.url)
	if err != nil {
		return err
	}
	amqpChan, err := conn.Channel()
	if err != nil {
		return err
	}

	m.conn = conn
	m.amqpChan = amqpChan

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
