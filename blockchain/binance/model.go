package binance

import "time"

type (
	NodeInfoResponse struct {
		SyncInfo struct {
			LatestBlockHeight int `json:"latest_block_height"`
		} `json:"sync_info"`
	}

	TransactionsInBlockResponse struct {
		BlockHeight int  `json:"blockHeight"`
		Tx          []Tx `json:"tx"`
	}

	TxType string

	Tx struct {
		TxHash          string            `json:"txHash"`
		BlockHeight     int               `json:"blockHeight"`
		TxType          TxType            `json:"txType"`
		TimeStamp       time.Time         `json:"timeStamp"`
		FromAddr        interface{}       `json:"fromAddr"`
		ToAddr          interface{}       `json:"toAddr"`
		Value           string            `json:"value"`
		TxAsset         string            `json:"txAsset"`
		TxFee           string            `json:"txFee"`
		OrderID         string            `json:"orderId,omitempty"`
		Code            int               `json:"code"`
		Data            string            `json:"data"`
		Memo            string            `json:"memo"`
		Source          int               `json:"source"`
		SubTransactions []SubTransactions `json:"subTransactions,omitempty"`
		Sequence        int               `json:"sequence"`
	}

	TransactionData struct {
		OrderData struct {
			Symbol      string `json:"symbol"`
			OrderType   string `json:"orderType"`
			Side        string `json:"side"`
			Price       string `json:"price"`
			Quantity    string `json:"quantity"`
			TimeInForce string `json:"timeInForce"`
			OrderID     string `json:"orderId"`
		} `json:"orderData"`
	}

	SubTransactions struct {
		TxHash      string `json:"txHash"`
		BlockHeight int    `json:"blockHeight"`
		TxType      string `json:"txType"`
		FromAddr    string `json:"fromAddr"`
		ToAddr      string `json:"toAddr"`
		TxAsset     string `json:"txAsset"`
		TxFee       string `json:"txFee"`
		Value       string `json:"value"`
	}

	AccountMeta struct {
		Balances []TokenBalance `json:"balances"`
	}

	TokenBalance struct {
		Free   string `json:"free"`
		Frozen string `json:"frozen"`
		Locked string `json:"locked"`
		Symbol string `json:"symbol"`
	}

	Tokens []Token

	Token struct {
		ContractAddress string `json:"contract_address"`
		Name            string `json:"name"`
		OriginalSymbol  string `json:"original_symbol"`
		Owner           string `json:"owner"`
		Symbol          string `json:"symbol"`
		TotalSupply     string `json:"total_supply"`
	}

	MarketPair struct {
		BaseAssetSymbol  string `json:"base_asset_symbol"`
		LotSize          string `json:"lot_size"`
		QuoteAssetSymbol string `json:"quote_asset_symbol"`
		TickSize         string `json:"tick_size"`
	}
)
