package mq

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type queue struct {
	name     string
	amqpChan *amqp.Channel

	consumer     Consumer
	consumerOpts ConsumerOptions
}

func InitQueue(name string, amqpChan *amqp.Channel) *queue {
	return &queue{
		name:     name,
		amqpChan: amqpChan,
	}
}

type Queue interface {
	Declare() error
	Publish(body []byte) error
	Name() string
	WithConsumer(consumer Consumer, consumerOptions ConsumerOptions) Queue
	Consume(ctx context.Context)
	Reconnect(amqpChan *amqp.Channel)
}

func (q *queue) Name() string {
	return q.name
}

func (q *queue) Declare() error {
	_, err := q.amqpChan.QueueDeclare(q.name, true, false, false, false, nil)
	return err
}

func (q *queue) Publish(body []byte) error {
	return publish(q.amqpChan, "", q.name, body)
}

func (q *queue) Reconnect(amqpChan *amqp.Channel) {
	q.amqpChan = amqpChan
}

func (q *queue) WithConsumer(consumer Consumer, consumerOptions ConsumerOptions) Queue {
	q.consumer = consumer
	q.consumerOpts = consumerOptions

	return q
}

func (q *queue) Consume(ctx context.Context) {
	messages := q.messageChannel()
	for w := 1; w <= q.consumerOpts.Workers; w++ {
		go q.consume(ctx, messages)
	}

	log.Infof("Started %d MQ consumer workers for queue %s", q.consumerOpts.Workers, q.name)
}

func (q *queue) consume(ctx context.Context, messages <-chan amqp.Delivery) {
	for {
		select {
		case <-ctx.Done():
			log.Infof("Stopped consuming queue %s", q.name)
			return
		case msg := <-messages:
			if msg.Body == nil {
				continue
			}

			err := q.consumer.Callback(msg)
			if err != nil {
				log.Error(err)
			}
			if err != nil && q.consumerOpts.RetryOnError {
				time.Sleep(q.consumerOpts.RetryDelay)
				if err := msg.Reject(true); err != nil {
					log.Error(err)
				}
			} else {
				if err := msg.Ack(false); err != nil {
					log.Error(err)
				}
			}
		}
	}
}

func (q *queue) messageChannel() <-chan amqp.Delivery {
	messageChannel, err := q.amqpChan.Consume(
		q.name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("MQ issue" + err.Error() + " for queue: " + q.name)
	}

	err = q.amqpChan.Qos(
		q.consumerOpts.PrefetchLimit,
		0,
		true,
	)
	if err != nil {
		log.Error("No qos limit ", err)
	}

	return messageChannel
}
