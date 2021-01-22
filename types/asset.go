package types

type Asset struct {
	Id       string    `json:"asset"`
	Name     string    `json:"name"`
	Symbol   string    `json:"symbol"`
	Type     TokenType `json:"type"`
	Decimals uint      `json:"decimals"`
}

func GetAssetsIds(assets []Token) []string {
	assetIds := make([]string, 0)
	for _, asset := range assets {
		assetIds = append(assetIds, asset.AssetId())
	}
	return assetIds
}
