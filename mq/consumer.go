package mq

import (
	"context"
	"fmt"
	"time"

	"github.com/trustwallet/go-libs/pkg/nullable"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const headerRemainingRetries = "x-remaining-retries"

type consumer struct {
	client *Client

	queue   Queue
	fn      func(Message) error
	options *ConsumerOptions

	messages <-chan amqp.Delivery
	stopChan chan struct{}
}

type Consumer interface {
	Start(ctx context.Context) error
	Reconnect(ctx context.Context) error
}

func (c *consumer) Start(ctx context.Context) error {
	c.stopChan = make(chan struct{})

	var err error
	c.messages, err = c.messageChannel()
	if err != nil {
		return fmt.Errorf("get message channel: %v", err)
	}
	for w := 1; w <= c.options.Workers; w++ {
		go c.consume(ctx)
	}

	log.Infof("Started %d MQ consumer workers for queue %s", c.options.Workers, c.queue.Name())

	return nil
}

func (c *consumer) Reconnect(ctx context.Context) error {
	c.messages = nil
	if c.stopChan != nil {
		close(c.stopChan)
	}

	err := c.Start(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *consumer) consume(ctx context.Context) {
	metric := c.options.PerformanceMetric
	queueName := string(c.queue.Name())
	lvs := []string{queueName}

	for {
		select {
		case <-ctx.Done():
			log.Infof("Finished consuming queue %s", queueName)
			return
		case <-c.stopChan:
			log.Infof("Force stopped consuming queue %s", queueName)
			return
		case msg := <-c.messages:
			if msg.Body == nil {
				continue
			}

			t, _ := metric.Start(lvs)
			err := c.fn(msg.Body)
			metric.Duration(t, lvs)

			if err != nil {
				metric.Failure(lvs)
				log.Error(err)
			} else {
				metric.Success(lvs)
			}

			if err != nil && c.options.RetryOnError {
				time.Sleep(c.options.RetryDelay)
				remainingRetries := c.getRemainingRetries(msg)

				switch {
				case remainingRetries > 0:
					if err := c.queue.PublishWithConfig(msg.Body, PublishConfig{
						MaxRetries: nullable.Int(int(remainingRetries - 1)),
					}); err != nil {
						log.Error(err)
					}
				case remainingRetries == 0:
					break
				default:
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

func (c *consumer) getRemainingRetries(delivery amqp.Delivery) int32 {
	remainingRetriesRaw, exists := delivery.Headers[headerRemainingRetries]
	if !exists {
		return int32(c.options.MaxRetries)
	}

	remainingRetries, ok := remainingRetriesRaw.(int32)
	if !ok {
		return int32(c.options.MaxRetries)
	}

	return remainingRetries
}
