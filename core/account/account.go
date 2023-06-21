package account

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto"
)

var _ core.Entity = &Account{}

//go:generate msgp
type Account struct {
	ID             core.IdType      `serialize:"true"`
	Name           name.AccountName `serialize:"true"`
	VmType         uint8            `serialize:"true"`
	VmVersion      uint8            `serialize:"true"`
	Privileged     bool             `serialize:"true"`
	LastCodeUpdate core.TimePoint   `serialize:"true"`
	CodeVersion    crypto.Sha256    `serialize:"true"`
	CreationDate   core.TimePoint   `serialize:"true"`
	Code           core.HexBytes    `serialize:"true"`
	Abi            core.HexBytes    `serialize:"true"`
	AbiVersion     crypto.Sha256    `serialize:"true"`
	RecvSequence   uint64           `serialize:"true"`
	AuthSequence   uint64           `serialize:"true"`
	CodeSequence   uint64           `serialize:"true"`
	AbiSequence    uint64           `serialize:"true"`
}

func (a Account) GetId() []byte {
	return a.ID.ToBytes()
}

func (a Account) GetIndexes() map[string]core.EntityIndex {
	return map[string]core.EntityIndex{
		"id": {
			Fields: []string{"ID"},
		},
		"byName": {
			Fields: []string{"Name"},
		},
	}
}

func (a Account) GetObjectType() uint8 {
	return core.AccountType
}
