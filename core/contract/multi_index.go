package contract

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/math"
)

var (
	_ core.Entity = &Index64Object{}
	_ core.Entity = &Index128Object{}
	_ core.Entity = &Index256Object{}
	_ core.Entity = &IndexDoubleObject{}
	_ core.Entity = &IndexLongDoubleObject{}
)

//go:generate msgp
type Index64Object struct {
	ID           core.IdType
	TableID      core.IdType
	PrimaryKey   uint64
	Payer        name.AccountName
	SecondaryKey uint64
}

func (kv Index64Object) GetId() []byte {
	return kv.ID.ToBytes()
}

func (kv Index64Object) GetIndexes() map[string]core.EntityIndex {
	return map[string]core.EntityIndex{
		"id": {
			Name:   "id",
			Fields: []string{"ID"},
		},
		"byPrimary": {
			Name:   "byPrimary",
			Fields: []string{"TableID", "PrimaryKey"},
		},
		"bySecondary": {
			Name:   "bySecondary",
			Fields: []string{"TableID", "SecondaryKey", "PrimaryKey"},
		},
	}
}

func (kv Index64Object) GetObjectType() uint8 {
	return core.IndexObjectType
}

type Index128Object struct {
	ID           core.IdType
	TableID      core.IdType
	PrimaryKey   uint64
	Payer        name.AccountName
	SecondaryKey math.Uint128
}

func (kv Index128Object) GetId() []byte {
	return kv.ID.ToBytes()
}

func (kv Index128Object) GetIndexes() map[string]core.EntityIndex {
	return map[string]core.EntityIndex{
		"id": {
			Name:   "id",
			Fields: []string{"ID"},
		},
		"byPrimary": {
			Name:   "byPrimary",
			Fields: []string{"TableID", "PrimaryKey"},
		},
		"bySecondary": {
			Name:   "bySecondary",
			Fields: []string{"TableID", "SecondaryKey", "PrimaryKey"},
		},
	}
}

func (kv Index128Object) GetObjectType() uint8 {
	return core.IndexObjectType
}

type Index256Object struct {
	ID           core.IdType
	TableID      core.IdType
	PrimaryKey   uint64
	Payer        name.AccountName
	SecondaryKey math.Uint256
}

func (kv Index256Object) GetId() []byte {
	return kv.ID.ToBytes()
}

func (kv Index256Object) GetIndexes() map[string]core.EntityIndex {
	return map[string]core.EntityIndex{
		"id": {
			Name:   "id",
			Fields: []string{"ID"},
		},
		"byPrimary": {
			Name:   "byPrimary",
			Fields: []string{"TableID", "PrimaryKey"},
		},
		"bySecondary": {
			Name:   "bySecondary",
			Fields: []string{"TableID", "SecondaryKey", "PrimaryKey"},
		},
	}
}

func (kv Index256Object) GetObjectType() uint8 {
	return core.IndexObjectType
}

type IndexDoubleObject struct {
	ID           core.IdType
	TableID      core.IdType
	PrimaryKey   uint64
	Payer        name.AccountName
	SecondaryKey float64
}

func (kv IndexDoubleObject) GetId() []byte {
	return kv.ID.ToBytes()
}

func (kv IndexDoubleObject) GetIndexes() map[string]core.EntityIndex {
	return map[string]core.EntityIndex{
		"id": {
			Name:   "id",
			Fields: []string{"ID"},
		},
		"byPrimary": {
			Name:   "byPrimary",
			Fields: []string{"TableID", "PrimaryKey"},
		},
		"bySecondary": {
			Name:   "bySecondary",
			Fields: []string{"TableID", "SecondaryKey", "PrimaryKey"},
		},
	}
}

func (kv IndexDoubleObject) GetObjectType() uint8 {
	return core.IndexObjectType
}

type IndexLongDoubleObject struct {
	ID           core.IdType
	TableID      core.IdType
	PrimaryKey   uint64
	Payer        name.AccountName
	SecondaryKey math.Float128
}

func (kv IndexLongDoubleObject) GetId() []byte {
	return kv.ID.ToBytes()
}

func (kv IndexLongDoubleObject) GetIndexes() map[string]core.EntityIndex {
	return map[string]core.EntityIndex{
		"id": {
			Name:   "id",
			Fields: []string{"ID"},
		},
		"byPrimary": {
			Name:   "byPrimary",
			Fields: []string{"TableID", "PrimaryKey"},
		},
		"bySecondary": {
			Name:   "bySecondary",
			Fields: []string{"TableID", "SecondaryKey", "PrimaryKey"},
		},
	}
}

func (kv IndexLongDoubleObject) GetObjectType() uint8 {
	return core.IndexObjectType
}
