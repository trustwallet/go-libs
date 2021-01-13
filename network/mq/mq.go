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

	DefaultConsumerOptions = ConsumerOptions{
		Workers:      1,
		RetryOnError: false,
		RetryDelay:   0,
	}
)

const ()

type Consumer interface {
	Callback(msg amqp.Delivery) error
}

type ConsumerOptions struct {
	Workers      int
	RetryOnError bool
	RetryDelay   time.Duration
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

func publish(exchange, queue string, body []byte) error {
	return amqpChan.Publish(exchange, queue, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
}

// Queue

func (q Queue) Declare() error {
	_, err := amqpChan.QueueDeclare(string(q), true, false, false, false, nil)
	return err
}

func (q Queue) Publish(body []byte) error {
	return publish("", string(q), body)
}

// Exchange
func (e Exchange) Declare(kind string) error {
	return amqpChan.ExchangeDeclare(string(e), kind, true, false, false, false, nil)
}

func (e Exchange) Bind(queues []Queue) error {
	for _, queue := range queues {
		err := amqpChan.QueueBind(string(queue), "", string(e), false, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e Exchange) Publish(body []byte) error {
	return publish(string(e), "", body)
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
		50,
		0,
		true,
	)
	if err != nil {
		log.Fatal("No qos limit ", err)
	}

	return messageChannel
}

func worker(messages <-chan amqp.Delivery, consumer Consumer, options ConsumerOptions) {
	for msg := range messages {
		err := consumer.Callback(msg)
		if err != nil {
			log.Error(err)
			if options.RetryOnError {
				if err := msg.Nack(false, true); err != nil {
					log.Error(err)
				}
				time.Sleep(options.RetryDelay)
			}
			return
		}
		if err := msg.Ack(false); err != nil {
			log.Error(err)
		}
	}
}

func (q Queue) RunConsumer(consumer Consumer, options ConsumerOptions, ctx context.Context) {
	messages := make(chan amqp.Delivery)
	for w := 1; w <= options.Workers; w++ {
		go worker(messages, consumer, options)
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
			messages <- message
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

func FatalWorker(timeout time.Duration) {
	log.Info("Run MQ FatalWorker")
	for {
		if conn.IsClosed() {
			log.Fatal("MQ is not available now")
		}
		time.Sleep(timeout)
	}
}
