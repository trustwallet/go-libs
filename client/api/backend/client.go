package backend

import (
	"context"
	"fmt"
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
		WriteTo(&result).
		PathStatic(fmt.Sprintf("/v1/assets/%s", assetID)).
		Build())
	return result, err
}
