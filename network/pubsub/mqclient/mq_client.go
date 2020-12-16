package mqclient

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"github.com/trustwallet/golibs/network/pubsub"
	"go.uber.org/atomic"
)

var (
	reconnectDelay = 5 * time.Second
	resendDelay    = 5 * time.Second
)

type Client struct {
	uri           string
	conn          *amqp.Connection
	channel       *amqp.Channel
	prefetchCount int
	streams       []*pubsub.Stream
	ctx           context.Context
	notifyClose   chan *amqp.Error
	notifyConfirm chan amqp.Confirmation
	isConnected   *atomic.Bool
	alive         *atomic.Bool
	locker        sync.Mutex
}

func New(uri string, prefetchCount int, ctx context.Context) (client pubsub.Client) {
	client = &Client{
		uri:           uri,
		ctx:           ctx,
		alive:         atomic.NewBool(false),
		isConnected:   atomic.NewBool(false),
		prefetchCount: prefetchCount,
		streams:       []*pubsub.Stream{},
	}
	return client
}

func (c *Client) Connect() error {
	conn, err := amqp.Dial(c.uri)
	if err != nil {
		return errors.Wrapf(err, "Client.Connect::Dial: %s", c.uri)
	}
	ch, err := conn.Channel()
	if err != nil {
		return errors.Wrapf(err, "Client.Connect::Channel: %s", c.uri)
	}
	err = ch.Confirm(false)
	if err != nil {
		return errors.Wrapf(err, "Client.Connect::Confirm: %s", c.uri)
	}
	c.conn = conn
	c.channel = ch
	c.notifyClose = make(chan *amqp.Error)
	c.notifyConfirm = make(chan amqp.Confirmation)
	c.channel.NotifyClose(c.notifyClose)
	c.channel.NotifyPublish(c.notifyConfirm)
	c.isConnected.Store(true)
	err = c.channel.Qos(c.prefetchCount, 0, false)
	if err != nil {
		return errors.Wrapf(err, "Client.Connect::Qos: %s", c.uri)
	}

	for _, stream := range c.streams {
		go (*stream).Connect(c.ctx)
	}
	return nil
}

func (c *Client) Run() error {
	if c.conn == nil {
		return errors.New("connect firstly")
	}
	go c.handleReconnect()
	return nil
}

func (c *Client) IsConnected() bool {
	return c.isConnected.Load()
}

func (c *Client) AddStream(consumer *pubsub.Consumer, isWriteOnly bool) error {
	var stream pubsub.Stream = &Stream{
		consumer:    consumer,
		client:      c,
		isConnected: atomic.NewBool(false),
		isWriteOnly: isWriteOnly,
		channel:     c.channel,
	}
	go stream.Connect(c.ctx) // Try connect, if client isn't run it will wait run
	c.locker.Lock()
	defer c.locker.Unlock()
	c.streams = append(c.streams, &stream)
	return nil
}

func (c *Client) Push(queue string, data []byte, isWaitStream bool) error {
	if !c.isConnected.Load() {
		return errors.New("failed to push push: not connected")
	}
	if isWaitStream {
		stream := c.findStream(queue)
		for !(*stream).IsConnected() {
		}
	}
	for {
		err := c.PushUnsafe(queue, data)
		if err != nil {
			if err == pubsub.ErrDisconnected {
				continue
			}
			return errors.Wrapf(err, "Client.Push: %s", c.uri)
		}
		select {
		case confirm := <-c.notifyConfirm:
			if confirm.Ack {
				return nil
			}
		case <-time.After(resendDelay):
		}
	}
}

func (c *Client) PushUnsafe(queue string, data []byte) error {
	if !c.isConnected.Load() {
		return pubsub.ErrDisconnected
	}
	return c.channel.Publish(
		"",    // Exchange
		queue, // Routing key
		false, // Mandatory
		false, // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
}

func (c *Client) Close() error {
	err := c.channel.Close()
	if err != nil {
		return errors.Wrapf(err, "Client.Close::Channel: %s", c.uri)
	}

	err = c.conn.Close()
	if err != nil {
		return errors.Wrapf(err, "Client.Close::Connect: %s", c.uri)
	}
	return nil
}

func (c *Client) SetUri(uri string) {
	c.uri = uri
}

func (c *Client) handleReconnect() {
	for c.alive.Load() {
		c.isConnected.Store(false)
		for c.Connect() != nil {
			if !c.alive.Load() {
				return
			}
			select {
			case <-c.ctx.Done():
				return
			case <-time.After(reconnectDelay + time.Second):
				// Add metric
			}
		}
		select {
		case <-c.ctx.Done():
			return
		case <-c.notifyClose:
		}
	}
}

func (c *Client) findStream(queue string) *pubsub.Stream {
	c.locker.Lock()
	defer c.locker.Unlock()
	for _, s := range c.streams {
		if (*(*s).GetConsumer()).GetQueue() == queue {
			return s
		}
	}
	return nil
}
