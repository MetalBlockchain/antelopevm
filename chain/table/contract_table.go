package table

import (
	"github.com/MetalBlockchain/antelopevm/chain/entity"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/resource"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/config"
)

var _ entity.Entity = &Table{}
var _ entity.Entity = &KeyValue{}

var (
	TableIdObjectBillableSize  = resource.NewBillableSize(44 + uint64(config.OverheadPerRowPerIndexRamBytes*2))
	KeyValueObjectBillableSize = resource.NewBillableSize(32 + 8 + 4 + uint64(config.OverheadPerRowPerIndexRamBytes*2))
)

type Table struct {
	ID    types.IdType     `serialize:"true"`
	Code  name.AccountName `serialize:"true"`
	Scope name.ScopeName   `serialize:"true"`
	Table name.TableName   `serialize:"true"`
	Payer name.AccountName `serialize:"true"`
	Count uint32           `serialize:"true"`
}

func (t Table) GetId() []byte {
	return t.ID.ToBytes()
}

func (t Table) GetIndexes() map[string]entity.EntityIndex {
	return map[string]entity.EntityIndex{
		"id": {
			Fields: []string{"ID"},
		},
		"byCodeScopeTable": {
			Fields: []string{"Code", "Scope", "Table"},
		},
	}
}

func (t Table) GetObjectType() uint8 {
	return entity.TableType
}

type KeyValue struct {
	ID         types.IdType     `serialize:"true"`
	TableID    types.IdType     `serialize:"true"`
	PrimaryKey uint64           `serialize:"true"`
	Payer      name.AccountName `serialize:"true"`
	Value      types.HexBytes   `serialize:"true"`
}

func (kv KeyValue) GetId() []byte {
	return kv.ID.ToBytes()
}

func (kv KeyValue) GetIndexes() map[string]entity.EntityIndex {
	return map[string]entity.EntityIndex{
		"id": {
			Name:   "id",
			Fields: []string{"ID"},
		},
		"byScopePrimary": {
			Name:   "byScopePrimary",
			Fields: []string{"TableID", "PrimaryKey"},
		},
	}
}

func (kv KeyValue) GetObjectType() uint8 {
	return entity.KeyValueType
}
