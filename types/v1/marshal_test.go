package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalling(t *testing.T) {
	tests := []struct {
		Name     string
		Type     TransactionType
		Metadata interface{}
	}{
		{
			Name:     "transfer",
			Type:     TxTransfer,
			Metadata: &Transfer{},
		},
		{
			Name:     "contract_call",
			Type:     TxContractCall,
			Metadata: &Transfer{},
		},
		{
			Name:     "claim_rewards",
			Type:     TxStakeClaimRewards,
			Metadata: &Transfer{},
		},
		{
			Name:     "delegate",
			Type:     TxStakeDelegate,
			Metadata: &Transfer{},
		},
		{
			Name:     "undelegate",
			Type:     TxStakeUndelegate,
			Metadata: &Transfer{},
		},
		{
			Name:     "redelegate",
			Type:     TxStakeRedelegate,
			Metadata: &Transfer{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			tx := Tx{
				Type:     tc.Type,
				Metadata: tc.Metadata,
			}

			data, err := json.Marshal(tx)
			assert.NoError(t, err)

			var receiver Tx
			err = json.Unmarshal(data, &receiver)
			assert.NoError(t, err)
		})
	}
}
