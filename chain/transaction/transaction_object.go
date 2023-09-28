package transaction

import (
	"github.com/MetalBlockchain/antelopevm/chain/entity"
	"github.com/MetalBlockchain/antelopevm/chain/time"
	"github.com/MetalBlockchain/antelopevm/chain/types"
)

var _ entity.Entity = &TransactionObject{}

type TransactionObject struct {
	ID         types.IdType
	Expiration time.TimePointSec
	TrxId      TransactionIdType
}

func (p TransactionObject) GetId() []byte {
	return p.ID.ToBytes()
}

func (p TransactionObject) GetIndexes() map[string]entity.EntityIndex {
	return map[string]entity.EntityIndex{
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
	return entity.TransactionObjectType
}
