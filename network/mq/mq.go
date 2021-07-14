package mq

import (
	"github.com/streadway/amqp"
)

var (
	conn *amqp.Connection
)

type Consumer interface {
	Callback(msg amqp.Delivery) error
}

type ConsumerDefaultCallback struct {
	Delivery func(amqp.Delivery) error
}

func (c ConsumerDefaultCallback) Callback(msg amqp.Delivery) error {
	return c.Delivery(msg)
}

func publish(amqpChan *amqp.Channel, exchange, queue string, body []byte) error {
	return amqpChan.Publish(exchange, queue, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
}

// queue

// Exchange
