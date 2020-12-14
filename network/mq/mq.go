package mq

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/streadway/amqp"
)

var (
	prefetchCount int
	amqpChan      *amqp.Channel
	conn          *amqp.Connection
)

type (
	Queue              string
	Consumer           func(amqp.Delivery)
	MessageChannel     <-chan amqp.Delivery
)

func Init(url string, prefetchChannelCount int) (err error) {
	conn, err = amqp.Dial(url)
	if err != nil {
		return err
	}
	amqpChan, err = conn.Channel()
	prefetchCount = prefetchChannelCount
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

func (q Queue) Declare() error {
	_, err := amqpChan.QueueDeclare(string(q), true, false, false, false, nil)
	return err
}

func (q Queue) Publish(body []byte) error {
	return amqpChan.Publish("", string(q), false, false, amqp.Publishing{
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

func (q Queue) RunConsumerWithCancel(consumer Consumer, async bool, ctx context.Context) {
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
			if async {
				go consumer(message)
			} else {
				consumer(message)
			}
		}
	}
}
func RestoreConnectionWorker(url string, queue Queue, timeout time.Duration) {
	log.Info("Run MQ RestoreConnectionWorker")
	for {
		if conn.IsClosed() {
			for {
				log.Warn("MQ is not available now")
				log.Warn("Trying to connect to MQ...")
				if err := Init(url, prefetchCount); err != nil {
					log.Warn("MQ is still unavailable")
					time.Sleep(timeout)
					continue
				}
				if err := queue.Declare(); err != nil {
					log.Warn("Can't declare queues:", queue)
					time.Sleep(timeout)
					continue
				} else {
					log.Info("MQ connection restored")
					break
				}
			}
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
