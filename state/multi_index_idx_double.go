package state

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/contract"
	"github.com/dgraph-io/badger/v3"
)

func (s *Session) FindIdxDoubleObject(id core.IdType) (*contract.IndexDoubleObject, error) {
	if obj, found := s.indexObjectCache.Get(id); found {
		return obj.(*contract.IndexDoubleObject), nil
	}

	key := getObjectKeyByIndex(&contract.IndexDoubleObject{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &contract.IndexDoubleObject{}

	if _, err := out.UnmarshalMsg(data); err != nil {
		return nil, err
	}

	s.indexObjectCache.Put(id, out)

	return out, nil
}

func (s *Session) FindIdxDoubleObjectBySecondary(tableId core.IdType, secondaryKey float64) (*contract.IndexDoubleObject, error) {
	key := getPartialKey("bySecondary", &contract.IndexDoubleObject{}, tableId, secondaryKey)
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

		return s.FindIdxDoubleObject(core.NewIdType(key))
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) FindIdxDoubleObjectByPrimary(tableId core.IdType, primaryKey uint64) (*contract.IndexDoubleObject, error) {
	key := getObjectKeyByIndex(&contract.IndexDoubleObject{TableID: tableId, PrimaryKey: primaryKey}, "byPrimary")
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

		return s.FindIdxDoubleObject(core.NewIdType(key))
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) CreateIdxDoubleObject(in *contract.IndexDoubleObject) error {
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

func (s *Session) ModifyIndexDoubleObject(in *contract.IndexDoubleObject, modifyFunc func()) error {
	if err := s.modify(in, modifyFunc); err != nil {
		return err
	}

	s.indexObjectCache.Put(in.ID, in)

	return nil
}

func (s *Session) RemoveIndexDoubleObject(in *contract.IndexDoubleObject) error {
	if err := s.remove(in); err != nil {
		return err
	}

	s.indexObjectCache.Evict(in.ID)

	return nil
}

func (s *Session) LowerboundSecondaryIndexDouble(tableId core.IdType, secondary float64) (*contract.IndexDoubleObject, error) {
	requiredPrefix := getPartialKey("bySecondary", &contract.IndexDoubleObject{}, tableId)
	prefix := getPartialKey("bySecondary", &contract.IndexDoubleObject{}, tableId, secondary)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*contract.IndexDoubleObject, error) {
		return s.FindIdxDoubleObject(core.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(prefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) UpperboundSecondaryIndexDouble(values ...interface{}) (*contract.IndexDoubleObject, error) {
	requiredPrefix := getPartialKey("bySecondary", &contract.IndexDoubleObject{}, values[:len(values)-1]...)
	prefix := getPartialKey("bySecondary", &contract.IndexDoubleObject{}, values...)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*contract.IndexDoubleObject, error) {
		return s.FindIdxDoubleObject(core.NewIdType(b))
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

func (s *Session) NextSecondaryIndexDouble(kv *contract.IndexDoubleObject) (*contract.IndexDoubleObject, error) {
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

	return s.FindIdxDoubleObject(core.NewIdType(data))
}

func (s *Session) PreviousSecondaryIndexDouble(kv *contract.IndexDoubleObject) (*contract.IndexDoubleObject, error) {
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

	out, err := s.FindIdxDoubleObject(core.NewIdType(data))

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) LowerboundPrimaryIndexDouble(tableId core.IdType, primary uint64) (*contract.IndexDoubleObject, error) {
	requiredPrefix := getPartialKey("byPrimary", &contract.IndexDoubleObject{}, tableId, primary)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*contract.IndexDoubleObject, error) {
		return s.FindIdxDoubleObject(core.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(requiredPrefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) UpperboundPrimaryIndexDouble(values ...interface{}) (*contract.IndexDoubleObject, error) {
	requiredPrefix := getPartialKey("byPrimary", &contract.IndexDoubleObject{}, values...)
	iterator := newReverseIterator(s, requiredPrefix, func(b []byte) (*contract.IndexDoubleObject, error) {
		return s.FindIdxDoubleObject(core.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(requiredPrefix)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) NextPrimaryIndexDouble(kv *contract.IndexDoubleObject) (*contract.IndexDoubleObject, error) {
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

	return s.FindIdxDoubleObject(core.NewIdType(data))
}

func (s *Session) PreviousPrimaryIndexDouble(kv *contract.IndexDoubleObject) (*contract.IndexDoubleObject, error) {
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

	return s.FindIdxDoubleObject(core.NewIdType(data))
}
