package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalling(t *testing.T) {
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
			Metadata:     &Transfer{},
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
