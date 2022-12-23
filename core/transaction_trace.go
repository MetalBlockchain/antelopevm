package core

import (
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
)

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

type TransactionTrace struct {
	Id           TransactionIdType  `serialize:"true" json:"id"`
	BlockNum     uint32             `serialize:"true" json:"block_num"`
	BlockTime    TimePoint          `serialize:"true" json:"block_time"`
	Receipt      TransactionReceipt `serialize:"true" json:"-"`
	Elapsed      Microseconds       `serialize:"true" json:"elapsed"`
	NetUsage     uint64             `serialize:"true" json:"-"`
	Scheduled    bool               `serialize:"true" json:"-"`
	ActionTraces []ActionTrace      `serialize:"true" json:"traces"`
	Except       error              `json:"-"`
	ErrorCode    uint64             `json:"-"`
}

type ActionTrace struct {
	ActionOrdinal        Vuint32           `serialize:"true" json:"action_ordinal"`
	CreatorActionOrdinal Vuint32           `serialize:"true" json:"creator_action_ordinal"`
	Receipt              ActionReceipt     `serialize:"true" json:"receipt"`
	Receiver             ActionName        `serialize:"true" json:"receiver"`
	Action               Action            `serialize:"true" json:"act"`
	ContextFree          bool              `serialize:"true" json:"context_free"`
	Elapsed              uint64            `serialize:"true" json:"elapsed"`
	Console              string            `serialize:"true" json:"-"`
	TransactionId        TransactionIdType `serialize:"true" json:"trx_id"`
	BlockNum             uint32            `serialize:"true" json:"block_num"`
	BlockTime            TimePoint         `serialize:"true" json:"block_time"`
	Except               error             `json:"-"`
	ErrorCode            uint64            `serialize:"true" json:"-"`
}

func NewActionTrace(trace *TransactionTrace, action Action, receiver AccountName, contextFree bool, actionOrdinal Vuint32, creatorActionOrdinal Vuint32) *ActionTrace {
	return &ActionTrace{
		ActionOrdinal:        actionOrdinal,
		CreatorActionOrdinal: creatorActionOrdinal,
		Receiver:             receiver,
		Action:               action,
		ContextFree:          contextFree,
		TransactionId:        trace.Id,
		BlockNum:             trace.BlockNum,
		BlockTime:            trace.BlockTime,
	}
}

func (a *ActionTrace) Digest() crypto.Sha256 {
	return *crypto.Hash256(a)
}
