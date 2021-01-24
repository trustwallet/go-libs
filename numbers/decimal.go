package numbers

import (
	"errors"
	"math/big"
	"strings"
	"unicode"
)

// DecimalToSatoshis removes the comma in a decimal string
// "12.345" => "12345"
// "0.0230" => "230"
func DecimalToSatoshis(dec string) (string, error) {
	out := strings.TrimLeft(dec, " ")
	out = strings.TrimRight(out, " ")
	out = strings.Replace(out, ".", "", 1)
	// trim left 0's but keep last
	if l := len(out); l >= 2 {
		out = strings.TrimLeft(out[:l-1], "0") + out[l-1:l]
	}
	if len(out) == 0 {
		return "", errors.New("Invalid empty input: " + dec)
	}
	for _, c := range out {
		if !unicode.IsNumber(c) {
			return "", errors.New("not a number: " + dec)
		}
	}
	return out, nil
}

// DecimalExp calculates dec * 10^exp in decimal string representation
func DecimalExp(dec string, exp int) string {
	// 0 * n = 0
	if dec == "0" {
		return "0"
	}
	// Get comma position
	i := strings.IndexRune(dec, '.')
	if i == -1 {
		// Virtual comma at the end of the string
		i = len(dec)
	} else {
		// Remove comma from underlying number
		dec = strings.Replace(dec, ".", "", 1)
	}
	// Shift comma by exponent
	i += exp
	// Remove leading zeros
	origSize := len(dec)
	dec = strings.TrimLeft(dec, "0")
	i -= origSize - len(dec)
	// Fix bounds
	if i <= 0 {
		zeros := ""
		for ; i < 0; i++ {
			zeros += "0"
		}
		return "0." + zeros + dec
	} else if i >= len(dec) {
		for i > len(dec) {
			dec += "0"
		}
		return dec
	}
	// No bound fix needed
	return dec[:i] + "." + dec[i:]
}

// HexToDecimal converts a hexadecimal integer to a base-10 integer
// "0x1fbad5f2e25570000" => "36582000000000000000"
func HexToDecimal(hex string) (string, error) {
	if len(hex) == 0 || hex == "0x" {
		return "0", nil
	}
	var i big.Int
	if _, ok := i.SetString(hex, 0); !ok {
		return "", errors.New("invalid hex: " + hex)
	}
	return i.String(), nil
}

// CutZeroFractional cuts off a decimal separator and zeros to the right.
// Fails if the fractional part contains contains other digits than zeros.
//  - CutZeroFractional("123.00000") => ("123", true)
//  - CutZeroFractional("123.456") => ("", false)
func CutZeroFractional(dec string) (integer string, ok bool) {
	// Get comma position
	comma := strings.IndexRune(dec, '.')
	if comma == -1 {
		return dec, true
	}

	for i := len(dec) - 1; i > comma; i-- {
		if dec[i] != '0' {
			return "", false
		}
	}

	if comma == 0 {
		return "0", true
	} else {
		return dec[:comma], true
	}
}
