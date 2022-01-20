package assetsmanager

import (
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

func (c *Client) GetTagValues() (result TagValuesResp, err error) {
	err = c.req.Get(&result, "/v1/values/tags", nil)
	return result, err
}
