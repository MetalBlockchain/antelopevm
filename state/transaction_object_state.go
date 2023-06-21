package state

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/transaction"
)

func (s *Session) FindTransactionObject(id core.IdType) (*transaction.TransactionObject, error) {
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

	if _, err := out.UnmarshalMsg(data); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) FindTransactionByHash(hash core.TransactionIdType) (*core.TransactionTrace, error) {
	key := getObjectKeyByIndex(&core.TransactionTrace{Hash: hash}, "byHash")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return s.FindTransaction(core.NewIdType(data))
}

func (s *Session) CreateTransaction(in *core.TransactionTrace) error {
	return s.create(false, nil, in)
}
