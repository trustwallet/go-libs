package mq

import "github.com/streadway/amqp"

type Exchange struct {
	name     string
	amqpChan *amqp.Channel
}

func NewExchange(name string, amqpChan *amqp.Channel) *Exchange {
	return &Exchange{
		name:     name,
		amqpChan: amqpChan,
	}
}

func (e *Exchange) Declare(kind string) error {
	return e.amqpChan.ExchangeDeclare(e.name, kind, true, false, false, false, nil)
}

func (e *Exchange) Bind(queues []Queue) error {
	for _, q := range queues {
		err := e.amqpChan.QueueBind(q.Name(), "", e.name, false, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Exchange) BindWithKey(queues []Queue, key string) error {
	for _, q := range queues {
		err := e.amqpChan.QueueBind(q.Name(), key, e.name, false, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Exchange) Publish(body []byte) error {
	return publish(e.amqpChan, e.name, "", body)
}

func (e *Exchange) PublishWithKey(body []byte, key string) error {
	return publish(e.amqpChan, e.name, key, body)
}
