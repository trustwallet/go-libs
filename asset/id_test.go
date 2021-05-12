package asset

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseID(t *testing.T) {
	testStruct := []struct {
		givenID     string
		wantedCoin  uint
		wantedToken string
		wantedType  CoinType
		wantedError error
	}{
		{"c714_tTWT-8C2",
			714,
			"TWT-8C2",
			Token,
			nil,
		},
		{"tTWT-8C2_c714",
			714,
			"TWT-8C2",
			Token,
			nil,
		},
		{"c714",
			714,
			"",
			Coin,
			nil,
		},
		{"tTWT-8C2",
			0,
			"",
			Coin,
			errors.New("bad ID"),
		},
		{"c714_TWT-8C2",
			714,
			"",
			Coin,
			nil,
		},
	}

	for _, tt := range testStruct {
		coin, token, err := ParseID(tt.givenID)
		assert.Equal(t, tt.wantedCoin, coin)
		assert.Equal(t, tt.wantedToken, token)
		assert.Equal(t, tt.wantedError, err)
	}
}

func TestBuildID(t *testing.T) {
	testStruct := []struct {
		wantedID   string
		givenCoin  uint
		givenToken string
	}{
		{"c714_tTWT-8C2",
			714,
			"TWT-8C2",
		},
		{"c60",
			60,
			"",
		},
		{"c0",
			0,
			"",
		},
		{"c0_t:fnfjunwpiucU#*0! 02",
			0,
			":fnfjunwpiucU#*0! 02",
		},
	}

	for _, tt := range testStruct {
		id := BuildID(tt.givenCoin, tt.givenToken)
		assert.Equal(t, tt.wantedID, id)
	}
}

func Test_removeFirstChar(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Normal case", "Bob", "ob"},
		{"Empty String", "", ""},
		{"One Char String Test", "A", ""},
		{"Another normaal", "abcdef", "bcdef"},
	}

	for _, tt := range tests {
		var got = removeFirstChar(tt.input)
		if got != tt.expected {
			t.Fatalf("Got %v, Expected %v.", got, tt.expected)
		}
	}
}

func Test_findCoinID(t *testing.T) {
	tests := []struct {
		name        string
		words       []string
		expected    uint
		expectedErr error
	}{
		{"Normal case", []string{"c100", "t60", "e30"}, 100, nil},
		{"Empty coin", []string{"d100", "t60", "e30"}, 0, errors.New("no coin")},
		{"Empty words", []string{}, 0, errors.New("no coin")},
		{"Bad coin", []string{"cd100", "t60", "e30"}, 0, errors.New("bad coin")},
		{"Bad coin #2", []string{"c", "t60", "e30"}, 0, errors.New("bad coin")},
	}

	for _, tt := range tests {
		got, err := FindCoinID(tt.words)
		assert.Equal(t, tt.expected, got)
		assert.Equal(t, tt.expectedErr, err)
	}
}

func Test_findTokenID(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		expected string
	}{
		{"Normal case", []string{"c100", "t60", "e30"}, "60"},
		{"Empty token", []string{"d100", "a", "e30"}, ""},
		{"Empty words", []string{}, ""},
		{"Bad token", []string{"cd100", "t", "e30"}, ""},
	}

	for _, tt := range tests {
		got := FindTokenID(tt.words)
		assert.Equal(t, tt.expected, got)
	}
}
