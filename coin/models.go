package coin

import "errors"

type ExternalCoin struct {
	Coin     uint   `json:"coin"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Decimals uint   `json:"decimals"`
}

func (c *Coin) External() *ExternalCoin {
	return &ExternalCoin{
		Coin:     c.ID,
		Name:     c.Name,
		Symbol:   c.Symbol,
		Decimals: c.Decimals,
	}
}

func GetCoinForId(id string) (Coin, error) {
	for _, c := range Coins {
		if c.Handle == id {
			return c, nil
		}
	}
	return Coin{}, errors.New("unknown id: " + id)
}
