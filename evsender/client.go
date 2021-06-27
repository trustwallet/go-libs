package evsender

import (
	"github.com/trustwallet/golibs/client"
	"github.com/trustwallet/golibs/network/middleware"
)

type Client struct {
	client.Request
}

type Status struct {
	Status bool `json:"status"`
}

type Event struct {
	Name    string  `json:"name"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	CreatedAt uint              `json:"created_at"`
	Payload   map[string]string `json:"payload"`
}

var senderClient *Client
var batchLimit = 100

func Init(url string, limit int) {
	senderClient = &Client{client.InitJSONClient(url, middleware.SentryErrorHandler)}
	batchLimit = limit
}

func (c Client) SendBatch(events interface{}) (status Status, err error) {
	err = senderClient.Post(&status, "", events)
	return
}
