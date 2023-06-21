package resource

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/name"
)

var _ core.Entity = &ResourceLimits{}

//go:generate msgp
type ResourceLimits struct {
	ID        core.IdType
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
func (*ResourceLimits) GetIndexes() map[string]core.EntityIndex {
	return map[string]core.EntityIndex{
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
	return core.ResourceLimitType
}
