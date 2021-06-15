package address

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/golibs/coin"
)

func TestEIP55Checksum(t *testing.T) {
	type args struct {
		unchecksummed string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test checksum number", args{"16345785d8a00000"}, "0x16345785D8A00000", false},
		{"test checksum hex", args{"fffdefefed"}, "0xFfFDEfeFeD", false},
		{"test checksum 3", args{"0x0000000000000000003731342d4f4e452d354639"}, "0x0000000000000000003731342d4f4E452d354639", false},
		{"test checksum 4", args{"0000000000000000003731342d4f4e452d354639"}, "0x0000000000000000003731342d4f4E452d354639", false},
		{"test checksum Ethereum address", args{"0x84a0d77c693adabe0ebc48f88b3ffff010577051"}, "0x84A0d77c693aDAbE0ebc48F88b3fFFF010577051", false},
		{"test invalid address format", args{"https://bscscan.com/token/0x959229D94c9060552daea25AC17193bcA65D7884"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EIP55Checksum(tt.args.unchecksummed)
			if (err != nil) != tt.wantErr {
				t.Errorf("EIP55Checksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EIP55Checksum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemove0x(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"remove 0x from addres", "0x158079ee67fce2f58472a96584a73c7ab9ac95c1", "158079ee67fce2f58472a96584a73c7ab9ac95c1"},
		{"remove 0x from hash", "0x230798fe22abff459b004675bf827a4089326a296fa4165d0c2ad27688e03e0c", "230798fe22abff459b004675bf827a4089326a296fa4165d0c2ad27688e03e0c"},
		{"remove 0x hex value", "0xfffdefefed", "fffdefefed"},
		{"remove 0x hex number", "0x16345785d8a0000", "16345785d8a0000"},
		{"remove hex without 0x", "trustwallet", "trustwallet"},
		{"remove hex number without 0x", "16345785d8a0000", "16345785d8a0000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Remove0x(tt.input); got != tt.want {
				t.Errorf("Remove0x() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToEIP55ByCoinID(t *testing.T) {
	var (
		addr1                        = "0xea674fdde714fd979de3edf0f56aa9716b898ec8"
		addr1EIP55                   = "0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8"
		wanAddrLowercase             = "0xae96137e0e05681ed2f5d1af272c3ee512939d0f"
		wanAddrEIP55Checksum         = "0xAe96137E0e05681eD2F5D1AF272C3ee512939D0F"
		wanAddrEIP55ChecksumWanchain = "0xAe96137E0e05681eD2F5D1AF272C3ee512939D0F"
		tests                        = []struct {
			name, address, expectedAddress string
			coinID                         uint
		}{
			{"Ethereum", addr1, addr1EIP55, coin.ETHEREUM},
			{"Ethereum Classic", addr1, addr1EIP55, coin.CLASSIC},
			{"POA", addr1, addr1EIP55, coin.POA},
			{"Callisto", addr1, addr1EIP55, coin.CALLISTO},
			{"Tomochain", addr1, addr1EIP55, coin.TOMOCHAIN},
			{"Thunder", addr1, addr1EIP55, coin.THUNDERTOKEN},
			{"Thunder", addr1, addr1EIP55, coin.THUNDERTOKEN},
			{"GoChain", addr1, addr1EIP55, coin.GOCHAIN},
			{"Wanchain 1", wanAddrLowercase, wanAddrEIP55ChecksumWanchain, coin.WANCHAIN},
			{"Wanchain 2", wanAddrEIP55Checksum, wanAddrEIP55ChecksumWanchain, coin.WANCHAIN},
			{"Non Ethereum like chain 1", "", "", coin.TRON},
			{"Non Ethereum like chain 2", addr1, addr1, coin.BINANCE},
		}
	)

	t.Run("Test TestToEIP55ByCoinID", func(t *testing.T) {
		for _, tt := range tests {
			actual, _ := ToEIP55ByCoinID(tt.address, tt.coinID)
			assert.Equal(t, tt.expectedAddress, actual)
		}
	})
}
