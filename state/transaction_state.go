package state

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/transaction"
)

func (s *Session) FindTransaction(id core.IdType) (*core.TransactionTrace, error) {
	key := getObjectKeyByIndex(&core.TransactionTrace{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &core.TransactionTrace{}

	if _, err := out.UnmarshalMsg(data); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) CreateTransactionObject(in *transaction.TransactionObject) error {
	return s.create(false, nil, in)
}
