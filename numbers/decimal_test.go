package numbers

import "testing"

func TestDecimalToSatoshis(t *testing.T) {
	assertSatEquals := func(expected string, input string) {
		actual, err := DecimalToSatoshis(input)
		if err != nil {
			t.Error(err)
		}
		if expected != actual {
			t.Errorf("expected %s, got %s, input %s", expected, actual, input)
		}
	}

	assertSatError := func(input string) {
		actual, err := DecimalToSatoshis(input)
		if err == nil {
			t.Errorf("Expected error but no error: got %s, input %s", actual, input)
		}
	}

	assertSatEquals("10", "1.0")
	assertSatEquals("1", "0.1")
	assertSatEquals("13602", "136.02")
	assertSatEquals("13602", "0136.02")
	assertSatEquals("1500000", "0.01500000")
	assertSatEquals("0", "0")
	assertSatEquals("2030", "0.002030")
	assertSatEquals("101010", "0101010")
	assertSatEquals("11001100", "0011001100")
	assertSatEquals("376", " 376")
	assertSatEquals("376", "376 ")

	assertSatError("12NotNumber34")
	assertSatError("12,34")
	assertSatError("")
	assertSatError(" ")
	assertSatError("37 6")
	assertSatError("37,6")
}

func TestDecimalExp(t *testing.T) {
	assertEquals := func(inputDec string, inputExp int, expected string) {
		actual := DecimalExp(inputDec, inputExp)
		if expected != actual {
			t.Errorf("expected: %s * (10^%d) = %s, got %s",
				inputDec, inputExp, expected, actual)
		}
	}

	// No-Op
	assertEquals("0", 300, "0")
	assertEquals("0", 8, "0")
	assertEquals("123", 0, "123")
	assertEquals("0.456", 0, "0.456")
	assertEquals("123.456", 0, "123.456")

	// In-Bounds, comma
	assertEquals("12.34", -1, "1.234")
	assertEquals("12.34", 1, "123.4")

	// 1 past bounds, comma
	assertEquals("12.34", -2, "0.1234")
	assertEquals("12.34", 2, "1234")

	// n past bounds, comma
	assertEquals("12.34", -4, "0.001234")
	assertEquals("12.34", 4, "123400")

	// Integer
	assertEquals("1234", -1, "123.4")
	assertEquals("1234", 1, "12340")

	// Denormalized
	assertEquals("0.1234", -1, "0.01234")
	assertEquals("0.1234", 1, "1.234")

	// Tiny
	assertEquals("0.001234", -1, "0.0001234")
	assertEquals("0.001234", 1, "0.01234")
	assertEquals("0.000375", 8, "37500")
}

func TestCutZeroFractional(t *testing.T) {
	assertEquals := func(inputDec string, expected string, expOk bool) {
		actual, ok := CutZeroFractional(inputDec)
		if expected != actual || ok != expOk {
			t.Errorf("expected: %s => (%s, %v), actual: (%s, %v)",
				inputDec, expected, expOk, actual, ok)
		}
	}

	// No comma
	assertEquals("", "", true)
	assertEquals("eee", "eee", true)

	// Length 1
	assertEquals(".", "0", true)
	assertEquals(".3", "", false)
	assertEquals(".0", "0", true)
	assertEquals("0.", "0", true)
	assertEquals("1.0", "1", true)
	assertEquals("1.1", "", false)
	assertEquals("1.0.0", "", false)

	// Arbitrary content left to comma
	assertEquals("eee.000", "eee", true)
	assertEquals("eee.001", "", false)
	assertEquals("eee.100", "", false)

	// Long strings
	assertEquals("163056848705309039018274728757999527956626319283048085297785610.238523", "", false)
	assertEquals("11434397695550368380599182733571088333799363173941798154.0000000000000", "11434397695550368380599182733571088333799363173941798154", true)
}

func TestHexToDecimal(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		want    string
		wantErr bool
	}{
		{
			name:    "Empty value",
			hex:     "",
			want:    "0",
			wantErr: false,
		},
		{
			name:    "Empty value 2",
			hex:     "0x",
			want:    "0",
			wantErr: false,
		},
		{
			name:    "Empty value 3",
			hex:     "0x0",
			want:    "0",
			wantErr: false,
		},
		{
			name:    "Hex to decimal",
			hex:     "0x1fbad5f2e25570000",
			want:    "36582000000000000000",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HexToDecimal(tt.hex)
			if (err != nil) != tt.wantErr {
				t.Errorf("HexToDecimal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HexToDecimal() got = %v, want %v", got, tt.want)
			}
		})
	}
}
