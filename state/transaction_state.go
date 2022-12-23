package state

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/metalgo/database"
	"github.com/inconshreveable/log15"
)

var (
	transactionIdKey = []byte("Trx__id__")
)

type transactionState struct {
	db database.Database
}

var _ TransactionState = &transactionState{}

type TransactionState interface {
	GetTransaction(core.TransactionIdType) (*core.TransactionTrace, error)
	PutTransaction(*core.TransactionTrace) error
}

func NewTransactionState(db database.Database) TransactionState {
	return &transactionState{
		db: db,
	}
}

func (t *transactionState) GetTransaction(id core.TransactionIdType) (*core.TransactionTrace, error) {
	key := append(transactionIdKey, id.Bytes()...)
	wrappedBytes, err := t.db.Get(key)

	if err != nil {
		return nil, err
	}

	transaction := &core.TransactionTrace{}

	if _, err := Codec.Unmarshal(wrappedBytes, transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *transactionState) PutTransaction(trx *core.TransactionTrace) error {
	log15.Info("storing tx", "tx", trx)
	wrappedBytes, err := Codec.Marshal(CodecVersion, &trx)

	if err != nil {
		return err
	}

	batch := t.db.NewBatch()
	key := append(transactionIdKey, trx.Id.Bytes()...)
	batch.Put(key, wrappedBytes)

	return batch.Write()
}
