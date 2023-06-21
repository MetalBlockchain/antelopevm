package state

import (
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/dgraph-io/badger/v3"
)

var (
	initializedKey = []byte("initialized")
)

type State struct {
	db           *badger.DB
	sequences    map[string]*badger.Sequence
	lastAccepted ids.ID
	vm           VM
}

func NewState(vm VM, db *badger.DB) *State {
	return &State{
		vm:        vm,
		db:        db,
		sequences: make(map[string]*badger.Sequence),
	}
}

func (s *State) CreateSession(update bool) *Session {
	transaction := s.db.NewTransaction(update)

	return NewSession(s, transaction)
}

func (s *State) Close() error {
	return s.db.Close()
}

func (s *State) IsInitialized() (bool, error) {
	err := s.db.View(func(txn *badger.Txn) error {
		if _, err := txn.Get(initializedKey); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *State) SetInitialized() error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(initializedKey, []byte{1})
	})
}
