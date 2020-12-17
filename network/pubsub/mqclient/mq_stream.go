package mqclient

import (
	"context"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"github.com/trustwallet/golibs/network/pubsub"
	"go.uber.org/atomic"
)

type Stream struct {
	consumer    *pubsub.Consumer
	pool        int
	client      pubsub.Client
	channel     *amqp.Channel
	isConnected *atomic.Bool
	isWriteOnly bool
}

func (s *Stream) Connect(cancelCtx context.Context) {
	s.isConnected.Store(true)
	for {
		if s.client.IsConnected() {
			break
		}
		time.Sleep(1 * time.Second)
	}
	_, err := s.channel.QueueDeclare((*s.consumer).GetQueue(), true, false, false, false, nil)
	if err != nil {
		fmt.Printf("Stream.Connect::QueueDeclare::%s %s", (*s.consumer).GetQueue(), err)
	}
	if s.isWriteOnly {
		return
	}
	messageChannel, err := s.channel.Consume(
		(*s.consumer).GetQueue(),
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		s.isConnected.Store(false)
		fmt.Printf("Stream.Connect::Consume::%s %s", (*s.consumer).GetQueue(), err)
		return
	}
	messages := make(chan amqp.Delivery)
	for w := 1; w <= s.pool; w++ {
		go s.worker(messages)
	}
	for {
		select {
		case <-cancelCtx.Done():
			return
		case msg, ok := <-messageChannel:
			if !ok {
				s.isConnected.Store(false)
				return
			}
			if msg.Body != nil {
				messages <- msg
			}
		}
	}
}
func (s *Stream) GetConsumer() *pubsub.Consumer {
	return s.consumer
}

func (s *Stream) GetClient() *pubsub.Client {
	return &s.client
}

func (s *Stream) IsConnected() bool {
	return s.isConnected.Load()
}

func (s *Stream) IsWriteOnly() bool {
	return s.isWriteOnly
}

func (s *Stream) worker(messages <-chan amqp.Delivery) {
	for msg := range messages {
		if (*s.consumer).Callback(msg) == nil {
			err := msg.Ack(false)
			if err != nil {
				fmt.Printf("Stream::ack::%s %s", (*s.consumer).GetQueue(), err)
			}
		}
	}
}
