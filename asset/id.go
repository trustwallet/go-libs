package asset

import (
	"errors"
	"strconv"
	"strings"
)

type CoinType string

const (
	Coin  CoinType = "coin"
	Token CoinType = "token"

	coinPrefix  = 'c'
	tokenPrefix = 't'
)

func ParseID(id string) (uint, string, error) {
	rawResult := strings.Split(id, "_")
	resLen := len(rawResult)
	if resLen < 1 {
		return 0, "", errors.New("bad ID")
	}

	coin, err := FindCoinID(rawResult)
	if err != nil {
		return 0, "", errors.New("bad ID")
	}

	token := FindTokenID(rawResult)

	if token != "" {
		return coin, token, nil
	}

	return coin, "", nil
}

func BuildID(coin uint, token string) string {
	c := strconv.Itoa(int(coin))
	if token != "" {
		return string(coinPrefix) + c + "_" + string(tokenPrefix) + token
	}
	return string(coinPrefix) + c
}

func FindCoinID(words []string) (uint, error) {
	for _, w := range words {
		if len(w) == 0 {
			return 0, errors.New("empty coin")
		}

		if w[0] == coinPrefix {
			rawCoin := removeFirstChar(w)
			coin, err := strconv.Atoi(rawCoin)
			if err != nil {
				return 0, errors.New("bad coin")
			}
			return uint(coin), nil
		}
	}
	return 0, errors.New("no coin")
}

func FindTokenID(words []string) string {
	for _, w := range words {
		if w[0] == tokenPrefix {
			token := removeFirstChar(w)
			if len(token) > 0 {
				return token
			}
			return ""
		}
	}
	return ""
}

func removeFirstChar(input string) string {
	if len(input) <= 1 {
		return ""
	}
	return string([]rune(input)[1:])
}
