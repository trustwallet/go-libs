package address

import (
	"strings"

	"github.com/trustwallet/golibs/coin"
	"golang.org/x/crypto/sha3"
)

// Decode decodes a hex string with 0x prefix.
func Remove0x(input string) string {
	if strings.HasPrefix(input, "0x") {
		return input[2:]
	}
	return input
}

// Hex returns an EIP55-compliant hex string representation of the address.
func EIP55Checksum(unchecksummed string) (string, error) {
	v := []byte(Remove0x(strings.ToLower(unchecksummed)))
	sha := sha3.NewLegacyKeccak256()
	_, err := sha.Write(v)
	if err != nil {
		return "", err
	}
	hash := sha.Sum(nil)

	result := v
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	val := string(result)
	return "0x" + val, nil
}

// Returns an EIP55 Wanchain compliant hex string representation of the address.
// See https://wandevs.org/docs/difference-between-wanchain-and-ethereum/
// https://github.com/wanchain/go-wanchain/blob/master/common/types.go#L173
func EIP55ChecksumWanchain(address string) (string, error) {
	v := []byte(Remove0x(strings.ToLower(address)))
	sha := sha3.NewLegacyKeccak256()
	_, err := sha.Write(v)
	if err != nil {
		return "", err
	}
	hash := sha.Sum(nil)

	result := v
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte <= 7 {
			result[i] -= 32
		}
	}
	return "0x" + string(result), nil
}

func ToEIP55ByCoinID(str string, coinID uint) (string, error) {
	switch coinID {
	case coin.ETH, coin.POA, coin.ETC, coin.TOMO, coin.CLO, coin.TT, coin.GO:
		return EIP55Checksum(str)
	case coin.WAN:
		return EIP55ChecksumWanchain(str)
	default:
		return str, nil
	}
}
