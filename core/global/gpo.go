package global

import (
	"github.com/MetalBlockchain/antelopevm/config"
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/producer"
)

var _ core.Entity = &GlobalPropertyObject{}

//go:generate msgp
type GlobalPropertyObject struct {
	ID                       core.IdType
	ProposedScheduleBlockNum uint64
	ProposedSchedule         producer.ProducerSchedule
	Configuration            config.ChainConfig
	ChainId                  core.ChainIdType
	WasmConfiguration        config.WasmConfig
}

// GetId implements core.Entity
func (gpo *GlobalPropertyObject) GetId() []byte {
	return gpo.ID.ToBytes()
}

// GetIndexes implements core.Entity
func (gpo *GlobalPropertyObject) GetIndexes() map[string]core.EntityIndex {
	return map[string]core.EntityIndex{
		"id": {
			Fields: []string{"ID"},
		},
	}
}

// GetObjectType implements core.Entity
func (gpo *GlobalPropertyObject) GetObjectType() uint8 {
	return core.GlobalPropertyObjectType
}
