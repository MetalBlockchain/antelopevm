package resource

import (
	"github.com/MetalBlockchain/antelopevm/chain/entity"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/types"
)

var _ entity.Entity = &ResourceLimits{}

type ResourceLimits struct {
	ID        types.IdType
	Owner     name.AccountName
	Pending   bool
	NetWeight int64
	CpuWeight int64
	RamBytes  int64
}

// GetId implements core.Entity
func (rl *ResourceLimits) GetId() []byte {
	return rl.ID.ToBytes()
}

// GetIndexes implements core.Entity
func (*ResourceLimits) GetIndexes() map[string]entity.EntityIndex {
	return map[string]entity.EntityIndex{
		"id": {
			Fields: []string{"ID"},
		},
		"byOwner": {
			Fields: []string{"Pending", "Owner"},
		},
	}
}

// GetObjectType implements core.Entity
func (*ResourceLimits) GetObjectType() uint8 {
	return entity.ResourceLimitType
}
