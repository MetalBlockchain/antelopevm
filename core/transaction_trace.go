package core

import (
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
)

//go:generate msgp
type TransactionReceiptHeader struct {
	Status        TransactionStatus `serialize:"true" json:"status"`
	CpuUsageUs    uint32            `serialize:"true" json:"cpu_usage_us"`
	NetUsageWords Vuint32           `serialize:"true" json:"net_usage_words"`
}

type TransactionReceipt struct {
	TransactionReceiptHeader `serialize:"true"`
	Transaction              PackedTransaction `serialize:"true" json:"trx" eos:"-"`
}

func (t *TransactionReceipt) Digest() crypto.Sha256 {
	enc := crypto.NewSha256()
	status, _ := rlp.EncodeToBytes(t.Status)
	cpuUsageUs, _ := rlp.EncodeToBytes(t.CpuUsageUs)
	netUsageWords, _ := rlp.EncodeToBytes(t.NetUsageWords)

	enc.Write(status)
	enc.Write(cpuUsageUs)
	enc.Write(netUsageWords)

	if t.Transaction.UnpackedTrx.ID() != crypto.NewSha256Nil() {
		trxID, _ := rlp.EncodeToBytes(t.Transaction.UnpackedTrx.ID())
		enc.Write(trxID)
	} else {
		packedTrx, _ := rlp.EncodeToBytes(t.Transaction.PackedDigest())
		enc.Write(packedTrx)
	}

	return *crypto.NewSha256Byte(enc.Sum(nil))
}

var _ Entity = &TransactionTrace{}

type TransactionTrace struct {
	ID           IdType             `serialize:"true" json:"-" eos:"-"`
	Hash         TransactionIdType  `serialize:"true" json:"id"`
	BlockNum     uint64             `serialize:"true" json:"block_num"`
	BlockTime    TimePoint          `serialize:"true" json:"block_time"`
	Receipt      TransactionReceipt `serialize:"true" json:"receipt"`
	Elapsed      Microseconds       `serialize:"true" json:"elapsed"`
	NetUsage     uint64             `serialize:"true" json:"-"`
	Scheduled    bool               `serialize:"true" json:"scheduled"`
	ActionTraces []ActionTrace      `serialize:"true" json:"action_traces"`
	Except       error              `msg:"-" json:"-"`
	ErrorCode    uint64             `json:"-"`
}

func (a TransactionTrace) GetId() []byte {
	return a.ID.ToBytes()
}

func (a TransactionTrace) GetIndexes() map[string]EntityIndex {
	return map[string]EntityIndex{
		"id": {
			Fields: []string{"ID"},
		},
		"byHash": {
			Fields: []string{"Hash"},
		},
	}
}

func (a TransactionTrace) GetObjectType() uint8 {
	return TransactionType
}

type RamDelta struct {
	Account name.AccountName `serialize:"true" json:"account"`
	Delta   int64            `serialize:"true" json:"delta"`
}

type ActionTrace struct {
	AccountRamDeltas     []RamDelta        `serialize:"true" json:"account_ram_deltas"`
	ActionOrdinal        Vuint32           `serialize:"true" json:"action_ordinal"`
	CreatorActionOrdinal Vuint32           `serialize:"true" json:"creator_action_ordinal"`
	Receipt              ActionReceipt     `serialize:"true" json:"receipt"`
	Receiver             name.ActionName   `serialize:"true" json:"receiver"`
	Action               Action            `serialize:"true" json:"act"`
	ContextFree          bool              `serialize:"true" json:"context_free"`
	Elapsed              uint64            `serialize:"true" json:"elapsed"`
	Console              string            `serialize:"true" json:"-"`
	TransactionId        TransactionIdType `serialize:"true" json:"trx_id"`
	BlockNum             uint64            `serialize:"true" json:"block_num"`
	BlockTime            TimePoint         `serialize:"true" json:"block_time"`
	Except               error             `msg:"-" json:"-"`
	ErrorCode            uint64            `serialize:"true" json:"-"`
}

func NewActionTrace(trace *TransactionTrace, action Action, receiver name.AccountName, contextFree bool, actionOrdinal Vuint32, creatorActionOrdinal Vuint32) *ActionTrace {
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
