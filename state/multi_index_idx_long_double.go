package state

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/contract"
	"github.com/MetalBlockchain/antelopevm/math"
	"github.com/dgraph-io/badger/v3"
)

func (s *Session) FindIdxLongDoubleObject(id core.IdType) (*contract.IndexLongDoubleObject, error) {
	if obj, found := s.indexObjectCache.Get(id); found {
		return obj.(*contract.IndexLongDoubleObject), nil
	}

	key := getObjectKeyByIndex(&contract.IndexLongDoubleObject{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &contract.IndexLongDoubleObject{}

	if _, err := out.UnmarshalMsg(data); err != nil {
		return nil, err
	}

	s.indexObjectCache.Put(id, out)

	return out, nil
}

func (s *Session) FindIdxLongDoubleObjectBySecondary(tableId core.IdType, secondaryKey math.Float128) (*contract.IndexLongDoubleObject, error) {
	key := getPartialKey("bySecondary", &contract.IndexLongDoubleObject{}, tableId, secondaryKey)
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

		return s.FindIdxLongDoubleObject(core.NewIdType(key))
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) FindIdxLongDoubleObjectByPrimary(tableId core.IdType, primaryKey uint64) (*contract.IndexLongDoubleObject, error) {
	key := getObjectKeyByIndex(&contract.IndexLongDoubleObject{TableID: tableId, PrimaryKey: primaryKey}, "byPrimary")
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

		return s.FindIdxLongDoubleObject(core.NewIdType(key))
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) CreateIdxLongDoubleObject(in *contract.IndexLongDoubleObject) error {
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

func (s *Session) ModifyIndexLongDoubleObject(in *contract.IndexLongDoubleObject, modifyFunc func()) error {
	if err := s.modify(in, modifyFunc); err != nil {
		return err
	}

	s.indexObjectCache.Put(in.ID, in)

	return nil
}

func (s *Session) RemoveIndexLongDoubleObject(in *contract.IndexLongDoubleObject) error {
	if err := s.remove(in); err != nil {
		return err
	}

	s.indexObjectCache.Evict(in.ID)

	return nil
}

func (s *Session) LowerboundSecondaryIndexLongDouble(tableId core.IdType, secondary math.Float128) (*contract.IndexLongDoubleObject, error) {
	requiredPrefix := getPartialKey("bySecondary", &contract.IndexLongDoubleObject{}, tableId)
	prefix := getPartialKey("bySecondary", &contract.IndexLongDoubleObject{}, tableId, secondary)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*contract.IndexLongDoubleObject, error) {
		return s.FindIdxLongDoubleObject(core.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(prefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) UpperboundSecondaryIndexLongDouble(values ...interface{}) (*contract.IndexLongDoubleObject, error) {
	requiredPrefix := getPartialKey("bySecondary", &contract.IndexLongDoubleObject{}, values[:len(values)-1]...)
	prefix := getPartialKey("bySecondary", &contract.IndexLongDoubleObject{}, values...)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*contract.IndexLongDoubleObject, error) {
		return s.FindIdxLongDoubleObject(core.NewIdType(b))
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

func (s *Session) NextSecondaryIndexLongDouble(kv *contract.IndexLongDoubleObject) (*contract.IndexLongDoubleObject, error) {
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

	return s.FindIdxLongDoubleObject(core.NewIdType(data))
}

func (s *Session) PreviousSecondaryIndexLongDouble(kv *contract.IndexLongDoubleObject) (*contract.IndexLongDoubleObject, error) {
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

	out, err := s.FindIdxLongDoubleObject(core.NewIdType(data))

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) LowerboundPrimaryIndexLongDouble(tableId core.IdType, primary uint64) (*contract.IndexLongDoubleObject, error) {
	requiredPrefix := getPartialKey("byPrimary", &contract.IndexLongDoubleObject{}, tableId, primary)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*contract.IndexLongDoubleObject, error) {
		return s.FindIdxLongDoubleObject(core.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(requiredPrefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) UpperboundPrimaryIndexLongDouble(values ...interface{}) (*contract.IndexLongDoubleObject, error) {
	requiredPrefix := getPartialKey("byPrimary", &contract.IndexLongDoubleObject{}, values...)
	iterator := newReverseIterator(s, requiredPrefix, func(b []byte) (*contract.IndexLongDoubleObject, error) {
		return s.FindIdxLongDoubleObject(core.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(requiredPrefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) NextPrimaryIndexLongDouble(kv *contract.IndexLongDoubleObject) (*contract.IndexLongDoubleObject, error) {
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

	return s.FindIdxLongDoubleObject(core.NewIdType(data))
}

func (s *Session) PreviousPrimaryIndexLongDouble(kv *contract.IndexLongDoubleObject) (*contract.IndexLongDoubleObject, error) {
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

	return s.FindIdxLongDoubleObject(core.NewIdType(data))
}
