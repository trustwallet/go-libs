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

			assert.Equal(t, tc.expectedMemo, txs[0].Memo)
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
				Type:     TxTransfer,
				From:     "from",
				To:       "to",
				Metadata: &Transfer{},
			},
			expected: []string{"from", "to"},
		},
		{
			name: "delegation",
			tx: Tx{
				Type:     TxStakeDelegate,
				From:     "from",
				To:       "to",
				Metadata: &Transfer{},
			},
			expected: []string{"from"},
		},
		{
			name: "undelegation",
			tx: Tx{
				Type:     TxStakeUndelegate,
				From:     "from",
				To:       "to",
				Metadata: &Transfer{},
			},
			expected: []string{"from"},
		},
		{
			name: "claim_rewards",
			tx: Tx{
				Type:     TxStakeClaimRewards,
				From:     "from",
				To:       "to",
				Metadata: &Transfer{},
			},
			expected: []string{"from"},
		},
		{
			name: "contract_call",
			tx: Tx{
				Type:     TxContractCall,
				From:     "from",
				To:       "to",
				Metadata: &ContractCall{},
			},
			expected: []string{"from", "to"},
		},
		{
			name: "any_action",
			tx: Tx{
				Type:     TxTransfer,
				From:     "from",
				To:       "to",
				Metadata: &Transfer{},
			},
			expected: []string{"from", "to"},
		},
		{
			name: "redelegation",
			tx: Tx{
				Type:     TxStakeRedelegate,
				From:     "from",
				To:       "to_validator",
				Metadata: &Transfer{},
			},
			expected: []string{"from"},
		},
		{
			name: "undefined",
			tx: Tx{
				From: "from",
				To:   "to",
			},
			expected: []string{},
		},
		{
			name: "utxo",
			tx: Tx{
				Type:     TxTransfer,
				From:     "from_utxo",
				To:       "from_utxo",
				Inputs:   []TxOutput{{Address: "from_utxo"}},
				Outputs:  []TxOutput{{Address: "from_utxo"}, {Address: "to_utxo"}},
				Metadata: &Transfer{},
			},
			expected: []string{"from_utxo", "to_utxo"},
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
		{
			name: "stake_delegate",
			tx: Tx{
				From: "sender",
				To:   "sender",
			},
			address:  "sender",
			expected: DirectionSelf,
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
		{
			name: "stake_undelegate",
			tx: Tx{
				From: "delegator",
				To:   "owner",
				Type: TxStakeUndelegate,
			},
			address:  "owner",
			expected: DirectionIncoming,
		},
		{
			name: "stake_redelegate",
			tx: Tx{
				From: "delegator1",
				To:   "delegator2",
				Type: TxStakeRedelegate,
			},
			address:  "owner",
			expected: DirectionOutgoing,
		},
		{
			name: "stake_delegate",
			tx: Tx{
				From: "owner",
				To:   "delegator",
				Type: TxStakeDelegate,
			},
			address:  "owner",
			expected: DirectionOutgoing,
		},
		{
			name: "stake_claim_rewards",
			tx: Tx{
				From: "delegator",
				To:   "sender",
				Type: TxStakeClaimRewards,
			},
			address:  "sender",
			expected: DirectionIncoming,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.tx.GetDirection(tc.address)
			assert.Equal(t, tc.expected, result)
		})
	}

}

func TestUTXOValueByAddress(t *testing.T) {
	tests := []struct {
		name                 string
		tx                   Tx
		address              string
		expected             Amount
		expectedErrAssertion assert.ErrorAssertionFunc
	}{
		{
			name: "transfer_self",
			tx: Tx{
				Inputs: []TxOutput{{
					Address: "addr",
					Value:   "1000",
				}},
				Outputs: []TxOutput{
					{
						Address: "addr",
						Value:   "900",
					},
					{
						Address: "addr",
						Value:   "100",
					},
				},
			},
			address:              "addr",
			expected:             "1000",
			expectedErrAssertion: assert.NoError,
		},
		{
			name: "transfer_in",
			tx: Tx{
				Outputs: []TxOutput{{
					Address: "addr",
					Value:   "1000",
				}},
			},
			address:              "addr",
			expected:             "1000",
			expectedErrAssertion: assert.NoError,
		},
		{
			name: "transfer_out_with_utxo",
			tx: Tx{
				Inputs: []TxOutput{{
					Address: "addr",
					Value:   "1000",
				}},
				Outputs: []TxOutput{
					{
						Address: "addr",
						Value:   "100",
					},
					{
						Address: "addr1",
						Value:   "800",
					},
				},
			},
			address:              "addr",
			expected:             "800",
			expectedErrAssertion: assert.NoError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.tx.GetUTXOValueFor(tc.address)
			tc.expectedErrAssertion(t, err)

			assert.Equal(t, tc.expected, result)
		})
	}
}
