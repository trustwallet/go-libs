package mq

type exchange struct {
	name    ExchangeName
	manager *Manager
}

type Exchange interface {
	Declare(kind string) error
	Bind(queues []Queue) error
	BindWithKey(queues []Queue, key ExchangeKey) error
	Publish(body []byte) error
	PublishWithKey(body []byte, key ExchangeKey) error
}

func (e *exchange) Declare(kind string) error {
	return e.manager.amqpChan.ExchangeDeclare(string(e.name), kind, true, false, false, false, nil)
}

func (e *exchange) Bind(queues []Queue) error {
	for _, q := range queues {
		err := e.manager.amqpChan.QueueBind(string(q.Name()), "", string(e.name), false, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *exchange) BindWithKey(queues []Queue, key ExchangeKey) error {
	for _, q := range queues {
		err := e.manager.amqpChan.QueueBind(string(q.Name()), string(key), string(e.name), false, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *exchange) Publish(body []byte) error {
	return publish(e.manager.amqpChan, e.name, "", body)
}

func (e *exchange) PublishWithKey(body []byte, key ExchangeKey) error {
	return publish(e.manager.amqpChan, e.name, key, body)
}
