package mq

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type consumer struct {
	client *Client

	queue   Queue
	fn      func(Message) error
	options ConsumerOptions
}

type Consumer interface {
	Start(ctx context.Context) error
	Reconnect(ctx context.Context) error
}

func (c *consumer) Start(ctx context.Context) error {
	messages, err := c.messageChannel()
	if err != nil {
		return fmt.Errorf("get message channel: %v", err)
	}
	for w := 1; w <= c.options.Workers; w++ {
		go c.consume(ctx, messages)
	}

	log.Infof("Started %d MQ consumer workers for queue %s", c.options.Workers, c.queue.Name())

	return nil
}

func (c *consumer) Reconnect(ctx context.Context) error {
	err := c.Start(ctx)
	if err != nil {
		return err
	}

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
	messageChannel, err := c.client.amqpChan.Consume(
		string(c.queue.Name()),
		"",
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
