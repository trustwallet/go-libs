package types

import "github.com/trustwallet/golibs/coin"

type (
	TokenType string
)

const (
	Coin    TokenType = "coin"
	Gas     TokenType = "gas"
	ERC20   TokenType = "ERC20"
	ERC721  TokenType = "ERC721"
	ERC1155 TokenType = "ERC1155"
	BEP2    TokenType = "BEP2"
	BEP8    TokenType = "BEP8"
	BEP20   TokenType = "BEP20"
	TRC10   TokenType = "TRC10"
	ETC20   TokenType = "ETC20"
	POA20   TokenType = "POA20"
	TRC20   TokenType = "TRC20"
	TRC21   TokenType = "TRC21"
	CLO20   TokenType = "CLO20"
	GO20    TokenType = "G020"
	WAN20   TokenType = "WAN20"
	TT20    TokenType = "TT20"
	KAVA    TokenType = "KAVA"
	SPL     TokenType = "SPL"
	POLYGON TokenType = "POLYGON"
)

func GetEthereumTokenTypeByIndex(coinIndex uint) TokenType {
	var tokenType TokenType
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
	case coin.Smartchain().ID:
		tokenType = BEP20
	case coin.Solana().ID:
		tokenType = SPL
	case coin.Polygon().ID:
		tokenType = POLYGON
	default:
		tokenType = ERC20
	}
	return tokenType
}
