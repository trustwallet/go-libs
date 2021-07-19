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

type MQ struct {
	url      string
	conn     *amqp.Connection
	amqpChan *amqp.Channel
}

type Option func(amqpChan *amqp.Channel) error

func Open(url string, options ...Option) (*MQ, error) {
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

	return &MQ{
		url:      url,
		conn:     conn,
		amqpChan: amqpChan,
	}, nil
}

func (mq *MQ) Close() error {
	if mq.amqpChan != nil {
		err := mq.amqpChan.Close()
		if err != nil {
			log.Error(err)
		}
	}

	if mq.conn != nil && !mq.conn.IsClosed() {
		err := mq.conn.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (mq *MQ) NewQueue(name QueueName) Queue {
	return &queue{
		name: name,
		mq:   mq,
	}
}

func (mq *MQ) NewExchange(name ExchangeName) Exchange {
	return &exchange{
		name: name,
		mq:   mq,
	}
}

func (mq *MQ) NewPool(options ...PoolOption) Pool {
	return initPool(mq, options...)
}

func connect(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}
	amqpChan, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, amqpChan, nil
}

func publish(amqpChan *amqp.Channel, exchange ExchangeName, key ExchangeKey, body []byte) error {
	return amqpChan.Publish(string(exchange), string(key), false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
}
