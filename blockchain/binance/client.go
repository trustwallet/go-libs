package binance

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/trustwallet/go-libs/client"
)

// Client is a binance dex API client
type Client struct {
	req client.Request
}

func InitClient(url, apiKey string, errorHandler client.HttpErrorHandler) Client {
	request := client.InitJSONClient(url, errorHandler)
	request.AddHeader("apikey", apiKey)
	return Client{
		req: request,
	}
}

func (c Client) FetchNodeInfo() (result NodeInfoResponse, err error) {
	err = c.req.Get(&result, "/api/v1/node-info", nil)
	return result, err
}

func (c Client) FetchTransactionsInBlock(blockNumber int64) (result TransactionsInBlockResponse, err error) {
	err = c.req.Get(&result, fmt.Sprintf("api/v2/transactions-in-block/%d", blockNumber), nil)
	return result, err
}

func (c Client) FetchTransactionsByAddressAndTokenID(address, tokenID string, limit int) ([]Tx, error) {
	startTime := strconv.Itoa(int(time.Now().AddDate(0, -3, 0).Unix() * 1000))
	params := url.Values{
		"address":   {address},
		"txAsset":   {tokenID},
		"startTime": {startTime},
		"limit":     {strconv.Itoa(limit)},
	}
	var result TransactionsInBlockResponse
	err := c.req.Get(&result, "/api/v1/transactions", params)
	return result.Tx, err
}

func (c Client) FetchAccountMeta(address string) (result AccountMeta, err error) {
	err = c.req.Get(&result, fmt.Sprintf("/api/v1/account/%s", address), nil)
	return result, err
}

func (c Client) FetchTokens(limit int) (result Tokens, err error) {
	params := url.Values{"limit": {strconv.Itoa(limit)}}
	err = c.req.Get(&result, "/api/v1/tokens", params)
	return result, err
}

func (c Client) FetchMarketPairs(limit int) (pairs []MarketPair, err error) {
	params := url.Values{"limit": {strconv.Itoa(limit)}}
	err = c.req.Get(&pairs, "/api/v1/markets", params)
	return pairs, err
}
