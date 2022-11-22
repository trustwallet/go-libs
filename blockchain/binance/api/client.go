package api

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/trustwallet/go-libs/client"
)

// Client is a binance API client
type Client struct {
	req client.Request
}

func InitClient(url string, errorHandler client.HttpErrorHandler) Client {
	request := client.InitJSONClient(url, errorHandler)

	return Client{
		req: request,
	}
}

func (c *Client) GetTransactionsByAddress(address string, limit int) ([]Tx, error) {
	startTime := strconv.Itoa(int(time.Now().AddDate(0, 0, -7).Unix() * 1000))
	endTime := strconv.Itoa(int(time.Now().Unix() * 1000))
	params := url.Values{
		"address":   {address},
		"startTime": {startTime},
		"endTime":   {endTime},
		"limit":     {strconv.Itoa(limit)},
	}

	var result TransactionsResponse

	_, err := c.req.Execute(context.TODO(), client.NewReqBuilder().
		Method(http.MethodGet).
		WriteTo(&result).
		PathStatic("bc/api/v1/txs").
		Query(params).
		Build())
	return result.Tx, err
}
