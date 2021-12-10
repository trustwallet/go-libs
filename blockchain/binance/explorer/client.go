package explorer

import (
	"net/url"
	"strconv"

	"github.com/trustwallet/go-libs/client"
)

// URL: https://explorer.binance.org
type Client struct {
	req client.Request
}

func InitClient(url string, errorHandler client.HttpErrorHandler) Client {
	request := client.InitJSONClient(url, errorHandler)

	return Client{
		req: request,
	}
}

func (c Client) FetchBep2Assets(page, rows int) (assets Bep2Assets, err error) {
	params := url.Values{
		"page": {strconv.Itoa(page)},
		"rows": {strconv.Itoa(rows)},
	}
	err = c.req.Get(&assets, "/api/v1/assets", params)

	return assets, err
}
