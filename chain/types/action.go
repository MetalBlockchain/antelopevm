package types

import (
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
)

type Action struct {
	Account       AccountName       `json:"account"`
	Name          ActionName        `json:"name"`
	Authorization []PermissionLevel `json:"authorization,omitempty"`
	Data          HexBytes          `json:"data"`
}

func (a Action) DataAs(t interface{}) {
	err := rlp.DecodeBytes(a.Data, t)
	if err != nil {
		panic(err.Error())
	}
}

type ContractTypesInterface interface {
	GetAccount() AccountName
	GetName() ActionName
}

type ActionReceipt struct {
	Receiver       AccountName   `json:"receiver"`
	ActDigest      crypto.Sha256 `json:"act_digest"`
	GlobalSequence uint64        `json:"global_sequence"`
	RecvSequence   uint64        `json:"recv_sequence"`
	AuthSequence   map[AccountName]uint64
	CodeSequence   Vuint32 `json:"code_sequence"`
	AbiSequence    Vuint32 `json:"abi_sequence"`
}

func (a *ActionReceipt) Digest() crypto.Sha256 {
	return *crypto.Hash256(a)
}
