package mq

type queue struct {
	name         QueueName
	manager      *Manager
	consumer     Consumer
	consumerOpts ConsumerOptions
}

type Queue interface {
	Declare() error
	Publish(body []byte) error
	Name() QueueName
}

func (q *queue) Name() QueueName {
	return q.name
}

func (q *queue) Declare() error {
	_, err := q.manager.amqpChan.QueueDeclare(string(q.name), true, false, false, false, nil)
	return err
}

func (q *queue) Publish(body []byte) error {
	return publish(q.manager.amqpChan, "", ExchangeKey(q.name), body)
}
