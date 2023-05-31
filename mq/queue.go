package mq

import "fmt"

type queue struct {
	name   QueueName
	client *Client
}

type Queue interface {
	Declare() error
	DeclareWithConfig(cfg DeclareConfig) error
	Publish(body []byte) error
	PublishWithConfig(body []byte, cfg PublishConfig) error
	Name() QueueName
}

func (q *queue) Name() QueueName {
	return q.name
}

func (q *queue) Declare() error {
	return q.DeclareWithConfig(DeclareConfig{Durable: true})
}

func (q *queue) DeclareWithConfig(cfg DeclareConfig) error {
	_, err := q.client.amqpChan.QueueDeclare(
		string(q.name),
		cfg.Durable,
		cfg.AutoDelete,
		cfg.Exclusive,
		cfg.NoWait,
		cfg.Args,
	)
	return err
}

func (q *queue) Publish(body []byte) error {
	return publish(q.client.amqpChan, "", ExchangeKey(q.name), body)
}

func (q *queue) PublishWithConfig(body []byte, cfg PublishConfig) error {
	return publishWithConfig(q.client.amqpChan, "", ExchangeKey(q.name), body, cfg)
}

func (q *queue) HealthCheck() error {
	if err := q.client.HealthCheck(); err != nil {
		return fmt.Errorf("client health check: %v", err)
	}

	return nil
}

type DeclareConfig struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       map[string]interface{}
}

type DeliveryMode uint8

const (
	DeliveryModeTransient  DeliveryMode = 1
	DeliveryModePersistent DeliveryMode = 2
)

type PublishConfig struct {
	// MaxRetries defines the maximum number of retries after processing failures.
	// Overrides the value of consumer's config.
	MaxRetries   *int
	DeliveryMode DeliveryMode
}
