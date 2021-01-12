package types

import (
	"testing"

	"github.com/trustwallet/golibs/coin"
)

func TestGetEthereumTokenTypeByIndex(t *testing.T) {
	type args struct {
		coinIndex uint
	}
	tests := []struct {
		name string
		args args
		want TokenType
	}{
		{
			"Ethereum ERC20",
			args{coinIndex: coin.ETH},
			ERC20,
		},
		{
			"Ethereum Classic ETC20",
			args{coinIndex: coin.ETC},
			ETC20,
		},
		{
			"TomoChain TRC21",
			args{coinIndex: coin.TOMO},
			TRC21,
		},
		{
			"WanChain WAN20",
			args{coinIndex: coin.WAN},
			WAN20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEthereumTokenTypeByIndex(tt.args.coinIndex); got != tt.want {
				t.Errorf("GetEthereumTokenTypeByIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
