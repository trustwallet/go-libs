package mq

import (
	"context"
	"os"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
)

var (
	prefetchCount int
	amqpChan      *amqp.Channel
	conn          *amqp.Connection
)

type Consumer interface {
	Callback(msg amqp.Delivery) error
}

type (
	Queue          string
	Exchange       string
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

type ConsumerDefaultCallback struct {
	Delivery func(amqp.Delivery) error
}

func (c ConsumerDefaultCallback) Callback(msg amqp.Delivery) error {
	return c.Delivery(msg)
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

func (q Queue) Declare() error {
	_, err := amqpChan.QueueDeclare(string(q), true, false, false, false, nil)
	return err
}

func (e Exchange) Declare(kind string) error {
	return amqpChan.ExchangeDeclare(string(e), kind, true, false, false, false, nil)
}

func (q Queue) Publish(body []byte) error {
	return publish("", string(q), body)
}

func (e Exchange) Publish(body []byte) error {
	return publish(string(e), "", body)
}

func publish(exchange, queue string, body []byte) error {
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
			err := consumer.Callback(message)
			if err != nil {
				log.Error(err)
				continue
			}
			if err := message.Ack(false); err != nil {
				log.Error(err)
			}
		}
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
