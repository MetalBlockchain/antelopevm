package state

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/contract"
	"github.com/MetalBlockchain/antelopevm/math"
	"github.com/dgraph-io/badger/v3"
)

func (s *Session) FindIdx128Object(id core.IdType) (*contract.Index128Object, error) {
	if obj, found := s.indexObjectCache.Get(id); found {
		return obj.(*contract.Index128Object), nil
	}

	key := getObjectKeyByIndex(&contract.Index128Object{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &contract.Index128Object{}

	if _, err := out.UnmarshalMsg(data); err != nil {
		return nil, err
	}

	s.indexObjectCache.Put(id, out)

	return out, nil
}

func (s *Session) FindIdx128ObjectBySecondary(tableId core.IdType, secondaryKey math.Uint128) (*contract.Index128Object, error) {
	key := getPartialKey("bySecondary", &contract.Index128Object{}, tableId, secondaryKey)
	opts := badger.DefaultIteratorOptions
	opts.Prefix = key
	iterator := s.transaction.NewIterator(opts)
	defer iterator.Close()
	iterator.Rewind()

	if iterator.ValidForPrefix(key) {
		key, err := iterator.Item().ValueCopy(nil)

		if err != nil {
			return nil, err
		}

		return s.FindIdx128Object(core.NewIdType(key))
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) FindIdx128ObjectByPrimary(tableId core.IdType, primaryKey uint64) (*contract.Index128Object, error) {
	key := getObjectKeyByIndex(&contract.Index128Object{TableID: tableId, PrimaryKey: primaryKey}, "byPrimary")
	opts := badger.DefaultIteratorOptions
	opts.Prefix = key
	iterator := s.transaction.NewIterator(opts)
	defer iterator.Close()
	iterator.Rewind()

	if iterator.ValidForPrefix(key) {
		key, err := iterator.Item().ValueCopy(nil)

		if err != nil {
			return nil, err
		}

		return s.FindIdx128Object(core.NewIdType(key))
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) CreateIdx128Object(in *contract.Index128Object) error {
	err := s.create(true, func(id core.IdType) error {
		in.ID = id
		return nil
	}, in)

	if err != nil {
		return err
	}

	s.indexObjectCache.Put(in.ID, in)

	return nil
}

func (s *Session) ModifyIndex128Object(in *contract.Index128Object, modifyFunc func()) error {
	if err := s.modify(in, modifyFunc); err != nil {
		return err
	}

	s.indexObjectCache.Put(in.ID, in)

	return nil
}

func (s *Session) RemoveIndex128Object(in *contract.Index128Object) error {
	if err := s.remove(in); err != nil {
		return err
	}

	s.indexObjectCache.Evict(in.ID)

	return nil
}

func (s *Session) LowerboundSecondaryIndex128(tableId core.IdType, secondary math.Uint128) (*contract.Index128Object, error) {
	requiredPrefix := getPartialKey("bySecondary", &contract.Index128Object{}, tableId)
	prefix := getPartialKey("bySecondary", &contract.Index128Object{}, tableId, secondary)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*contract.Index128Object, error) {
		return s.FindIdx128Object(core.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(prefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) UpperboundSecondaryIndex128(values ...interface{}) (*contract.Index128Object, error) {
	requiredPrefix := getPartialKey("bySecondary", &contract.Index128Object{}, values[:len(values)-1]...)
	prefix := getPartialKey("bySecondary", &contract.Index128Object{}, values...)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*contract.Index128Object, error) {
		return s.FindIdx128Object(core.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(prefix)

	// Should be key > than given key
	if iterator.ValidForPrefix(prefix) {
		iterator.Next()
	}

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) NextSecondaryIndex128(kv *contract.Index128Object) (*contract.Index128Object, error) {
	key := getObjectKeyByIndex(kv, "bySecondary")
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 2
	opts.Prefix = key
	iterator := s.transaction.NewIterator(opts)
	defer iterator.Close()
	iterator.Rewind()
	iterator.Next()
	item := iterator.Item()
	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return s.FindIdx128Object(core.NewIdType(data))
}

func (s *Session) PreviousSecondaryIndex128(kv *contract.Index128Object) (*contract.Index128Object, error) {
	key := getObjectKeyByIndex(kv, "bySecondary")
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 2
	opts.Prefix = key
	opts.Reverse = true
	iterator := s.transaction.NewIterator(opts)
	defer iterator.Close()
	iterator.Rewind()
	iterator.Next()
	item := iterator.Item()
	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out, err := s.FindIdx128Object(core.NewIdType(data))

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) LowerboundPrimaryIndex128(tableId core.IdType, primary uint64) (*contract.Index128Object, error) {
	requiredPrefix := getPartialKey("byPrimary", &contract.Index128Object{}, tableId, primary)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*contract.Index128Object, error) {
		return s.FindIdx128Object(core.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(requiredPrefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) UpperboundPrimaryIndex128(values ...interface{}) (*contract.Index128Object, error) {
	requiredPrefix := getPartialKey("byPrimary", &contract.Index128Object{}, values...)
	iterator := newReverseIterator(s, requiredPrefix, func(b []byte) (*contract.Index128Object, error) {
		return s.FindIdx128Object(core.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(requiredPrefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) NextPrimaryIndex128(kv *contract.Index128Object) (*contract.Index128Object, error) {
	key := getObjectKeyByIndex(kv, "byPrimary")
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 2
	opts.Prefix = key
	iterator := s.transaction.NewIterator(opts)
	defer iterator.Close()
	iterator.Rewind()
	iterator.Next()
	item := iterator.Item()
	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return s.FindIdx128Object(core.NewIdType(data))
}

func (s *Session) PreviousPrimaryIndex128(kv *contract.Index128Object) (*contract.Index128Object, error) {
	key := getObjectKeyByIndex(kv, "byPrimary")
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 2
	opts.Prefix = key
	opts.Reverse = true
	iterator := s.transaction.NewIterator(opts)
	defer iterator.Close()
	iterator.Rewind()
	iterator.Next()
	item := iterator.Item()
	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return s.FindIdx128Object(core.NewIdType(data))
}
