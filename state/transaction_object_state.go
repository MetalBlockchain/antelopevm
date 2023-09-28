package state

import (
	"github.com/MetalBlockchain/antelopevm/chain/transaction"
	"github.com/MetalBlockchain/antelopevm/chain/types"
)

func (s *Session) FindTransactionObject(id types.IdType) (*transaction.TransactionObject, error) {
	key := getObjectKeyByIndex(&transaction.TransactionObject{ID: id}, "id")
	item, err := s.transaction.Get(key)
	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	out := &transaction.TransactionObject{}
	if _, err := Codec.Unmarshal(data, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) FindTransactionByHash(hash transaction.TransactionIdType) (*transaction.TransactionTrace, error) {
	key := getObjectKeyByIndex(&transaction.TransactionTrace{Hash: hash}, "byHash")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return s.FindTransaction(types.NewIdType(data))
}

func (s *Session) CreateTransaction(in *transaction.TransactionTrace) error {
	return s.create(false, nil, in)
}
