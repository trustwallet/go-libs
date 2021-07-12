package types

import (
	"fmt"
	"sort"
	"strconv"

	mapset "github.com/deckarep/golang-set"
	"github.com/trustwallet/golibs/asset"
)

const (
	StatusCompleted Status = "completed"
	StatusError     Status = "error"

	DirectionOutgoing Direction = "outgoing"
	DirectionIncoming Direction = "incoming"
	DirectionSelf     Direction = "yourself"

	TxTransfer          TransactionType = "transfer"
	TxContractCall      TransactionType = "contract_call"
	TxStakeClaimRewards TransactionType = "stake_claim_rewards"
	TxStakeDelegate     TransactionType = "stake_delegate"
	TxStakeUndelegate   TransactionType = "stake_undelegate"
	TxStakeRedelegate   TransactionType = "stake_redelegate"
)

var SupportedTypes = []TransactionType{
	TxTransfer, TxContractCall, TxStakeClaimRewards, TxStakeDelegate, TxStakeUndelegate, TxStakeRedelegate,
}

// Transaction fields
type (
	Asset           string
	Direction       string
	Status          string
	TransactionType string
	KeyType         string
	KeyTitle        string

	// Amount is a positive decimal integer string.
	// It is written in the smallest possible unit (e.g. Wei, Satoshis)
	Amount string
)

// Data objects
type (
	Block struct {
		Number int64 `json:"number"`
		Txs    []Tx  `json:"txs"`
	}

	TxPage struct {
		Total int  `json:"total"`
		Docs  []Tx `json:"docs"`
	}

	// Tx describes an on-chain transaction generically
	Tx struct {
		// Unique identifier
		ID string `json:"id"`

		// Address of the transaction sender
		From string `json:"from"`

		// Address of the transaction recipient
		To string `json:"to"`

		// Unix timestamp of the block the transaction was included in
		BlockCreatedAt int64 `json:"block_created_at"`

		// Height of the block the transaction was included in
		Block uint64 `json:"block"`

		// Status of the transaction e.g: "completed", "pending", "error"
		Status Status `json:"status"`

		// Empty if the transaction "completed" or "pending", else error explaining why the transaction failed (optional)
		Error string `json:"error,omitempty"`

		// Transaction nonce or sequence
		Sequence uint64 `json:"sequence"`

		// Type of metadata
		Type TransactionType `json:"type"`

		// Transaction Direction
		Direction Direction `json:"direction,omitempty"`

		Inputs  []TxOutput `json:"inputs,omitempty"`
		Outputs []TxOutput `json:"outputs,omitempty"`

		Memo string `json:"memo,omitempty"`

		Fee Fee `json:"fee"`

		// Metadata data object
		Metadata interface{} `json:"metadata"`

		// Create At indicates transactions creation time in database, Unix
		CreatedAt int64 `json:"created_at"`
	}

	// Every transaction consumes some Fee
	Fee struct {
		Asset Asset  `json:"asset"`
		Value Amount `json:"value"`
	}

	// UTXO transactions consist of a set of inputs and a set of outputs
	// both represented by TxOutput model
	TxOutput struct {
		Address string `json:"address"`
		Value   Amount `json:"value"`
	}

	// Transfer describes the transfer of currency
	Transfer struct {
		Asset Asset  `json:"asset"`
		Value Amount `json:"value"`
	}

	// ContractCall describes a
	ContractCall struct {
		Asset Asset  `json:"asset"`
		Value Amount `json:"value"`
		Input string `json:"input"`
	}

	Txs []Tx

	AssetHolder interface {
		GetAsset() Asset
	}

	Validator interface {
		Validate() error
	}
)

var (
	EmptyTxPage = TxPage{Total: 0, Docs: Txs{}}
)

func NewTxPage(txs Txs) TxPage {
	if txs == nil {
		txs = Txs{}
	}
	return TxPage{
		Total: len(txs),
		Docs:  txs,
	}
}

func (txs Txs) FilterUniqueID() Txs {
	keys := make(map[string]bool)
	list := make(Txs, 0)
	for _, entry := range txs {
		if _, value := keys[entry.ID]; !value {
			keys[entry.ID] = true
			list = append(list, entry)
		}
	}
	return list
}

func (txs Txs) CleanMemos() {
	for i := range txs {
		txs[i].Memo = cleanMemo(txs[i].Memo)
	}
}

func (txs Txs) SortByBlockCreationTime() Txs {
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].BlockCreatedAt > txs[j].BlockCreatedAt
	})
	return txs
}

func (txs Txs) FilterTransactionsByType(types []TransactionType) Txs {
	result := make(Txs, 0)
	for _, tx := range txs {
		for _, t := range types {
			if tx.Type == t {
				result = append(result, tx)
			}
		}
	}

	return result
}

func (t *Transfer) GetAsset() Asset {
	return t.Asset
}

func (t *Transfer) Validate() error {
	if t.Value == "" {
		return fmt.Errorf("emtpy transfer value")
	}

	if t.Asset == "" {
		return fmt.Errorf("empty transfer asset")
	}

	return nil
}

func (cc *ContractCall) GetAsset() Asset {
	return cc.Asset
}

func (cc *ContractCall) Validate() error {
	if cc.Value == "" {
		return fmt.Errorf("empty contract call value")
	}

	if cc.Asset == "" {
		return fmt.Errorf("empty contract call asset")
	}

	return nil
}

func cleanMemo(memo string) string {
	if len(memo) == 0 {
		return ""
	}

	_, err := strconv.ParseFloat(memo, 64)
	if err != nil {
		return ""
	}

	return memo
}

func (t *Tx) GetAddresses() []string {
	addresses := make([]string, 0)
	switch t.Type {
	case TxTransfer:
		if len(t.Inputs) > 0 || len(t.Outputs) > 0 {
			uniqueAddresses := make(map[string]struct{})
			for _, input := range t.Inputs {
				uniqueAddresses[input.Address] = struct{}{}
			}

			for _, output := range t.Outputs {
				uniqueAddresses[output.Address] = struct{}{}
			}

			for address := range uniqueAddresses {
				addresses = append(addresses, address)
			}

			return addresses
		}

		return append(addresses, t.From, t.To)
	case TxContractCall:
		return append(addresses, t.From, t.To)
	case TxStakeDelegate, TxStakeRedelegate, TxStakeUndelegate, TxStakeClaimRewards:
		return append(addresses, t.From)
	default:
		return addresses
	}
}

func (t *Tx) GetSubscriptionAddresses() ([]string, error) {
	coin, _, err := asset.ParseID(string(t.Metadata.(AssetHolder).GetAsset()))
	if err != nil {
		return nil, err
	}

	addresses := t.GetAddresses()
	result := make([]string, len(addresses))
	for i, a := range addresses {
		result[i] = fmt.Sprintf("%d_%s", coin, a)
	}

	return result, nil
}

func (t *Tx) GetDirection(address string) Direction {
	if len(t.Direction) > 0 {
		return t.Direction
	}

	if len(t.Inputs) > 0 && len(t.Outputs) > 0 {
		addressSet := mapset.NewSet(address)
		return InferDirection(t, addressSet)
	}

	return t.determineTransactionDirection(address, t.From, t.To)
}

func (t *Tx) determineTransactionDirection(address, from, to string) Direction {
	if t.Type == TxStakeUndelegate || t.Type == TxStakeClaimRewards {
		return DirectionIncoming
	}

	if address == to {
		if from == to {
			return DirectionSelf
		}
		return DirectionIncoming
	}
	return DirectionOutgoing
}

func (t *Tx) IsUTXO() bool {
	return t.Type == TxTransfer && len(t.Outputs) > 0
}

func (t *Tx) GetUTXOValueFor(address string) (Amount, error) {
	isTransferOut := false
	isSelf := true

	var totalInputValue uint64
	var addressInputValue uint64
	for _, input := range t.Inputs {
		value, err := strconv.ParseUint(string(input.Value), 10, 64)
		if err != nil {
			return "0", fmt.Errorf("input value for address %s: %v", input.Address, err)
		}

		totalInputValue += value

		if input.Address == address {
			addressInputValue = value
			isTransferOut = true
		} else {
			isSelf = false
		}
	}

	var addressOutputValue uint64
	var totalOutputValue uint64
	for _, output := range t.Outputs {
		value, err := strconv.ParseUint(string(output.Value), 10, 64)
		if err != nil {
			return "0", fmt.Errorf("output value for address %s: %v", output.Address, err)
		}
		totalOutputValue += value
		if output.Address == address {
			addressOutputValue += value
		} else {
			isSelf = false
		}
	}

	var result uint64
	if isTransferOut && !isSelf {
		result = addressInputValue - (totalInputValue-totalOutputValue)/uint64(len(t.Inputs)) - addressOutputValue
	} else {
		result = addressOutputValue
	}

	return Amount(fmt.Sprintf("%d", result)), nil
}

func InferDirection(tx *Tx, addressSet mapset.Set) Direction {
	inputSet := mapset.NewSet()
	for _, address := range tx.Inputs {
		inputSet.Add(address.Address)
	}
	outputSet := mapset.NewSet()
	for _, address := range tx.Outputs {
		outputSet.Add(address.Address)
	}
	intersect := addressSet.Intersect(inputSet)
	if intersect.Cardinality() == 0 {
		return DirectionIncoming
	}
	if outputSet.IsProperSubset(addressSet) || outputSet.Equal(inputSet) {
		return DirectionSelf
	}
	return DirectionOutgoing
}

func IsTxTypeAmong(txType TransactionType, types []TransactionType) bool {
	result := false
	for _, t := range types {
		if txType == t {
			result = true
			break
		}
	}

	return result
}
