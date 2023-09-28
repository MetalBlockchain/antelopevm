package producer

import (
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
)

type ProducerKey struct {
	ProducerName    name.AccountName `json:"producer_name"`
	BlockSigningKey ecc.PublicKey    `json:"block_signing_key"`
}

func (p ProducerKey) Equal(other ProducerKey) bool {
	return p.ProducerName == other.ProducerName && p.BlockSigningKey.Compare(other.BlockSigningKey)
}

type SharedBlockSigningAuthority struct {
	Threshold uint32
}
