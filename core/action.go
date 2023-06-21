package core

import (
	"github.com/MetalBlockchain/antelopevm/core/authority"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
)

//go:generate msgp
type Action struct {
	Account       name.AccountName            `serialize:"true" json:"account"`
	Name          name.ActionName             `serialize:"true" json:"name"`
	Authorization []authority.PermissionLevel `serialize:"true" json:"authorization,omitempty"`
	Data          HexBytes                    `serialize:"true" json:"hex_data"`
	ParsedData    map[string]interface{}      `json:"data,omitempty" eos:"-"`
}

func (a Action) DataAs(t interface{}) {
	err := rlp.DecodeBytes(a.Data, t)
	if err != nil {
		panic(err.Error())
	}
}

type ContractTypesInterface interface {
	GetAccount() name.AccountName
	GetName() name.ActionName
}

type ActionReceipt struct {
	Receiver       name.AccountName `serialize:"true" json:"receiver"`
	ActDigest      crypto.Sha256    `serialize:"true" json:"act_digest"`
	GlobalSequence uint64           `serialize:"true" json:"global_sequence"`
	RecvSequence   uint64           `serialize:"true" json:"recv_sequence"`
	AuthSequence   AuthSequenceSet  `serialize:"true" json:"auth_sequence"`
	CodeSequence   Vuint32          `serialize:"true" json:"code_sequence"`
	AbiSequence    Vuint32          `serialize:"true" json:"abi_sequence"`
}

func (a *ActionReceipt) Digest() crypto.Sha256 {
	return *crypto.Hash256(a)
}
