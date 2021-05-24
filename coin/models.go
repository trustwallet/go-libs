package coin

import "errors"

func GetCoinForId(id string) (Coin, error) {
	for _, c := range Coins {
		if c.Handle == id {
			return c, nil
		}
	}
	return Coin{}, errors.New("unknown id: " + id)
}
