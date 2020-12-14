package binance

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/imroc/req"
)

type Client struct {
	url     string
	headers http.Header
}

func InitClient(url, apiKey string) Client {
	header := make(http.Header)
	if apiKey != "" {
		header.Set("apikey", apiKey)
	}
	return Client{
		url:     url,
		headers: header,
	}
}

func (c Client) Get(path string, params interface{}) (*req.Resp, error) {
	return req.Get(c.url+path, c.headers, params)
}

func (c Client) FetchLatestBlockNumber() (int64, error) {
	resp, err := c.Get("/api/v1/node-info", nil)
	if err != nil {
		return 0, err
	}
	var result NodeInfoResponse
	if err := resp.ToJSON(&result); err != nil {
		return 0, errors.Wrap(err, "URL: "+resp.Request().URL.String()+"; Code: "+resp.Response().Status)
	}
	return int64(result.SyncInfo.LatestBlockHeight), nil
}

func (c Client) FetchTransactionsInBlock(blockNumber int64) (TransactionsInBlockResponse, error) {
	resp, err := c.Get(fmt.Sprintf("/api/v2/transactions-in-block/%d", blockNumber), nil)
	if err != nil {
		return TransactionsInBlockResponse{}, err
	}
	var result TransactionsInBlockResponse
	if err := resp.ToJSON(&result); err != nil {
		return TransactionsInBlockResponse{}, errors.Wrap(err, "URL: "+resp.Request().URL.String()+"; Code: "+resp.Response().Status)
	}
	return result, nil
}

func (c Client) FetchTransactionsByAddressAndTokenID(address, tokenID string, limit int) ([]Tx, error) {
	startTime := strconv.Itoa(int(time.Now().AddDate(0, -3, 0).Unix() * 1000))
	params := url.Values{"address": {address}, "txAsset": {tokenID}, "startTime": {startTime}, "limit": {strconv.Itoa(limit)}}
	resp, err := c.Get("/api/v1/transactions", params)
	if err != nil {
		return nil, err
	}
	var result TransactionsInBlockResponse
	if err := resp.ToJSON(&result); err != nil {
		return nil, errors.Wrap(err, "URL: "+resp.Request().URL.String()+"; Code: "+resp.Response().Status)
	}
	return result.Tx, nil
}

func (c Client) FetchAccountMeta(address string) (AccountMeta, error) {
	resp, err := c.Get(fmt.Sprintf("/api/v1/account/%s", address), nil)
	if err != nil {
		return AccountMeta{}, err
	}
	var result AccountMeta
	if err := resp.ToJSON(&result); err != nil {
		return AccountMeta{}, errors.Wrap(err, "URL: "+resp.Request().URL.String()+"; Code: "+resp.Response().Status)
	}
	return result, nil
}

func (c Client) FetchTokens(limit int) (Tokens, error) {
	result := new(Tokens)
	query := url.Values{"limit": {strconv.Itoa(limit)}}
	resp, err := c.Get("/api/v1/tokens", query)
	if err != nil {
		return nil, err
	}
	if err := resp.ToJSON(&result); err != nil {
		return nil, errors.Wrap(err, "URL: "+resp.Request().URL.String()+"; Code: "+resp.Response().Status)
	}
	return *result, nil
}
