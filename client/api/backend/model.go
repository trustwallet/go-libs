package backend

type (
	AssetInfoResp struct {
		Name     string `json:"name"`
		Symbol   string `json:"symbol"`
		Type     string `json:"type"`
		Decimals int    `json:"decimals"`
		AssetID  string `json:"asset_id"`
	}
)
