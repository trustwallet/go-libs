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
	TxDelegation        TransactionType = "delegation"
	TxUndelegation      TransactionType = "undelegation"
	TxRedelegation      TransactionType = "redelegation"
	TxAnyAction         TransactionType = "any_action"

	KeyPlaceOrder        KeyType = "place_order"
	KeyCancelOrder       KeyType = "cancel_order"
	KeyIssueToken        KeyType = "issue_token"
	KeyBurnToken         KeyType = "burn_token"
	KeyMintToken         KeyType = "mint_token"
	KeyApproveToken      KeyType = "approve_token"
	KeyStakeDelegate     KeyType = "stake_delegate"
	KeyStakeClaimRewards KeyType = "stake_claim_rewards"

	KeyTitlePlaceOrder   KeyTitle = "Place Order"
	KeyTitleCancelOrder  KeyTitle = "Cancel Order"
	KeyTitleDelegation   KeyTitle = "Delegation"
	KeyTitleUndelegation KeyTitle = "Undelegation"
	KeyTitleClaimRewards KeyTitle = "Claim Rewards"
)

// Transaction fields
type (
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
		Total  int  `json:"total"`
		Docs   []Tx `json:"docs"`
		Status bool `json:"status"`
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
		Date int64 `json:"date"`

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

		Fee Fee `json:"tx_fee"`

		// Metadata data object
		Metadata interface{} `json:"metadata"`
	}

	// Every transaction consumes some Fee
	Fee struct {
		Asset string `json:"asset"`
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
		Value   Amount     `json:"value"`
		Asset   string     `json:"asset"`
		Memo    string     `json:"memo,omitempty"`
		Inputs  []TxOutput `json:"inputs,omitempty"`
		Outputs []TxOutput `json:"outputs,omitempty"`
	}

	// Delegation describes the blocking of a stacked asset
	Delegation struct {
		Asset     string `json:"asset"`
		Value     Amount `json:"value"`
		Validator string `json:"validator"`
		Memo      string `json:"memo,omitempty"`
	}

	// Undelegation describes the releasing of a stacked asset
	Undelegation struct {
		Asset     string `json:"asset"`
		Value     Amount `json:"value"`
		Validator string `json:"validator"`
		Memo      string `json:"memo,omitempty"`
	}

	// In staking there is a possibility to change a validator
	// For that tx of Redelegation type is created
	Redelegation struct {
		Asset         string `json:"asset"`
		Value         Amount `json:"value"`
		FromValidator string `json:"from_validator"`
		ToValidator   string `json:"to_validator"`
		Memo          string `json:"memo,omitempty"`
	}

	// When staking is completed user get rewards which are transferred
	// via ClaimRewards transaction
	ClaimRewards struct {
		Asset string `json:"asset"`
		Value Amount `json:"value"`
		Memo  string `json:"memo,omitempty"`
	}

	// ContractCall describes a
	ContractCall struct {
		Asset string `json:"asset"`
		Input string `json:"input"`
		Value Amount `json:"value"`
	}

	// AnyAction describes all other types
	AnyAction struct {
		Title KeyTitle `json:"title"`
		Key   KeyType  `json:"key"`
		Value Amount   `json:"value"`
		Asset string   `json:"asset"`
		Memo  string   `json:"memo,omitempty"`
	}

	Txs []Tx

	Memo interface {
		Clean()
		GetMemo() string
	}

	Asset interface {
		GetAsset() string
	}
)

var (
	EmptyTxPage = TxPage{Total: 0, Docs: Txs{}, Status: true}
)

func NewTxPage(txs Txs) TxPage {
	if txs == nil {
		txs = Txs{}
	}
	return TxPage{
		Total:  len(txs),
		Docs:   txs,
		Status: true,
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
	for _, tx := range txs {
		memo, ok := tx.Metadata.(Memo)
		if ok {
			memo.Clean()
		}
	}
}

func (txs Txs) SortByDate() Txs {
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Date > txs[j].Date
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

func (t *Transfer) Clean() {
	t.Memo = cleanMemo(t.Memo)
}

func (t *Transfer) GetMemo() string {
	return t.Memo
}

func (t *Transfer) GetAsset() string {
	return t.Asset
}

func (t *Transfer) Addresses() (addresses []string) {
	for _, input := range t.Inputs {
		addresses = append(addresses, input.Address)
	}

	for _, output := range t.Outputs {
		addresses = append(addresses, output.Address)
	}

	return addresses
}

func (d *Delegation) Clean() {
	d.Memo = cleanMemo(d.Memo)
}
func (d *Delegation) GetMemo() string {
	return d.Memo
}
func (t *Delegation) GetAsset() string {
	return t.Asset
}

func (d *Undelegation) Clean() {
	d.Memo = cleanMemo(d.Memo)
}
func (d *Undelegation) GetMemo() string {
	return d.Memo
}
func (t *Undelegation) GetAsset() string {
	return t.Asset
}

func (r *Redelegation) Clean() {
	r.Memo = cleanMemo(r.Memo)
}
func (r *Redelegation) GetMemo() string {
	return r.Memo
}
func (t *Redelegation) GetAsset() string {
	return t.Asset
}

func (cr *ClaimRewards) Clean() {
	cr.Memo = cleanMemo(cr.Memo)
}
func (cr *ClaimRewards) GetMemo() string {
	return cr.Memo
}
func (t *ClaimRewards) GetAsset() string {
	return t.Asset
}

func (сс *ContractCall) GetAsset() string {
	return сс.Asset
}

func (cr *AnyAction) Clean() {
	cr.Memo = cleanMemo(cr.Memo)
}
func (cr *AnyAction) GetMemo() string {
	return cr.Memo
}
func (t *AnyAction) GetAsset() string {
	return t.Asset
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
	switch t.Metadata.(type) {
	case *Transfer, *Delegation, *ContractCall, *AnyAction, *ClaimRewards:
		return append(addresses, t.From, t.To)
	case *Redelegation:
		metadata := t.Metadata.(*Redelegation)
		return append(addresses, metadata.FromValidator, metadata.ToValidator)
	default:
		return addresses
	}
}

func (t *Tx) GetSubscriptionAddresses() ([]string, error) {
	coin, _, err := asset.ParseID(t.Metadata.(Asset).GetAsset())
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

	transfer, ok := t.Metadata.(*Transfer)
	if ok && len(transfer.Inputs) > 0 && len(transfer.Outputs) > 0 {
		addressSet := mapset.NewSet(address)
		return InferDirection(transfer, addressSet)
	}

	return determineTransactionDirection(address, t.From, t.To)
}

func determineTransactionDirection(address, from, to string) Direction {
	if address == to {
		if from == to {
			return DirectionSelf
		}
		return DirectionIncoming
	}
	return DirectionOutgoing
}

func InferDirection(transfer *Transfer, addressSet mapset.Set) Direction {
	inputSet := mapset.NewSet()
	for _, address := range transfer.Inputs {
		inputSet.Add(address.Address)
	}
	outputSet := mapset.NewSet()
	for _, address := range transfer.Outputs {
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
