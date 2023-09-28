package state

import (
	"github.com/MetalBlockchain/antelopevm/chain/table"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/dgraph-io/badger/v3"
)

func (s *Session) FindIdxDoubleObject(id types.IdType) (*table.IndexDoubleObject, error) {
	if obj, found := s.indexObjectCache.Get(id); found {
		return obj.(*table.IndexDoubleObject), nil
	}

	key := getObjectKeyByIndex(&table.IndexDoubleObject{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &table.IndexDoubleObject{}

	if _, err := Codec.Unmarshal(data, out); err != nil {
		return nil, err
	}

	s.indexObjectCache.Put(id, out)

	return out, nil
}

func (s *Session) FindIdxDoubleObjectBySecondary(tableId types.IdType, secondaryKey float64) (*table.IndexDoubleObject, error) {
	key := getPartialKey("bySecondary", &table.IndexDoubleObject{}, tableId, secondaryKey)
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

		return s.FindIdxDoubleObject(types.NewIdType(key))
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) FindIdxDoubleObjectByPrimary(tableId types.IdType, primaryKey uint64) (*table.IndexDoubleObject, error) {
	key := getObjectKeyByIndex(&table.IndexDoubleObject{TableID: tableId, PrimaryKey: primaryKey}, "byPrimary")
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

		return s.FindIdxDoubleObject(types.NewIdType(key))
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) CreateIdxDoubleObject(in *table.IndexDoubleObject) error {
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

func (s *Session) ModifyIndexDoubleObject(in *table.IndexDoubleObject, modifyFunc func()) error {
	if err := s.modify(in, modifyFunc); err != nil {
		return err
	}

	s.indexObjectCache.Put(in.ID, in)

	return nil
}

func (s *Session) RemoveIndexDoubleObject(in *table.IndexDoubleObject) error {
	if err := s.remove(in); err != nil {
		return err
	}

	s.indexObjectCache.Evict(in.ID)

	return nil
}

func (s *Session) LowerboundSecondaryIndexDouble(tableId types.IdType, secondary float64) (*table.IndexDoubleObject, error) {
	requiredPrefix := getPartialKey("bySecondary", &table.IndexDoubleObject{}, tableId)
	prefix := getPartialKey("bySecondary", &table.IndexDoubleObject{}, tableId, secondary)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*table.IndexDoubleObject, error) {
		return s.FindIdxDoubleObject(types.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(prefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) UpperboundSecondaryIndexDouble(values ...interface{}) (*table.IndexDoubleObject, error) {
	requiredPrefix := getPartialKey("bySecondary", &table.IndexDoubleObject{}, values[:len(values)-1]...)
	prefix := getPartialKey("bySecondary", &table.IndexDoubleObject{}, values...)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*table.IndexDoubleObject, error) {
		return s.FindIdxDoubleObject(types.NewIdType(b))
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

func (s *Session) NextSecondaryIndexDouble(kv *table.IndexDoubleObject) (*table.IndexDoubleObject, error) {
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

	return s.FindIdxDoubleObject(types.NewIdType(data))
}

func (s *Session) PreviousSecondaryIndexDouble(kv *table.IndexDoubleObject) (*table.IndexDoubleObject, error) {
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

	out, err := s.FindIdxDoubleObject(types.NewIdType(data))

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) LowerboundPrimaryIndexDouble(tableId types.IdType, primary uint64) (*table.IndexDoubleObject, error) {
	requiredPrefix := getPartialKey("byPrimary", &table.IndexDoubleObject{}, tableId, primary)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*table.IndexDoubleObject, error) {
		return s.FindIdxDoubleObject(types.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(requiredPrefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) UpperboundPrimaryIndexDouble(values ...interface{}) (*table.IndexDoubleObject, error) {
	requiredPrefix := getPartialKey("byPrimary", &table.IndexDoubleObject{}, values...)
	iterator := newReverseIterator(s, requiredPrefix, func(b []byte) (*table.IndexDoubleObject, error) {
		return s.FindIdxDoubleObject(types.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(requiredPrefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) NextPrimaryIndexDouble(kv *table.IndexDoubleObject) (*table.IndexDoubleObject, error) {
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

	return s.FindIdxDoubleObject(types.NewIdType(data))
}

func (s *Session) PreviousPrimaryIndexDouble(kv *table.IndexDoubleObject) (*table.IndexDoubleObject, error) {
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

	return s.FindIdxDoubleObject(types.NewIdType(data))
}
