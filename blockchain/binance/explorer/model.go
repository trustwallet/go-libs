package explorer

type (
	Bep2Asset struct {
		Asset       string `json:"asset"`
		Name        string `json:"name"`
		AssetImg    string `json:"assetImg"`
		MappedAsset string `json:"mappedAsset"`
		Decimals    int    `json:"decimals"`
	}

	Bep2Assets struct {
		AssetInfoList []Bep2Asset `json:"assetInfoList"`
	}
)
