package binance

import (
	"context"
	"net/http"
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
	request := client.InitJSONClient(url, errorHandler, client.WithExtraHeader("apikey", apiKey))
	return Client{
		req: request,
	}
}

func (c Client) FetchNodeInfo() (result NodeInfoResponse, err error) {
	_, err = c.req.Execute(context.TODO(), client.NewReqBuilder().
		Method(http.MethodGet).
		PathStatic("/api/v1/node-info").
		WriteTo(&result).
		Build())
	return result, err
}

func (c Client) FetchTransactionsInBlock(blockNumber int64) (result TransactionsInBlockResponse, err error) {
	_, err = c.req.Execute(context.TODO(), client.NewReqBuilder().
		Method(http.MethodGet).
		Pathf("api/v2/transactions-in-block/%d", blockNumber).
		WriteTo(&result).
		Build())
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
	_, err := c.req.Execute(context.TODO(), client.NewReqBuilder().
		Method(http.MethodGet).
		PathStatic("/api/v1/transactions").
		Query(params).
		WriteTo(&result).
		Build())
	return result.Tx, err
}

func (c Client) FetchAccountMeta(address string) (result AccountMeta, err error) {
	_, err = c.req.Execute(context.TODO(), client.NewReqBuilder().
		Method(http.MethodGet).
		Pathf("/api/v1/account/%s", address).
		WriteTo(&result).
		Build())
	return result, err
}

func (c Client) FetchTokens(limit int) (result Tokens, err error) {
	params := url.Values{"limit": {strconv.Itoa(limit)}}
	_, err = c.req.Execute(context.TODO(), client.NewReqBuilder().
		Method(http.MethodGet).
		PathStatic("/api/v1/tokens").
		Query(params).
		WriteTo(&result).
		Build())
	return result, err
}

func (c Client) FetchMarketPairs(limit int) (pairs []MarketPair, err error) {
	params := url.Values{"limit": {strconv.Itoa(limit)}}
	_, err = c.req.Execute(context.TODO(), client.NewReqBuilder().
		Method(http.MethodGet).
		PathStatic("/api/v1/markets").
		Query(params).
		WriteTo(&pairs).
		Build())
	return pairs, err
}
