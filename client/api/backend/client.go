package backend

import (
	"fmt"

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

func (c *Client) ValidateAssetInfo(req *AssetValidationReq) (result AssetValidationResp, err error) {
	err = c.req.Post(&result, "/v1/validate/asset_info", req)
	return result, err
}

func (c *Client) GetAssetInfo(assetID string) (result AssetInfoResp, err error) {
	err = c.req.Get(&result, fmt.Sprintf("/v1/assets/%s", assetID), nil)
	return result, err
}
