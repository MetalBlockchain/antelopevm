package global

import (
	"github.com/MetalBlockchain/antelopevm/chain/entity"
	"github.com/MetalBlockchain/antelopevm/chain/producer"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/config"
)

var _ entity.Entity = &GlobalPropertyObject{}

//go:generate msgp
type GlobalPropertyObject struct {
	ID                       types.IdType
	ProposedScheduleBlockNum uint64
	ProposedSchedule         producer.ProducerSchedule
	Configuration            config.ChainConfig
	ChainId                  types.ChainIdType
	WasmConfiguration        config.WasmConfig
}

// GetId implements core.Entity
func (gpo *GlobalPropertyObject) GetId() []byte {
	return gpo.ID.ToBytes()
}

// GetIndexes implements core.Entity
func (gpo *GlobalPropertyObject) GetIndexes() map[string]entity.EntityIndex {
	return map[string]entity.EntityIndex{
		"id": {
			Fields: []string{"ID"},
		},
	}
}

// GetObjectType implements core.Entity
func (gpo *GlobalPropertyObject) GetObjectType() uint8 {
	return entity.GlobalPropertyObjectType
}
