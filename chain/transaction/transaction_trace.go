package transaction

import (
	"github.com/MetalBlockchain/antelopevm/chain/entity"
	"github.com/MetalBlockchain/antelopevm/chain/fc"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/time"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto"
)

var _ entity.Entity = &TransactionTrace{}

type TransactionTrace struct {
	ID           types.IdType       `serialize:"true" json:"-" eos:"-"`
	Hash         TransactionIdType  `serialize:"true" json:"id"`
	BlockNum     uint64             `serialize:"true" json:"block_num"`
	BlockTime    time.TimePoint     `serialize:"true" json:"block_time"`
	Receipt      TransactionReceipt `serialize:"true" json:"receipt"`
	Elapsed      time.Microseconds  `serialize:"true" json:"elapsed"`
	NetUsage     uint64             `serialize:"true" json:"-"`
	Scheduled    bool               `serialize:"true" json:"scheduled"`
	ActionTraces []ActionTrace      `serialize:"true" json:"action_traces"`
	Except       error              `msg:"-" json:"-"`
	ErrorCode    uint64             `json:"-"`
}

func (a TransactionTrace) GetId() []byte {
	return a.ID.ToBytes()
}

func (a TransactionTrace) GetIndexes() map[string]entity.EntityIndex {
	return map[string]entity.EntityIndex{
		"id": {
			Fields: []string{"ID"},
		},
		"byHash": {
			Fields: []string{"Hash"},
		},
	}
}

func (a TransactionTrace) GetObjectType() uint8 {
	return entity.TransactionObjectType
}

type RamDelta struct {
	Account name.AccountName `serialize:"true" json:"account"`
	Delta   int64            `serialize:"true" json:"delta"`
}

type ActionTrace struct {
	AccountRamDeltas     []RamDelta        `serialize:"true" json:"account_ram_deltas"`
	ActionOrdinal        fc.UnsignedInt    `serialize:"true" json:"action_ordinal"`
	CreatorActionOrdinal fc.UnsignedInt    `serialize:"true" json:"creator_action_ordinal"`
	Receipt              ActionReceipt     `serialize:"true" json:"receipt"`
	Receiver             name.ActionName   `serialize:"true" json:"receiver"`
	Action               Action            `serialize:"true" json:"act"`
	ContextFree          bool              `serialize:"true" json:"context_free"`
	Elapsed              uint64            `serialize:"true" json:"elapsed"`
	Console              string            `serialize:"true" json:"-"`
	TransactionId        TransactionIdType `serialize:"true" json:"trx_id"`
	BlockNum             uint64            `serialize:"true" json:"block_num"`
	BlockTime            time.TimePoint    `serialize:"true" json:"block_time"`
	Except               error             `msg:"-" json:"-"`
	ErrorCode            uint64            `serialize:"true" json:"-"`
}

func NewActionTrace(trace *TransactionTrace, action Action, receiver name.AccountName, contextFree bool, actionOrdinal fc.UnsignedInt, creatorActionOrdinal fc.UnsignedInt) *ActionTrace {
	return &ActionTrace{
		ActionOrdinal:        actionOrdinal,
		CreatorActionOrdinal: creatorActionOrdinal,
		Receiver:             receiver,
		Action:               action,
		ContextFree:          contextFree,
		TransactionId:        trace.Hash,
		BlockNum:             trace.BlockNum,
		BlockTime:            trace.BlockTime,
	}
}

func (a *ActionTrace) Digest() crypto.Sha256 {
	return *crypto.Hash256(a)
}
