package mq

import (
	"os"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	Conn     *amqp.Connection
	AMQPChan *amqp.Channel
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

func Connect(url string) {
	var err error
	Conn, AMQPChan, err = connect(url)
	if err != nil {
		log.Fatal(err)
	}
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

func publish(amqpChan *amqp.Channel, exchange, queue string, body []byte) error {
	return amqpChan.Publish(exchange, queue, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
}

func QuitWorker(timeout time.Duration, quit chan<- os.Signal) {
	log.Info("Run CancelWorker")
	for {
		if Conn.IsClosed() {
			log.Error("MQ is not available now")
			quit <- syscall.SIGTERM
			return
		}
		time.Sleep(timeout)
	}
}

func FatalWorker(timeout time.Duration) {
	log.Info("Run MQ FatalWorker")
	for {
		if Conn.IsClosed() {
			log.Fatal("MQ is not available now")
		}
		time.Sleep(timeout)
	}
}
