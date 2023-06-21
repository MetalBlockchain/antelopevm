package transaction

import "github.com/MetalBlockchain/antelopevm/core"

var _ core.Entity = &TransactionObject{}

//go:generate msgp
type TransactionObject struct {
	ID         core.IdType
	Expiration core.TimePointSec
	TrxId      core.TransactionIdType
}

func (p TransactionObject) GetId() []byte {
	return p.ID.ToBytes()
}

func (p TransactionObject) GetIndexes() map[string]core.EntityIndex {
	return map[string]core.EntityIndex{
		"id": {
			Fields: []string{"ID"},
		},
		"byTrxId": {
			Fields: []string{"TrxId"},
		},
		"byExpiration": {
			Fields: []string{"Expiration", "TrxId"},
		},
	}
}

func (p TransactionObject) GetObjectType() uint8 {
	return core.TransactionObjectType
}
