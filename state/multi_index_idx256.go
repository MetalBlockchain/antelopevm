package state

import (
	"github.com/MetalBlockchain/antelopevm/chain/table"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/math"
	"github.com/dgraph-io/badger/v3"
)

func (s *Session) FindIdx256Object(id types.IdType) (*table.Index256Object, error) {
	if obj, found := s.indexObjectCache.Get(id); found {
		return obj.(*table.Index256Object), nil
	}

	key := getObjectKeyByIndex(&table.Index256Object{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &table.Index256Object{}
	if _, err := Codec.Unmarshal(data, out); err != nil {
		return nil, err
	}

	s.indexObjectCache.Put(id, out)

	return out, nil
}

func (s *Session) FindIdx256ObjectBySecondary(tableId types.IdType, secondaryKey math.Uint256) (*table.Index256Object, error) {
	key := getPartialKey("bySecondary", &table.Index256Object{}, tableId, secondaryKey)
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

		return s.FindIdx256Object(types.NewIdType(key))
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) FindIdx256ObjectByPrimary(tableId types.IdType, primaryKey uint64) (*table.Index256Object, error) {
	key := getObjectKeyByIndex(&table.Index256Object{TableID: tableId, PrimaryKey: primaryKey}, "byPrimary")
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

		return s.FindIdx256Object(types.NewIdType(key))
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) CreateIdx256Object(in *table.Index256Object) error {
	err := s.create(true, func(id types.IdType) error {
		in.ID = id
		return nil
	}, in)

	if err != nil {
		return err
	}

	s.indexObjectCache.Put(in.ID, in)

	return nil
}

func (s *Session) ModifyIndex256Object(in *table.Index256Object, modifyFunc func()) error {
	if err := s.modify(in, modifyFunc); err != nil {
		return err
	}

	s.indexObjectCache.Put(in.ID, in)

	return nil
}

func (s *Session) RemoveIndex256Object(in *table.Index256Object) error {
	if err := s.remove(in); err != nil {
		return err
	}

	s.indexObjectCache.Evict(in.ID)

	return nil
}

func (s *Session) LowerboundSecondaryIndex256(tableId types.IdType, secondary math.Uint256) (*table.Index256Object, error) {
	requiredPrefix := getPartialKey("bySecondary", &table.Index256Object{}, tableId)
	prefix := getPartialKey("bySecondary", &table.Index256Object{}, tableId, secondary)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*table.Index256Object, error) {
		return s.FindIdx256Object(types.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(prefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) UpperboundSecondaryIndex256(values ...interface{}) (*table.Index256Object, error) {
	requiredPrefix := getPartialKey("bySecondary", &table.Index256Object{}, values[:len(values)-1]...)
	prefix := getPartialKey("bySecondary", &table.Index256Object{}, values...)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*table.Index256Object, error) {
		return s.FindIdx256Object(types.NewIdType(b))
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

func (s *Session) NextSecondaryIndex256(kv *table.Index256Object) (*table.Index256Object, error) {
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

	return s.FindIdx256Object(types.NewIdType(data))
}

func (s *Session) PreviousSecondaryIndex256(kv *table.Index256Object) (*table.Index256Object, error) {
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

	out, err := s.FindIdx256Object(types.NewIdType(data))

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) LowerboundPrimaryIndex256(tableId types.IdType, primary uint64) (*table.Index256Object, error) {
	requiredPrefix := getPartialKey("byPrimary", &table.Index256Object{}, tableId, primary)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*table.Index256Object, error) {
		return s.FindIdx256Object(types.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(requiredPrefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) UpperboundPrimaryIndex256(values ...interface{}) (*table.Index256Object, error) {
	requiredPrefix := getPartialKey("byPrimary", &table.Index256Object{}, values...)
	iterator := newReverseIterator(s, requiredPrefix, func(b []byte) (*table.Index256Object, error) {
		return s.FindIdx256Object(types.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(requiredPrefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) NextPrimaryIndex256(kv *table.Index256Object) (*table.Index256Object, error) {
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

	return s.FindIdx256Object(types.NewIdType(data))
}

func (s *Session) PreviousPrimaryIndex256(kv *table.Index256Object) (*table.Index256Object, error) {
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

	return s.FindIdx256Object(types.NewIdType(data))
}
