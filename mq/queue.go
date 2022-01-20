package mq

type queue struct {
	name   QueueName
	client *Client
}

type Queue interface {
	Declare() error
	Publish(body []byte) error
	PublishWithConfig(body []byte, cfg PublishConfig) error
	Name() QueueName
}

func (q *queue) Name() QueueName {
	return q.name
}

func (q *queue) Declare() error {
	_, err := q.client.amqpChan.QueueDeclare(string(q.name), true, false, false, false, nil)
	return err
}

func (q *queue) Publish(body []byte) error {
	return publish(q.client.amqpChan, "", ExchangeKey(q.name), body)
}

func (q *queue) PublishWithConfig(body []byte, cfg PublishConfig) error {
	return publishWithConfig(q.client.amqpChan, "", ExchangeKey(q.name), body, cfg)
}

type PublishConfig struct {
	// MaxRetries defines the maximum number of retries after processing failures.
	MaxRetries *int
}
