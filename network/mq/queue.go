package mq

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type queue struct {
	name         QueueName
	mq           *MQ
	consumer     Consumer
	consumerOpts ConsumerOptions
}

type Queue interface {
	Declare() error
	Publish(body []byte) error
	Name() QueueName
	WithConsumer(consumer Consumer, consumerOptions ConsumerOptions) Queue
	Consume(ctx context.Context) error
	Reconnect(ctx context.Context, amqpChan *amqp.Channel) error
}

type Consumer func([]byte) error

func (q *queue) Name() QueueName {
	return q.name
}

func (q *queue) Declare() error {
	_, err := q.mq.amqpChan.QueueDeclare(string(q.name), true, false, false, false, nil)
	return err
}

func (q *queue) Publish(body []byte) error {
	return publish(q.mq.amqpChan, "", ExchangeKey(q.name), body)
}

func (q *queue) Reconnect(ctx context.Context, amqpChan *amqp.Channel) error {
	q.mq.amqpChan = amqpChan

	err := q.Consume(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (q *queue) WithConsumer(consumer Consumer, consumerOptions ConsumerOptions) Queue {
	q.consumer = consumer
	q.consumerOpts = consumerOptions

	return q
}

func (q *queue) Consume(ctx context.Context) error {
	messages, err := q.messageChannel()
	if err != nil {
		return err
	}
	for w := 1; w <= q.consumerOpts.Workers; w++ {
		go q.consume(ctx, messages)
	}

	log.Infof("Started %d MQ consumer workers for queue %s", q.consumerOpts.Workers, q.name)

	return nil
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

			err := q.consumer(msg.Body)
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

func (q *queue) messageChannel() (<-chan amqp.Delivery, error) {
	messageChannel, err := q.mq.amqpChan.Consume(
		string(q.name),
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("MQ issue" + err.Error() + " for queue: " + string(q.name))
	}

	return messageChannel, nil
}
