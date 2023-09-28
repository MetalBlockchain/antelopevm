package table

import (
	"github.com/MetalBlockchain/antelopevm/chain/entity"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/resource"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/config"
	"github.com/MetalBlockchain/antelopevm/math"
)

var (
	_ entity.Entity = &Index64Object{}
	_ entity.Entity = &Index128Object{}
	_ entity.Entity = &Index256Object{}
	_ entity.Entity = &IndexDoubleObject{}
	_ entity.Entity = &IndexLongDoubleObject{}

	Index64ObjectBillableSize         = resource.NewBillableSize(24 + 8 + uint64(config.OverheadPerRowPerIndexRamBytes*3))
	Index128ObjectBillableSize        = resource.NewBillableSize(24 + 16 + uint64(config.OverheadPerRowPerIndexRamBytes*3))
	Index256ObjectBillableSize        = resource.NewBillableSize(24 + 32 + uint64(config.OverheadPerRowPerIndexRamBytes*3))
	IndexDoubleObjectBillableSize     = resource.NewBillableSize(24 + 8 + uint64(config.OverheadPerRowPerIndexRamBytes*3))
	IndexLongDoubleObjectBillableSize = resource.NewBillableSize(24 + 16 + uint64(config.OverheadPerRowPerIndexRamBytes*3))
)

type Index64Object struct {
	ID           types.IdType
	TableID      types.IdType
	PrimaryKey   uint64
	Payer        name.AccountName
	SecondaryKey uint64
}

func (kv Index64Object) GetId() []byte {
	return kv.ID.ToBytes()
}

func (kv Index64Object) GetIndexes() map[string]entity.EntityIndex {
	return map[string]entity.EntityIndex{
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
	return entity.IndexObjectType
}

type Index128Object struct {
	ID           types.IdType
	TableID      types.IdType
	PrimaryKey   uint64
	Payer        name.AccountName
	SecondaryKey math.Uint128
}

func (kv Index128Object) GetId() []byte {
	return kv.ID.ToBytes()
}

func (kv Index128Object) GetIndexes() map[string]entity.EntityIndex {
	return map[string]entity.EntityIndex{
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
	return entity.IndexObjectType
}

type Index256Object struct {
	ID           types.IdType
	TableID      types.IdType
	PrimaryKey   uint64
	Payer        name.AccountName
	SecondaryKey math.Uint256
}

func (kv Index256Object) GetId() []byte {
	return kv.ID.ToBytes()
}

func (kv Index256Object) GetIndexes() map[string]entity.EntityIndex {
	return map[string]entity.EntityIndex{
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
	return entity.IndexObjectType
}

type IndexDoubleObject struct {
	ID           types.IdType
	TableID      types.IdType
	PrimaryKey   uint64
	Payer        name.AccountName
	SecondaryKey float64
}

func (kv IndexDoubleObject) GetId() []byte {
	return kv.ID.ToBytes()
}

func (kv IndexDoubleObject) GetIndexes() map[string]entity.EntityIndex {
	return map[string]entity.EntityIndex{
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
	return entity.IndexObjectType
}

type IndexLongDoubleObject struct {
	ID           types.IdType
	TableID      types.IdType
	PrimaryKey   uint64
	Payer        name.AccountName
	SecondaryKey math.Float128
}

func (kv IndexLongDoubleObject) GetId() []byte {
	return kv.ID.ToBytes()
}

func (kv IndexLongDoubleObject) GetIndexes() map[string]entity.EntityIndex {
	return map[string]entity.EntityIndex{
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
	return entity.IndexObjectType
}
