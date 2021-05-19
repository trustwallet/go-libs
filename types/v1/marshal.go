package types

import (
	"encoding/json"
	"errors"
)

// Tx, but with default JSON marshalling methods
type wrappedTx Tx

// UnmarshalJSON creates a transaction along with metadata from a JSON object.
// Fails if the meta object can't be read.
func (t *Tx) UnmarshalJSON(data []byte) error {
	// Wrap the Tx type to avoid infinite recursion
	var wrapped wrappedTx

	var raw json.RawMessage
	wrapped.Metadata = &raw
	if err := json.Unmarshal(data, &wrapped); err != nil {
		return err
	}

	*t = Tx(wrapped)

	switch t.Type {
	case TxTransfer, TxStakeDelegate, TxStakeUndelegate, TxStakeRedelegate, TxStakeClaimRewards:
		t.Metadata = new(Transfer)
	case TxContractCall:
		t.Metadata = new(ContractCall)
	default:
		return errors.New("unsupported tx type")
	}

	err := json.Unmarshal(raw, t.Metadata)
	if err != nil {
		return err
	}
	return nil
}

// MarshalJSON creates a JSON object from a transaction.
func (t *Tx) MarshalJSON() ([]byte, error) {
	// validate metadata type
	switch t.Metadata.(type) {
	case *Transfer, *ContractCall:
		break
	default:
		return nil, errors.New("unsupported tx metadata")
	}

	// Set status to completed by default
	if t.Status == "" {
		t.Status = StatusCompleted
	}

	// Wrap the Tx type to avoid infinite recursion
	return json.Marshal(wrappedTx(*t))
}

// Sort sorts the response by date, descending
func (txs Txs) Len() int           { return len(txs) }
func (txs Txs) Less(i, j int) bool { return txs[i].Date > txs[j].Date }
func (txs Txs) Swap(i, j int)      { txs[i], txs[j] = txs[j], txs[i] }
