package api

type (
	TransactionsResponse struct {
		Total int  `json:"total"`
		Tx    []Tx `json:"txs"`
	}

	Type string

	Tx struct {
		Hash        string  `json:"hash"`
		BlockHeight int     `json:"blockHeight"`
		BlockTime   int64   `json:"blockTime"`
		Type        Type    `json:"type"`
		Fee         int     `json:"fee"`
		Code        int     `json:"code"`
		Source      int     `json:"source"`
		Sequence    int     `json:"sequence"`
		Memo        string  `json:"memo"`
		Log         string  `json:"log"`
		Data        string  `json:"data"`
		Asset       string  `json:"asset"`
		Amount      float64 `json:"amount"`
		FromAddr    string  `json:"fromAddr"`
		ToAddr      string  `json:"toAddr"`
	}
)
