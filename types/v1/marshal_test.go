package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTxMarshalling(t *testing.T) {
	tests := []struct {
		Name         string
		Type         TransactionType
		Metadata     interface{}
		marshalErr   assert.ErrorAssertionFunc
		unmarshalErr assert.ErrorAssertionFunc
	}{
		{
			Name:         "transfer",
			Type:         TxTransfer,
			Metadata:     &Transfer{},
			marshalErr:   assert.NoError,
			unmarshalErr: assert.NoError,
		},
		{
			Name:         "contract_call",
			Type:         TxContractCall,
			Metadata:     &ContractCall{},
			marshalErr:   assert.NoError,
			unmarshalErr: assert.NoError,
		},
		{
			Name:         "claim_rewards",
			Type:         TxStakeClaimRewards,
			Metadata:     &Transfer{},
			marshalErr:   assert.NoError,
			unmarshalErr: assert.NoError,
		},
		{
			Name:         "delegate",
			Type:         TxStakeDelegate,
			Metadata:     &Transfer{},
			marshalErr:   assert.NoError,
			unmarshalErr: assert.NoError,
		},
		{
			Name:         "undelegate",
			Type:         TxStakeUndelegate,
			Metadata:     &Transfer{},
			marshalErr:   assert.NoError,
			unmarshalErr: assert.NoError,
		},
		{
			Name:         "redelegate",
			Type:         TxStakeRedelegate,
			Metadata:     &Transfer{},
			marshalErr:   assert.NoError,
			unmarshalErr: assert.NoError,
		},
		{
			Name:         "redelegate",
			Type:         TxStakeRedelegate,
			Metadata:     &Transfer{},
			marshalErr:   assert.NoError,
			unmarshalErr: assert.NoError,
		},
		{
			Name:         "without_type",
			Metadata:     &Transfer{},
			marshalErr:   assert.Error,
			unmarshalErr: assert.Error,
		},
		{
			Name:         "unsupported_type",
			Metadata:     &Transfer{},
			marshalErr:   assert.Error,
			unmarshalErr: assert.Error,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			tx := Tx{
				Type:     tc.Type,
				Metadata: tc.Metadata,
			}

			data, err := json.Marshal(tx)
			tc.marshalErr(t, err)

			var receiver Tx
			err = json.Unmarshal(data, &receiver)
			tc.unmarshalErr(t, err)
		})
	}
}

func TestTxsMarshalling(t *testing.T) {
	tests := []struct {
		Name         string
		Type         TransactionType
		Metadata     interface{}
		marshalErr   assert.ErrorAssertionFunc
		unmarshalErr assert.ErrorAssertionFunc
		expectNil    bool
	}{
		{
			Name:         "transfer",
			Type:         TxTransfer,
			Metadata:     &Transfer{},
			marshalErr:   assert.NoError,
			unmarshalErr: assert.NoError,
		},
		{
			Name:         "contract_call",
			Type:         TxContractCall,
			Metadata:     &ContractCall{},
			marshalErr:   assert.NoError,
			unmarshalErr: assert.NoError,
		},
		{
			Name:         "claim_rewards",
			Type:         TxStakeClaimRewards,
			Metadata:     &Transfer{},
			marshalErr:   assert.NoError,
			unmarshalErr: assert.NoError,
		},
		{
			Name:         "delegate",
			Type:         TxStakeDelegate,
			Metadata:     &Transfer{},
			marshalErr:   assert.NoError,
			unmarshalErr: assert.NoError,
		},
		{
			Name:         "undelegate",
			Type:         TxStakeUndelegate,
			Metadata:     &Transfer{},
			marshalErr:   assert.NoError,
			unmarshalErr: assert.NoError,
		},
		{
			Name:         "redelegate",
			Type:         TxStakeRedelegate,
			Metadata:     &Transfer{},
			marshalErr:   assert.NoError,
			unmarshalErr: assert.NoError,
		},
		{
			Name:         "redelegate",
			Type:         TxStakeRedelegate,
			Metadata:     &Transfer{},
			marshalErr:   assert.NoError,
			unmarshalErr: assert.NoError,
		},
		{
			Name:         "without_type",
			Metadata:     &Transfer{},
			marshalErr:   assert.Error,
			unmarshalErr: assert.Error,
			expectNil:    true,
		},
		{
			Name:         "unsupported_type",
			Metadata:     &Transfer{},
			marshalErr:   assert.Error,
			unmarshalErr: assert.Error,
			expectNil:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			txs := Txs{{
				Type:     tc.Type,
				Metadata: tc.Metadata,
				Status:   StatusCompleted,
			}}

			data, err := json.Marshal(txs)
			tc.marshalErr(t, err)

			var receiver Txs
			err = json.Unmarshal(data, &receiver)
			tc.unmarshalErr(t, err)

			if tc.expectNil {
				assert.Equal(t, Txs(nil), receiver)
			} else {
				assert.Equal(t, txs[0], receiver[0])
			}
		})
	}
}
