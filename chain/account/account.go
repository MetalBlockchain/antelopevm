package account

import (
	"github.com/MetalBlockchain/antelopevm/chain/block"
	"github.com/MetalBlockchain/antelopevm/chain/entity"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/time"
	"github.com/MetalBlockchain/antelopevm/chain/types"
)

const (
	PrivilegedFlag uint32 = 1
)

var _ entity.Entity = &Account{}

type Account struct {
	ID           types.IdType         `serialize:"true"`
	Name         name.AccountName     `serialize:"true"`
	CreationDate block.BlockTimeStamp `serialize:"true"`
	Abi          types.HexBytes       `serialize:"true"`
}

func (a Account) GetId() []byte {
	return a.ID.ToBytes()
}

func (a Account) GetIndexes() map[string]entity.EntityIndex {
	return map[string]entity.EntityIndex{
		"id": {
			Fields: []string{"ID"},
		},
		"byName": {
			Fields: []string{"Name"},
		},
	}
}

func (a Account) GetObjectType() uint8 {
	return entity.AccountType
}

var _ entity.Entity = &AccountMetaDataObject{}

type AccountMetaDataObject struct {
	ID             types.IdType     `serialize:"true"`
	Name           name.AccountName `serialize:"true"`
	RecvSequence   uint64           `serialize:"true"`
	AuthSequence   uint64           `serialize:"true"`
	CodeSequence   uint64           `serialize:"true"`
	AbiSequence    uint64           `serialize:"true"`
	CodeHash       types.DigestType `serialize:"true"`
	LastCodeUpdate time.TimePoint   `serialize:"true"`
	Flags          uint32           `serialize:"true"`
	VmType         uint8            `serialize:"true"`
	VmVersion      uint8            `serialize:"true"`
}

func (a AccountMetaDataObject) IsPrivileged() bool {
	return (a.Flags & PrivilegedFlag) != 0
}

func (a *AccountMetaDataObject) SetPrivileged(privileged bool) {
	if privileged {
		a.Flags = a.Flags | PrivilegedFlag
	} else {
		a.Flags = a.Flags &^ PrivilegedFlag
	}
}

func (a AccountMetaDataObject) GetId() []byte {
	return a.ID.ToBytes()
}

func (a AccountMetaDataObject) GetIndexes() map[string]entity.EntityIndex {
	return map[string]entity.EntityIndex{
		"id": {
			Fields: []string{"ID"},
		},
		"byName": {
			Fields: []string{"Name"},
		},
	}
}

func (a AccountMetaDataObject) GetObjectType() uint8 {
	return entity.AccountMetaDataObjectType
}

var _ entity.Entity = &AccountRamCorrectionObject{}

type AccountRamCorrectionObject struct {
	ID            types.IdType     `serialize:"true"`
	Name          name.AccountName `serialize:"true"`
	RamCorrection uint64           `serialize:"true"`
}

func (a AccountRamCorrectionObject) GetId() []byte {
	return a.ID.ToBytes()
}

func (a AccountRamCorrectionObject) GetIndexes() map[string]entity.EntityIndex {
	return map[string]entity.EntityIndex{
		"id": {
			Fields: []string{"ID"},
		},
		"byName": {
			Fields: []string{"Name"},
		},
	}
}

func (a AccountRamCorrectionObject) GetObjectType() uint8 {
	return entity.AccountRamCorrectionObjectType
}
