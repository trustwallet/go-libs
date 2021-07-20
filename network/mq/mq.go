package mq

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type (
	QueueName    string
	ExchangeName string
	ExchangeKey  string
)

type Manager struct {
	url      string
	conn     *amqp.Connection
	amqpChan *amqp.Channel
}

type Option func(amqpChan *amqp.Channel) error

func Open(url string, options ...Option) (*Manager, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	amqpChan, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	for _, opt := range options {
		err = opt(amqpChan)
		if err != nil {
			return nil, err
		}
	}

	return &Manager{
		url:      url,
		conn:     conn,
		amqpChan: amqpChan,
	}, nil
}

func (m *Manager) Close() error {
	if m.amqpChan != nil {
		err := m.amqpChan.Close()
		if err != nil {
			log.Error(err)
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

func (m *Manager) InitConsumer(queueName QueueName, tag string, options ConsumerOptions, fn func([]byte) error) Consumer {
	return &consumer{
		manager: m,
		queue:   m.InitQueue(queueName),
		fn:      fn,
		options: options,
		tag:     tag,
	}
}

func (m *Manager) InitPool(options ...PoolOption) Pool {
	return initPool(m, options...)
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
