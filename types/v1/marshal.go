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
	case TxTransfer:
		t.Metadata = new(Transfer)
	case TxContractCall:
		t.Metadata = new(ContractCall)
	case TxAnyAction:
		t.Metadata = new(AnyAction)
	case TxDelegation:
		t.Metadata = new(Delegation)
	case TxUndelegation:
		t.Metadata = new(Undelegation)
	case TxRedelegation:
		t.Metadata = new(Redelegation)
	case TxStakeClaimRewards:
		t.Metadata = new(ClaimRewards)
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
// Sets the Type field to the correct value based on the Metadata type.
func (t *Tx) MarshalJSON() ([]byte, error) {
	// Set type from metadata content
	switch t.Metadata.(type) {
	case *Transfer:
		t.Type = TxTransfer
	case *ContractCall:
		t.Type = TxContractCall
	case *AnyAction:
		t.Type = TxAnyAction
	case *Delegation:
		t.Type = TxDelegation
	case *Undelegation:
		t.Type = TxUndelegation
	case *Redelegation:
		t.Type = TxRedelegation
	case *ClaimRewards:
		t.Type = TxStakeClaimRewards
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
