package account

import (
	"github.com/MetalBlockchain/antelopevm/chain/entity"
	"github.com/MetalBlockchain/antelopevm/chain/types"
)

var _ entity.Entity = &CodeObject{}

type CodeObject struct {
	ID             types.IdType     `serialize:"true"`
	CodeHash       types.DigestType `serialize:"true"`
	Code           types.HexBytes   `serialize:"true"`
	CodeRefCount   uint64           `serialize:"true"`
	FirstBlockUsed uint32           `serialize:"true"`
	VmType         uint8            `serialize:"true"`
	VmVersion      uint8            `serialize:"true"`
}

func (a CodeObject) GetId() []byte {
	return a.ID.ToBytes()
}

func (a CodeObject) GetIndexes() map[string]entity.EntityIndex {
	return map[string]entity.EntityIndex{
		"id": {
			Fields: []string{"ID"},
		},
		"byCodeHash": {
			Fields: []string{"Code", "VmType", "VmVersion"},
		},
	}
}

func (a CodeObject) GetObjectType() uint8 {
	return entity.CodeObjectType
}
