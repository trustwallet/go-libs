package types

type Asset struct {
	Id       string    `json:"asset"`
	Name     string    `json:"name"`
	Symbol   string    `json:"symbol"`
	Type     TokenType `json:"type"`
	Decimals uint      `json:"decimals"`
}
