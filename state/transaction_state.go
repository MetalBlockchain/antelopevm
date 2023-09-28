package state

import (
	"github.com/MetalBlockchain/antelopevm/chain/transaction"
	"github.com/MetalBlockchain/antelopevm/chain/types"
)

func (s *Session) FindTransaction(id types.IdType) (*transaction.TransactionTrace, error) {
	key := getObjectKeyByIndex(&transaction.TransactionTrace{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &transaction.TransactionTrace{}

	if _, err := Codec.Unmarshal(data, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) CreateTransactionObject(in *transaction.TransactionObject) error {
	return s.create(false, nil, in)
}
