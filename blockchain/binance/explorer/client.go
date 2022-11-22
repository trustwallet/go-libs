package explorer

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/trustwallet/go-libs/client"
)

// Client is a binance explorer API client
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
	_, err = c.req.Execute(context.TODO(), client.NewReqBuilder().
		Method(http.MethodGet).
		WriteTo(&assets).
		PathStatic("/api/v1/assets").
		Query(params).
		Build())

	return assets, err
}
