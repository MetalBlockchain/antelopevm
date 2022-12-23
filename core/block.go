package core

import (
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/metalgo/ids"
)

type BlockId = crypto.Sha256
type BlockTimeStamp TimePoint

type BlockHeader struct {
	Created               TimePoint     `serialize:"true"`
	Producer              AccountName   `serialize:"true"`
	Confirmed             uint16        `serialize:"true"`
	PreviousBlock         ids.ID        `serialize:"true"`
	Index                 uint64        `serialize:"true"`
	TransactionMerkleRoot crypto.Sha256 `serialize:"true"`
	ActionMerkleRoot      crypto.Sha256 `serialize:"true"`
	ScheduleVersion       uint32        `serialize:"true"`
	HeaderExtensions      []Extension   `serialize:"true"`
}

type SignedBlockHeader struct {
	BlockHeader       `serialize:"true"`
	ProducerSignature ecc.Signature `serialize:"true"`
}
