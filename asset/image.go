package asset

import (
	"fmt"

	"github.com/trustwallet/golibs/coin"
)

func GetImageURL(endpoint, asset string) string {
	coinId, tokenId, err := ParseID(asset)
	if err != nil {
		return ""
	}
	if c, ok := coin.Coins[coinId]; ok {
		if len(tokenId) > 0 {
			return fmt.Sprintf("%s/blockchains/%s/assets/%s/logo.png", endpoint, c.Handle, tokenId)
		}
		return fmt.Sprintf("%s/blockchains/%s/info/logo.png", endpoint, c.Handle)
	}
	return ""
}
