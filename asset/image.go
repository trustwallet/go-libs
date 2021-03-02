package asset

import (
	"fmt"

	"github.com/trustwallet/golibs/coin"
)

func GetImageURL(endpoint, asset string) string {
	c, tokenId, err := ParseID(asset)
	if err != nil {
		return ""
	}
	if cc, ok := coin.Coins[c]; ok {
		if len(tokenId) > 0 {
			return fmt.Sprintf("%s/blockchains/%s/assets/%s/logo.png", endpoint, cc.Handle, tokenId)
		}
		return fmt.Sprintf("%s/blockchains/%s/info/logo.png", endpoint, cc.Handle)
	}
	return ""
}
