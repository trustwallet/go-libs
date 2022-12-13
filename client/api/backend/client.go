package backend

import (
	"context"
	"net/http"

	"github.com/trustwallet/go-libs/client"
)

type Client struct {
	req client.Request
}

func InitClient(url string, errorHandler client.HttpErrorHandler) Client {
	return Client{
		req: client.InitJSONClient(url, errorHandler),
	}
}

func (c *Client) GetAssetInfo(assetID string) (result AssetInfoResp, err error) {
	_, err = c.req.Execute(context.TODO(), client.NewReqBuilder().
		Method(http.MethodGet).
		Pathf("/v1/assets/%s", assetID).
		WriteTo(&result).
		Build())
	return result, err
}
