package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTxs_CleanMemos(t *testing.T) {
	tests := []struct {
		name         string
		tx           Tx
		expectedMemo string
	}{
		{
			name:         "transfer_ok",
			tx:           Tx{Memo: "1"},
			expectedMemo: "1",
		},
		{
			name:         "transfer_empty_string",
			tx:           Tx{Metadata: &Transfer{}},
			expectedMemo: "",
		},
		{
			name:         "transfer_non_number",
			tx:           Tx{Memo: "non_number"},
			expectedMemo: "",
		},
		{
			name:         "delegation_ok",
			tx:           Tx{Memo: "1"},
			expectedMemo: "1",
		},
		{
			name:         "delegation_empty_string",
			tx:           Tx{Metadata: &Transfer{}},
			expectedMemo: "",
		},
		{
			name:         "delegation_non_number",
			tx:           Tx{Memo: "non_number"},
			expectedMemo: "",
		},
		{
			name:         "redelegation_ok",
			tx:           Tx{Memo: "1"},
			expectedMemo: "1",
		},
		{
			name:         "redelegation_empty_string",
			tx:           Tx{Metadata: &Transfer{}},
			expectedMemo: "",
		},
		{
			name:         "redelegation_non_number",
			tx:           Tx{Memo: "non_number"},
			expectedMemo: "",
		},
		{
			name:         "claim_rewards_ok",
			tx:           Tx{Memo: "1"},
			expectedMemo: "1",
		},
		{
			name:         "claim_rewards_empty_string",
			tx:           Tx{Metadata: &Transfer{}},
			expectedMemo: "",
		},
		{
			name:         "claim_rewards_non_number",
			tx:           Tx{Memo: "non_number"},
			expectedMemo: "",
		},
		{
			name:         "any_action_ok",
			tx:           Tx{Memo: "1"},
			expectedMemo: "1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			txs := Txs{tc.tx}
			txs.CleanMemos()

			memo, ok := txs[0].Metadata.(Memo)
			if ok {
				assert.Equal(t, tc.expectedMemo, memo.GetMemo())
			}
		})
	}
}

func TestCleanMemo(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "empty_value",
			value:    "",
			expected: "",
		},
		{
			name:     "string_value",
			value:    "test",
			expected: "",
		},
		{
			name:     "good_number_value",
			value:    "1",
			expected: "1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := cleanMemo(tc.value)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestTx_GetAddresses(t *testing.T) {
	tests := []struct {
		name     string
		tx       Tx
		expected []string
	}{
		{
			name: "transfer",
			tx: Tx{
				From:     "from",
				To:       "to",
				Metadata: &Transfer{},
			},
			expected: []string{"from", "to"},
		},
		{
			name: "delegation",
			tx: Tx{
				From:     "from",
				To:       "to",
				Metadata: &Transfer{},
			},
			expected: []string{"from", "to"},
		},
		{
			name: "contract_call",
			tx: Tx{
				From:     "from",
				To:       "to",
				Metadata: &ContractCall{},
			},
			expected: []string{"from", "to"},
		},
		{
			name: "any_action",
			tx: Tx{
				From:     "from",
				To:       "to",
				Metadata: &Transfer{},
			},
			expected: []string{"from", "to"},
		},
		{
			name: "redelegation",
			tx: Tx{
				From:     "from_validator",
				To:       "to_validator",
				Metadata: &Transfer{},
			},
			expected: []string{"from_validator", "to_validator"},
		},
		{
			name: "undefined",
			tx: Tx{
				From: "from",
				To:   "to",
			},
			expected: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.tx.GetAddresses()
			assert.Equal(t, tc.expected, result)
		})
	}

}

func TestTx_GetDirection(t *testing.T) {
	tests := []struct {
		name     string
		tx       Tx
		address  string
		expected Direction
	}{
		{
			name: "direction_defined_outgoing",
			tx: Tx{
				Direction: DirectionOutgoing,
			},
			expected: DirectionOutgoing,
		},
		{
			name: "direction_defined_incoming",
			tx: Tx{
				Direction: DirectionIncoming,
			},
			expected: DirectionIncoming,
		},
		{
			name: "utxo_outgoing",
			tx: Tx{
				Inputs: []TxOutput{
					{
						Address: "sender",
					},
				},
				Outputs: []TxOutput{
					{
						Address: "receiver",
					},
				},
			},
			address:  "sender",
			expected: DirectionOutgoing,
		},
		{
			name: "utxo_incoming",
			tx: Tx{
				Inputs: []TxOutput{
					{
						Address: "sender",
					},
				},
				Outputs: []TxOutput{
					{
						Address: "receiver",
					},
				},
			},
			address:  "receiver",
			expected: DirectionIncoming,
		},
		{
			name: "utxo_self",
			tx: Tx{
				Inputs: []TxOutput{
					{
						Address: "sender",
					},
				},
				Outputs: []TxOutput{
					{
						Address: "sender",
					},
				},
			},
			address:  "sender",
			expected: DirectionSelf,
		},
		{
			name: "outgoing",
			tx: Tx{
				From: "sender",
				To:   "receiver",
			},
			address:  "sender",
			expected: DirectionOutgoing,
		},
		{
			name: "incoming",
			tx: Tx{
				From: "sender",
				To:   "receiver",
			},
			address:  "receiver",
			expected: DirectionIncoming,
		},
		{
			name: "self",
			tx: Tx{
				From: "sender",
				To:   "sender",
			},
			address:  "sender",
			expected: DirectionSelf,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.tx.GetDirection(tc.address)
			assert.Equal(t, tc.expected, result)
		})
	}

}
