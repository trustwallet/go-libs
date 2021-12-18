package backend

type (
	AssetValidationReq struct {
		Name        *string `json:"name,omitempty"`
		Symbol      *string `json:"symbol,omitempty"`
		Type        *string `json:"type,omitempty"`
		Decimals    *int    `json:"decimals,omitempty"`
		Description *string `json:"description,omitempty"`
		Website     *string `json:"website,omitempty"`
		Explorer    *string `json:"explorer,omitempty"`
		Research    string  `json:"research,omitempty"`
		Status      *string `json:"status,omitempty"`
		ID          *string `json:"id,omitempty"`
		Links       []struct {
			Name *string `json:"name,omitempty"`
			URL  *string `json:"url,omitempty"`
		} `json:"links,omitempty"`
		ShortDesc     *string  `json:"short_desc,omitempty"`
		Audit         *string  `json:"audit,omitempty"`
		AuditReport   *string  `json:"audit_report,omitempty"`
		Tags          []string `json:"tags,omitempty"`
		Code          *string  `json:"code,omitempty"`
		Ticker        *string  `json:"ticker,omitempty"`
		ExplorerEth   *string  `json:"explorer-ETH,omitempty"`
		Address       *string  `json:"address,omitempty"`
		Twitter       *string  `json:"twitter,omitempty"`
		CoinMarketcap *string  `json:"coinmarketcap,omitempty"`
		DataSource    *string  `json:"data_source,omitempty"`
	}

	AssetValidationResp struct {
		Status string  `json:"status"`
		Errors []Error `json:"errors"`
	}

	Error struct {
		Message string `json:"message"`
	}
)

type (
	AssetInfoResp struct {
		Name     string `json:"name"`
		Symbol   string `json:"symbol"`
		Type     string `json:"type"`
		Decimals int    `json:"decimals"`
		AssetID  string `json:"asset_id"`
	}
)
