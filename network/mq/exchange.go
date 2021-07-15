package mq

import (
	"context"

	"github.com/streadway/amqp"
)

type exchange struct {
	name     ExchangeName
	amqpChan *amqp.Channel
}

func InitExchange(name ExchangeName, amqpChan *amqp.Channel) Exchange {
	return &exchange{
		name:     name,
		amqpChan: amqpChan,
	}
}

type Exchange interface {
	Declare(kind string) error
	Bind(queues []Queue) error
	BindWithKey(queues []Queue, key ExchangeKey) error
	Publish(body []byte) error
	PublishWithKey(body []byte, key ExchangeKey) error
	Reconnect(ctx context.Context, amqpChan *amqp.Channel)
}

func (e *exchange) Reconnect(_ context.Context, amqpChan *amqp.Channel) {
	e.amqpChan = amqpChan
}

func (e *exchange) Declare(kind string) error {
	return e.amqpChan.ExchangeDeclare(string(e.name), kind, true, false, false, false, nil)
}

func (e *exchange) Bind(queues []Queue) error {
	for _, q := range queues {
		err := e.amqpChan.QueueBind(string(q.Name()), "", string(e.name), false, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *exchange) BindWithKey(queues []Queue, key ExchangeKey) error {
	for _, q := range queues {
		err := e.amqpChan.QueueBind(string(q.Name()), string(key), string(e.name), false, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *exchange) Publish(body []byte) error {
	return publish(e.amqpChan, e.name, "", body)
}

func (e *exchange) PublishWithKey(body []byte, key ExchangeKey) error {
	return publish(e.amqpChan, e.name, key, body)
}
