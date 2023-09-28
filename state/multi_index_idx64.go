package state

import (
	"github.com/MetalBlockchain/antelopevm/chain/table"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/dgraph-io/badger/v3"
)

func (s *Session) FindIdx64Object(id types.IdType) (*table.Index64Object, error) {
	if obj, found := s.indexObjectCache.Get(id); found {
		return obj.(*table.Index64Object), nil
	}

	key := getObjectKeyByIndex(&table.Index64Object{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &table.Index64Object{}
	if _, err := Codec.Unmarshal(data, out); err != nil {
		return nil, err
	}

	s.indexObjectCache.Put(id, out)

	return out, nil
}

func (s *Session) FindIdx64ObjectBySecondary(tableId types.IdType, secondaryKey uint64) (*table.Index64Object, error) {
	key := getPartialKey("bySecondary", &table.Index64Object{}, tableId, secondaryKey)
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

		return s.FindIdx64Object(types.NewIdType(key))
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) FindIdx64ObjectByPrimary(tableId types.IdType, primaryKey uint64) (*table.Index64Object, error) {
	key := getObjectKeyByIndex(&table.Index64Object{TableID: tableId, PrimaryKey: primaryKey}, "byPrimary")
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

		return s.FindIdx64Object(types.NewIdType(key))
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) CreateIdx64Object(in *table.Index64Object) error {
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

func (s *Session) ModifyIndex64Object(in *table.Index64Object, modifyFunc func()) error {
	if err := s.modify(in, modifyFunc); err != nil {
		return err
	}

	s.indexObjectCache.Put(in.ID, in)

	return nil
}

func (s *Session) RemoveIndex64Object(in *table.Index64Object) error {
	if err := s.remove(in); err != nil {
		return err
	}

	s.indexObjectCache.Evict(in.ID)

	return nil
}

func (s *Session) LowerboundSecondaryIndex64(tableId types.IdType, secondary uint64) (*table.Index64Object, error) {
	requiredPrefix := getPartialKey("bySecondary", &table.Index64Object{}, tableId)
	prefix := getPartialKey("bySecondary", &table.Index64Object{}, tableId, secondary)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*table.Index64Object, error) {
		return s.FindIdx64Object(types.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(prefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) UpperboundSecondaryIndex64(values ...interface{}) (*table.Index64Object, error) {
	requiredPrefix := getPartialKey("bySecondary", &table.Index64Object{}, values[:len(values)-1]...)
	prefix := getPartialKey("bySecondary", &table.Index64Object{}, values...)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*table.Index64Object, error) {
		return s.FindIdx64Object(types.NewIdType(b))
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

func (s *Session) NextSecondaryIndex64(kv *table.Index64Object) (*table.Index64Object, error) {
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

	return s.FindIdx64Object(types.NewIdType(data))
}

func (s *Session) PreviousSecondaryIndex64(kv *table.Index64Object) (*table.Index64Object, error) {
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

	out, err := s.FindIdx64Object(types.NewIdType(data))

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) LowerboundPrimaryIndex64(tableId types.IdType, primary uint64) (*table.Index64Object, error) {
	requiredPrefix := getPartialKey("byPrimary", &table.Index64Object{}, tableId, primary)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*table.Index64Object, error) {
		return s.FindIdx64Object(types.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(requiredPrefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) UpperboundPrimaryIndex64(values ...interface{}) (*table.Index64Object, error) {
	requiredPrefix := getPartialKey("byPrimary", &table.Index64Object{}, values...)
	iterator := newReverseIterator(s, requiredPrefix, func(b []byte) (*table.Index64Object, error) {
		return s.FindIdx64Object(types.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(requiredPrefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) NextPrimaryIndex64(kv *table.Index64Object) (*table.Index64Object, error) {
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

	return s.FindIdx64Object(types.NewIdType(data))
}

func (s *Session) PreviousPrimaryIndex64(kv *table.Index64Object) (*table.Index64Object, error) {
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

	return s.FindIdx64Object(types.NewIdType(data))
}
