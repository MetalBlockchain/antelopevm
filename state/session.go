package state

import (
	"encoding/binary"

	"github.com/MetalBlockchain/antelopevm/chain/account"
	"github.com/MetalBlockchain/antelopevm/chain/entity"
	"github.com/MetalBlockchain/antelopevm/chain/resource"
	"github.com/MetalBlockchain/antelopevm/chain/table"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/metalgo/cache"
	"github.com/dgraph-io/badger/v3"
)

type Session struct {
	state               *State
	transaction         *badger.Txn
	accountCache        *cache.LRU[types.IdType, *account.Account]
	tableCache          *cache.LRU[types.IdType, *table.Table]
	kvCache             *cache.LRU[types.IdType, *table.KeyValue]
	indexObjectCache    *cache.LRU[types.IdType, interface{}]
	resourceUsageCache  *cache.LRU[types.IdType, *resource.ResourceUsage]
	resourceLimitsCache *cache.LRU[types.IdType, *resource.ResourceLimits]
}

func NewSession(state *State, transaction *badger.Txn) *Session {
	session := &Session{
		state:               state,
		transaction:         transaction,
		accountCache:        &cache.LRU[types.IdType, *account.Account]{Size: blockCacheSize},
		tableCache:          &cache.LRU[types.IdType, *table.Table]{Size: blockCacheSize},
		kvCache:             &cache.LRU[types.IdType, *table.KeyValue]{Size: blockCacheSize},
		indexObjectCache:    &cache.LRU[types.IdType, interface{}]{Size: blockCacheSize},
		resourceUsageCache:  &cache.LRU[types.IdType, *resource.ResourceUsage]{Size: blockCacheSize},
		resourceLimitsCache: &cache.LRU[types.IdType, *resource.ResourceLimits]{Size: blockCacheSize},
	}

	return session
}

func (s *Session) create(incrementId bool, setId func(types.IdType) error, in entity.Entity) error {
	if incrementId {
		id, err := s.increment([]byte{in.GetObjectType()})
		if err != nil {
			return err
		}

		if err := setId(types.IdType(id)); err != nil {
			return err
		}
	}

	bytes, err := Codec.Marshal(CodecVersion, in)
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

func (s *Session) modify(in entity.Entity, modifyFunc func()) error {
	keys := getObjectKeys(in)
	modifyFunc()
	bytes, err := Codec.Marshal(CodecVersion, in)
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

func (s *Session) remove(in entity.Entity) error {
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
