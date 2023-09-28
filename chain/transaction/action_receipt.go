package transaction

import (
	"github.com/MetalBlockchain/antelopevm/chain/authority"
	"github.com/MetalBlockchain/antelopevm/chain/fc"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto"
)

type ActionReceipt struct {
	Receiver       name.AccountName          `serialize:"true" json:"receiver"`
	ActDigest      types.DigestType          `serialize:"true" json:"act_digest"`
	GlobalSequence uint64                    `serialize:"true" json:"global_sequence"`
	RecvSequence   uint64                    `serialize:"true" json:"recv_sequence"`
	AuthSequence   authority.AuthSequenceSet `serialize:"true" json:"auth_sequence"`
	CodeSequence   fc.UnsignedInt            `serialize:"true" json:"code_sequence"`
	AbiSequence    fc.UnsignedInt            `serialize:"true" json:"abi_sequence"`
}

func (a *ActionReceipt) Digest() types.DigestType {
	return *crypto.Hash256(a)
}
