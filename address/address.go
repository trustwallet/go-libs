package address

import (
	"encoding/hex"
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

	_, err := hex.DecodeString(string(v))
	if err != nil {
		return "", err
	}

	sha := sha3.NewLegacyKeccak256()
	_, err = sha.Write(v)
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

func ToEIP55ByCoinID(str string, coinID uint) (string, error) {
	switch coinID {
	case coin.ETHEREUM, coin.POA, coin.CLASSIC, coin.TOMOCHAIN, coin.CALLISTO, coin.THUNDERTOKEN, coin.GOCHAIN, coin.WANCHAIN:
		return EIP55Checksum(str)
	default:
		return str, nil
	}
}
