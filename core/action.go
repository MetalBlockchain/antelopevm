package core

import (
	"encoding/json"

	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
)

type Action struct {
	Account       AccountName       `serialize:"true" json:"account"`
	Name          ActionName        `serialize:"true" json:"name"`
	Authorization []PermissionLevel `serialize:"true" json:"authorization,omitempty"`
	Data          HexBytes          `serialize:"true" json:"hex_data"`
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

type AuthSequence struct {
	Account  AccountName
	Sequence uint64
}

type AuthSequenceSet []AuthSequence

func (a *AuthSequenceSet) Set(account AccountName, sequence uint64) {
	for k, v := range a {
		if v.Account == account {
			a[k].Sequence = sequence
			return
		}
	}

	a = append(a, AuthSequence{account, sequence})
}

func (a AuthSequenceSet) MarshalJSON() ([]byte, error) {
	var out [][]interface{}

	for k, v := range a {
		out[k] = []interface{}{k, v}
	}

	return json.Marshal(out)
}

type ActionReceipt struct {
	Receiver       AccountName     `serialize:"true" json:"receiver"`
	ActDigest      crypto.Sha256   `serialize:"true" json:"act_digest"`
	GlobalSequence uint64          `serialize:"true" json:"global_sequence"`
	RecvSequence   uint64          `serialize:"true" json:"recv_sequence"`
	AuthSequence   AuthSequenceSet `serialize:"true" json:"auth_sequence"`
	CodeSequence   Vuint32         `serialize:"true" json:"code_sequence"`
	AbiSequence    Vuint32         `serialize:"true" json:"abi_sequence"`
}

func (a *ActionReceipt) Digest() crypto.Sha256 {
	return *crypto.Hash256(a)
}
