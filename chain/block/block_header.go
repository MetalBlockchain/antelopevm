package block

import (
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/producer"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/utils"
)

type BlockHeader struct {
	Timestamp             BlockTimeStamp             `json:"timestamp"`
	Producer              name.AccountName           `json:"producer"`
	Confirmed             uint16                     `json:"confirmed"`
	Previous              crypto.Sha256              `json:"previous"`
	TransactionMerkleRoot crypto.Sha256              `json:"transaction_mroot"`
	ActionMerkleRoot      crypto.Sha256              `json:"action_mroot"`
	ScheduleVersion       uint32                     `json:"schedule_version"`
	NewProducers          *producer.ProducerSchedule `json:"new_producers" eos:"optional"`
	Extensions            []types.Extension          `json:"header_extensions"`
}

func (b *BlockHeader) Digest() *crypto.Sha256 {
	return crypto.Hash256(b)
}

func (b *BlockHeader) NumFromId(id crypto.Sha256) uint32 {
	return utils.EndianReverseU32(uint32(id.Hash[0]))
}

func (b *BlockHeader) BlockNum() uint32 {
	return b.NumFromId(b.Previous) + 1
}

func (b *BlockHeader) CalculateId() *crypto.Sha256 {
	result := b.Digest()
	result.Hash[0] &= 0xffffffff00000000
	result.Hash[0] += uint64(utils.EndianReverseU32(b.BlockNum()))
	return result
}

type SignedBlockHeader struct {
	BlockHeader
	ProducerSignature ecc.Signature `json:"producer_signature"`
}
