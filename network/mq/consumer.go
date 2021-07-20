package mq

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type consumer struct {
	manager *Manager

	queue   Queue
	fn      func([]byte) error
	options ConsumerOptions

	tag string
}

type Consumer interface {
	Consume(ctx context.Context) error
}

func (c *consumer) Consume(ctx context.Context) error {
	messages, err := c.messageChannel()
	if err != nil {
		return err
	}
	for w := 1; w <= c.options.Workers; w++ {
		go c.consume(ctx, messages)
	}

	log.Infof("Started %d MQ consumer workers for queue %s", c.options.Workers, c.queue.Name())

	return nil
}

func (c *consumer) consume(ctx context.Context, messages <-chan amqp.Delivery) {
	for {
		select {
		case <-ctx.Done():
			log.Infof("Stopped consuming queue %s", c.queue.Name())
			return
		case msg := <-messages:
			if msg.Body == nil {
				continue
			}

			err := c.fn(msg.Body)
			if err != nil {
				log.Error(err)

				if c.options.RetryOnError {
					time.Sleep(c.options.RetryDelay)

					if err := msg.Reject(true); err != nil {
						log.Error(err)
					}

					continue
				}
			}

			if err := msg.Ack(false); err != nil {
				log.Error(err)
			}
		}
	}
}

func (c *consumer) messageChannel() (<-chan amqp.Delivery, error) {
	messageChannel, err := c.manager.amqpChan.Consume(
		string(c.queue.Name()),
		c.tag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("MQ issue" + err.Error() + " for queue: " + string(c.queue.Name()))
	}

	return messageChannel, nil
}
