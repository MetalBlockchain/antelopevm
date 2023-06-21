package resource

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/name"
)

var _ core.Entity = &ResourceUsage{}

//go:generate msgp
type ResourceUsage struct {
	ID       core.IdType      `serialize:"true"`
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
func (*ResourceUsage) GetIndexes() map[string]core.EntityIndex {
	return map[string]core.EntityIndex{
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
	return core.ResourceUsageType
}
