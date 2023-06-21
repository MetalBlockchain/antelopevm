package core

import (
	"encoding/hex"

	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
)

//go:generate msgp
type BlockHash [32]byte

func (b BlockHash) Hex() string {
	return hex.EncodeToString(b[:])
}

type BlockTimeStamp = TimePoint
type BlockStatus uint8

const (
	BlockStatusProcessing BlockStatus = 0
	BlockStatusAccepted   BlockStatus = 1
	BlockStatusRejected   BlockStatus = 2
)

func (t BlockTimeStamp) ToTimePoint(blockIntervalMs int64, blockTimestampEpochMs int64) TimePoint {
	msec := int64(t) * int64(blockIntervalMs)
	msec += int64(blockTimestampEpochMs)
	return TimePoint(Milliseconds(msec))
}

type BlockHeader struct {
	Created               TimePoint        `serialize:"true"`
	Producer              name.AccountName `serialize:"true"`
	Confirmed             uint16           `serialize:"true"`
	PreviousBlockHash     BlockHash        `serialize:"true"`
	Index                 uint64           `serialize:"true"`
	TransactionMerkleRoot crypto.Sha256    `serialize:"true"`
	ActionMerkleRoot      crypto.Sha256    `serialize:"true"`
	ScheduleVersion       uint32           `serialize:"true"`
	HeaderExtensions      []Extension      `serialize:"true"`
}

type SignedBlockHeader struct {
	BlockHeader       `serialize:"true"`
	ProducerSignature ecc.Signature `serialize:"true"`
}
