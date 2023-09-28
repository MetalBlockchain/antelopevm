package resource

import (
	"github.com/MetalBlockchain/antelopevm/chain/entity"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/types"
)

var _ entity.Entity = &ResourceUsage{}

type ResourceUsage struct {
	ID       types.IdType     `serialize:"true"`
	Owner    name.AccountName `serialize:"true"`
	NetUsage UsageAccumulator `serialize:"true"`
	CpuUsage UsageAccumulator `serialize:"true"`
	RamUsage uint64           `serialize:"true"`
}

// GetId implements core.Entity
func (r *ResourceUsage) GetId() []byte {
	return r.ID.ToBytes()
}

// GetIndexes implements core.Entity
func (*ResourceUsage) GetIndexes() map[string]entity.EntityIndex {
	return map[string]entity.EntityIndex{
		"id": {
			Fields: []string{"ID"},
		},
		"byOwner": {
			Fields: []string{"Owner"},
		},
	}
}

// GetObjectType implements core.Entity
func (*ResourceUsage) GetObjectType() uint8 {
	return entity.ResourceUsageType
}
