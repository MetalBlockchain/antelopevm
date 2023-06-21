package state

import (
	"encoding/binary"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/account"
	"github.com/MetalBlockchain/antelopevm/core/resource"
	"github.com/MetalBlockchain/metalgo/cache"
	"github.com/dgraph-io/badger/v3"
)

type Session struct {
	state               *State
	transaction         *badger.Txn
	accountCache        *cache.LRU[core.IdType, *account.Account]
	tableCache          *cache.LRU[core.IdType, *core.Table]
	kvCache             *cache.LRU[core.IdType, *core.KeyValue]
	indexObjectCache    *cache.LRU[core.IdType, interface{}]
	resourceUsageCache  *cache.LRU[core.IdType, *resource.ResourceUsage]
	resourceLimitsCache *cache.LRU[core.IdType, *resource.ResourceLimits]
}

func NewSession(state *State, transaction *badger.Txn) *Session {
	session := &Session{
		state:               state,
		transaction:         transaction,
		accountCache:        &cache.LRU[core.IdType, *account.Account]{Size: blockCacheSize},
		tableCache:          &cache.LRU[core.IdType, *core.Table]{Size: blockCacheSize},
		kvCache:             &cache.LRU[core.IdType, *core.KeyValue]{Size: blockCacheSize},
		indexObjectCache:    &cache.LRU[core.IdType, interface{}]{Size: blockCacheSize},
		resourceUsageCache:  &cache.LRU[core.IdType, *resource.ResourceUsage]{Size: blockCacheSize},
		resourceLimitsCache: &cache.LRU[core.IdType, *resource.ResourceLimits]{Size: blockCacheSize},
	}

	return session
}

func (s *Session) create(incrementId bool, setId func(core.IdType) error, in core.Entity) error {
	if incrementId {
		id, err := s.increment([]byte{in.GetObjectType()})

		if err != nil {
			return err
		}

		if err := setId(core.IdType(id)); err != nil {
			return err
		}
	}

	bytes, err := in.MarshalMsg(nil)

	if err != nil {
		return err
	}

	keys := getObjectKeys(in)

	for index, key := range keys {
		if index == "id" {
			if err := s.transaction.Set(key, bytes); err != nil {
				return err
			}
		} else {
			if err := s.transaction.Set(key, in.GetId()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Session) modify(in core.Entity, modifyFunc func()) error {
	keys := getObjectKeys(in)
	modifyFunc()
	bytes, err := in.MarshalMsg(nil)

	if err != nil {
		return err
	}

	for index, key := range keys {
		if index == "id" {
			if err := s.transaction.Set(key, bytes); err != nil {
				return err
			}
		} else {
			if err := s.transaction.Set(key, in.GetId()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Session) remove(in core.Entity) error {
	keys := getObjectKeys(in)

	for _, key := range keys {
		if err := s.transaction.Delete(key); err != nil {
			return err
		}
	}

	return nil
}

func (s *Session) increment(key []byte) (uint64, error) {
	item, err := s.transaction.Get(key)

	if err != nil {
		if err := s.transaction.Set(key, uint64ToBytes(0)); err != nil {
			return 0, err
		}

		return 0, nil
	}

	value, err := item.ValueCopy(nil)

	if err != nil {
		return 0, err
	}

	id := bytesToUint64(value) + 1

	if err := s.transaction.Set(key, uint64ToBytes(id)); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Session) Commit() error {
	return s.transaction.Commit()
}

func (s *Session) Discard() {
	s.transaction.Discard()
}

func uint64ToBytes(i uint64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], i)
	return buf[:]
}

func bytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
