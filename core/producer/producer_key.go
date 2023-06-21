package producer

import (
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
)

//go:generate msgp
type ProducerKey struct {
	ProducerName    name.AccountName
	BlockSigningKey ecc.PublicKey
}

func (p ProducerKey) Equal(other ProducerKey) bool {
	return p.ProducerName == other.ProducerName && p.BlockSigningKey.Compare(other.BlockSigningKey)
}

type SharedBlockSigningAuthority struct {
	Threshold uint32
}
