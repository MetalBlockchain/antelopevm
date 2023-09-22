package block

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/core/producer"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/utils"
)

type BlockHeader struct {
	Timestamp             uint32
	Producer              name.AccountName
	Confirmed             uint16
	Previous              crypto.Sha256
	TransactionMerkleRoot crypto.Sha256
	ActionMerkleRoot      crypto.Sha256
	ScheduleVersion       uint32
	Producers             *producer.ProducerSchedule `eos:"optional"`
	Extensions            []core.Extension           `eos:"optional"`
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
