package mq

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
)

var (
	prefetchCount int
	amqpChan      *amqp.Channel
	conn          *amqp.Connection
)

type Consumer interface {
	GetQueue() string
	Callback(msg amqp.Delivery) error
}

type (
	Queue          string
	MessageChannel <-chan amqp.Delivery
)

func Init(url string) (err error) {
	conn, err = amqp.Dial(url)
	if err != nil {
		return err
	}
	amqpChan, err = conn.Channel()
	return err
}

func Close() error {
	err := amqpChan.Close()
	if err != nil {
		log.Error(err)
	}

	return conn.Close()
}

func (mc MessageChannel) GetMessage() amqp.Delivery {
	return <-mc
}

func DeclareQueue(name string) error {
	_, err := amqpChan.QueueDeclare(name, true, false, false, false, nil)
	return err
}

func DeclareExchange(name, kind string) error {
	return amqpChan.ExchangeDeclare(name, kind, true, false, false, false, nil)
}

func Publish(exchange, queue string, body []byte) error {
	return amqpChan.Publish(exchange, queue, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
}

func (q Queue) GetMessageChannel() MessageChannel {
	messageChannel, err := amqpChan.Consume(
		string(q),
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("MQ issue " + err.Error())
	}

	err = amqpChan.Qos(
		prefetchCount,
		0,
		true,
	)
	if err != nil {
		log.Fatal("No qos limit ", err)
	}

	return messageChannel
}

func worker(messages <-chan amqp.Delivery, consumer Consumer) {
	for msg := range messages {
		consumer.Callback(msg)
	}
}

func (q Queue) RunConsumer(consumer Consumer, workers int, ctx context.Context) {
	messages := make(chan amqp.Delivery)
	for w := 1; w <= workers; w++ {
		go worker(messages, consumer)
	}
	messageChannel := q.GetMessageChannel()
	for {
		select {
		case <-ctx.Done():
			log.Info("Consumer stopped")
			return
		case message := <-messageChannel:
			if message.Body == nil {
				continue
			}
			consumer.Callback(message)
		}
	}
}
