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
