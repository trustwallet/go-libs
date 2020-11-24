package tokentype

import "github.com/trustwallet/golibs/coin"

type (
	Type       string
)

const (
	ERC20 Type = "ERC20"
	BEP2  Type = "BEP2"
	BEP8  Type = "BEP8"
	BEP20 Type = "BEP20"
	TRC10 Type = "TRC10"
	ETC20 Type = "ETC20"
	POA20 Type = "POA20"
	TRC20 Type = "TRC20"
	TRC21 Type = "TRC21"
	CLO20 Type = "CLO20"
	GO20  Type = "G020"
	WAN20 Type = "WAN20"
	TT20  Type = "TT20"
	KAVA Type = "KAVA"
)

func GetEthereumTokenTypeByIndex(coinIndex uint) Type {
	var tokenType Type
	switch coinIndex {
	case coin.Ethereum().ID:
		tokenType = ERC20
	case coin.Classic().ID:
		tokenType = ETC20
	case coin.Poa().ID:
		tokenType = POA20
	case coin.Callisto().ID:
		tokenType = CLO20
	case coin.Wanchain().ID:
		tokenType = WAN20
	case coin.Thundertoken().ID:
		tokenType = TT20
	case coin.Gochain().ID:
		tokenType = GO20
	case coin.Tomochain().ID:
		tokenType = TRC21
	case coin.Bsc().ID, coin.Smartchain().ID:
		tokenType = BEP20
	default:
		tokenType = ERC20
	}
	return tokenType
}

