package core

import "github.com/MetalBlockchain/antelopevm/core/name"

var _ Entity = &Table{}
var _ Entity = &KeyValue{}

//go:generate msgp
//msgp:ignore IndexObject
type Table struct {
	ID    IdType           `serialize:"true"`
	Code  name.AccountName `serialize:"true"`
	Scope name.ScopeName   `serialize:"true"`
	Table name.TableName   `serialize:"true"`
	Payer name.AccountName `serialize:"true"`
	Count uint32           `serialize:"true"`
}

func (t Table) GetId() []byte {
	return t.ID.ToBytes()
}

func (t Table) GetIndexes() map[string]EntityIndex {
	return map[string]EntityIndex{
		"id": {
			Fields: []string{"ID"},
		},
		"byCodeScopeTable": {
			Fields: []string{"Code", "Scope", "Table"},
		},
	}
}

func (t Table) GetObjectType() uint8 {
	return TableType
}

type KeyValue struct {
	ID         IdType           `serialize:"true"`
	TableID    IdType           `serialize:"true"`
	PrimaryKey uint64           `serialize:"true"`
	Payer      name.AccountName `serialize:"true"`
	Value      HexBytes         `serialize:"true"`
}

func (kv KeyValue) GetId() []byte {
	return kv.ID.ToBytes()
}

func (kv KeyValue) GetIndexes() map[string]EntityIndex {
	return map[string]EntityIndex{
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
	return KeyValueType
}
