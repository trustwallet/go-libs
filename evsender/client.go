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
	Name    string      `json:"name"`
	Payload interface{} `json:"payload"`
}

var senderClient *Client
var batchLimit = 100

func Init(url string, limit int) {
	senderClient = &Client{client.InitJSONClient(url, middleware.SentryErrorHandler)}
	batchLimit = limit
}

func (c Client) SendBatch(events []Event) (status Status, err error) {
	err = senderClient.Post(&status, "", events)
	return
}
