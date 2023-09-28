package state

import (
	"github.com/MetalBlockchain/antelopevm/chain/table"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/math"
	"github.com/dgraph-io/badger/v3"
)

func (s *Session) FindIdxLongDoubleObject(id types.IdType) (*table.IndexLongDoubleObject, error) {
	if obj, found := s.indexObjectCache.Get(id); found {
		return obj.(*table.IndexLongDoubleObject), nil
	}

	key := getObjectKeyByIndex(&table.IndexLongDoubleObject{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &table.IndexLongDoubleObject{}
	if _, err := Codec.Unmarshal(data, out); err != nil {
		return nil, err
	}

	s.indexObjectCache.Put(id, out)

	return out, nil
}

func (s *Session) FindIdxLongDoubleObjectBySecondary(tableId types.IdType, secondaryKey math.Float128) (*table.IndexLongDoubleObject, error) {
	key := getPartialKey("bySecondary", &table.IndexLongDoubleObject{}, tableId, secondaryKey)
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

		return s.FindIdxLongDoubleObject(types.NewIdType(key))
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) FindIdxLongDoubleObjectByPrimary(tableId types.IdType, primaryKey uint64) (*table.IndexLongDoubleObject, error) {
	key := getObjectKeyByIndex(&table.IndexLongDoubleObject{TableID: tableId, PrimaryKey: primaryKey}, "byPrimary")
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

		return s.FindIdxLongDoubleObject(types.NewIdType(key))
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) CreateIdxLongDoubleObject(in *table.IndexLongDoubleObject) error {
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

func (s *Session) ModifyIndexLongDoubleObject(in *table.IndexLongDoubleObject, modifyFunc func()) error {
	if err := s.modify(in, modifyFunc); err != nil {
		return err
	}

	s.indexObjectCache.Put(in.ID, in)

	return nil
}

func (s *Session) RemoveIndexLongDoubleObject(in *table.IndexLongDoubleObject) error {
	if err := s.remove(in); err != nil {
		return err
	}

	s.indexObjectCache.Evict(in.ID)

	return nil
}

func (s *Session) LowerboundSecondaryIndexLongDouble(tableId types.IdType, secondary math.Float128) (*table.IndexLongDoubleObject, error) {
	requiredPrefix := getPartialKey("bySecondary", &table.IndexLongDoubleObject{}, tableId)
	prefix := getPartialKey("bySecondary", &table.IndexLongDoubleObject{}, tableId, secondary)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*table.IndexLongDoubleObject, error) {
		return s.FindIdxLongDoubleObject(types.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(prefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) UpperboundSecondaryIndexLongDouble(values ...interface{}) (*table.IndexLongDoubleObject, error) {
	requiredPrefix := getPartialKey("bySecondary", &table.IndexLongDoubleObject{}, values[:len(values)-1]...)
	prefix := getPartialKey("bySecondary", &table.IndexLongDoubleObject{}, values...)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*table.IndexLongDoubleObject, error) {
		return s.FindIdxLongDoubleObject(types.NewIdType(b))
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

func (s *Session) NextSecondaryIndexLongDouble(kv *table.IndexLongDoubleObject) (*table.IndexLongDoubleObject, error) {
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

	return s.FindIdxLongDoubleObject(types.NewIdType(data))
}

func (s *Session) PreviousSecondaryIndexLongDouble(kv *table.IndexLongDoubleObject) (*table.IndexLongDoubleObject, error) {
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

	out, err := s.FindIdxLongDoubleObject(types.NewIdType(data))

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) LowerboundPrimaryIndexLongDouble(tableId types.IdType, primary uint64) (*table.IndexLongDoubleObject, error) {
	requiredPrefix := getPartialKey("byPrimary", &table.IndexLongDoubleObject{}, tableId, primary)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*table.IndexLongDoubleObject, error) {
		return s.FindIdxLongDoubleObject(types.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(requiredPrefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) UpperboundPrimaryIndexLongDouble(values ...interface{}) (*table.IndexLongDoubleObject, error) {
	requiredPrefix := getPartialKey("byPrimary", &table.IndexLongDoubleObject{}, values...)
	iterator := newReverseIterator(s, requiredPrefix, func(b []byte) (*table.IndexLongDoubleObject, error) {
		return s.FindIdxLongDoubleObject(types.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(requiredPrefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) NextPrimaryIndexLongDouble(kv *table.IndexLongDoubleObject) (*table.IndexLongDoubleObject, error) {
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

	return s.FindIdxLongDoubleObject(types.NewIdType(data))
}

func (s *Session) PreviousPrimaryIndexLongDouble(kv *table.IndexLongDoubleObject) (*table.IndexLongDoubleObject, error) {
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

	return s.FindIdxLongDoubleObject(types.NewIdType(data))
}
