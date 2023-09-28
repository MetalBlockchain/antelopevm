package state

import (
	"bytes"
	"math"

	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/table"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/dgraph-io/badger/v3"
)

func (s *Session) FindTable(id types.IdType) (*table.Table, error) {
	if obj, found := s.tableCache.Get(id); found {
		return obj, nil
	}

	key := getObjectKeyByIndex(&table.Table{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &table.Table{}
	if _, err := Codec.Unmarshal(data, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) FindTableByCodeScopeTable(code name.AccountName, scope name.ScopeName, tableName name.TableName) (*table.Table, error) {
	key := getObjectKeyByIndex(&table.Table{Code: code, Scope: scope, Table: tableName}, "byCodeScopeTable")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return s.FindTable(types.NewIdType(data))
}

func (s *Session) FindOrCreateTable(code name.AccountName, scope name.ScopeName, tableName name.TableName, payer name.AccountName) (*table.Table, error) {
	tab, err := s.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		if err == badger.ErrKeyNotFound {
			tab := &table.Table{
				Code:  code,
				Scope: scope,
				Table: tableName,
				Payer: payer,
				Count: 0,
			}

			if err := s.CreateTable(tab); err != nil {
				return nil, err
			}

			return tab, nil
		}

		return nil, err
	}

	return tab, nil
}

func (s *Session) CreateTable(in *table.Table) error {
	return s.create(true, func(id types.IdType) error {
		in.ID = id
		return nil
	}, in)
}

func (s *Session) ModifyTable(in *table.Table, modifyFunc func()) error {
	return s.modify(in, modifyFunc)
}

func (s *Session) RemoveTable(in *table.Table) error {
	return s.remove(in)
}

func (s *Session) FindKeyValue(id types.IdType) (*table.KeyValue, error) {
	if obj, found := s.kvCache.Get(id); found {
		return obj, nil
	}

	key := getObjectKeyByIndex(&table.KeyValue{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &table.KeyValue{}
	if _, err := Codec.Unmarshal(data, out); err != nil {
		return nil, err
	}

	s.kvCache.Put(id, out)

	return out, nil
}

func (s *Session) FindKeyValueByScopePrimary(tableId types.IdType, primaryKey uint64) (*table.KeyValue, error) {
	key := getObjectKeyByIndex(&table.KeyValue{TableID: tableId, PrimaryKey: primaryKey}, "byScopePrimary")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return s.FindKeyValue(types.NewIdType(data))
}

func (s *Session) FindKeyValuesByScope(tableId types.IdType) *Iterator[table.KeyValue] {
	key := getPartialKey("byScopePrimary", &table.KeyValue{}, tableId)
	opts := badger.DefaultIteratorOptions
	opts.Prefix = key

	return newIterator(s, key, func(b []byte) (*table.KeyValue, error) {
		return s.FindKeyValue(types.NewIdType(b))
	})
}

func (s *Session) CreateKeyValue(in *table.KeyValue) error {
	err := s.create(true, func(id types.IdType) error {
		in.ID = id
		return nil
	}, in)

	if err != nil {
		return err
	}

	s.kvCache.Put(in.ID, in)

	return nil
}

func (s *Session) ModifyKeyValue(in *table.KeyValue, modifyFunc func()) error {
	if err := s.modify(in, modifyFunc); err != nil {
		return err
	}

	s.kvCache.Put(in.ID, in)

	return nil
}

func (s *Session) RemoveKeyValue(in *table.KeyValue) error {
	if err := s.remove(in); err != nil {
		return err
	}

	s.kvCache.Evict(in.ID)

	return nil
}

func (s *Session) FindNextKeyValue(kv *table.KeyValue) (*table.KeyValue, error) {
	key := getObjectKeyByIndex(&table.KeyValue{TableID: kv.TableID, PrimaryKey: kv.PrimaryKey}, "byScopePrimary")
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

	return s.FindKeyValue(types.NewIdType(data))
}

func (s *Session) FindPreviousKeyValue(kv *table.KeyValue) (*table.KeyValue, error) {
	key := getObjectKeyByIndex(&table.KeyValue{TableID: kv.TableID, PrimaryKey: kv.PrimaryKey}, "byScopePrimary")
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

	return s.FindKeyValue(types.NewIdType(data))
}

// Find KV with primary key greater than or equal to the provided primary key
func (s *Session) LowerboundKeyValueByScopePrimary(tab *table.Table, primaryKey uint64) (*table.KeyValue, error) {
	// See if we have the object
	if kv, err := s.FindKeyValueByScopePrimary(tab.ID, primaryKey); err == nil {
		return kv, nil
	} else if err != badger.ErrKeyNotFound {
		return nil, err
	}

	key := getObjectKeyByIndex(&table.KeyValue{TableID: tab.ID, PrimaryKey: primaryKey}, "byScopePrimary")
	requiredPrefix := getPartialKey("byScopePrimary", &table.KeyValue{}, tab.ID)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*table.KeyValue, error) {
		return s.FindKeyValue(types.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(key)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

// Find KV with primary key less than the provided primary key
func (s *Session) UpperboundKeyValueByScope(tab *table.Table) (*table.KeyValue, error) {
	key := getObjectKeyByIndex(&table.KeyValue{TableID: tab.ID, PrimaryKey: math.MaxUint64}, "byScopePrimary")
	requiredPrefix := getPartialKey("byScopePrimary", &table.KeyValue{}, tab.ID)
	iterator := newReverseIterator(s, requiredPrefix, func(b []byte) (*table.KeyValue, error) {
		return s.FindKeyValue(types.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(key)

	if iterator.Valid() {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) UpperboundKeyValueByScopePrimary(tab *table.Table, primaryKey uint64) (*table.KeyValue, error) {
	key := getObjectKeyByIndex(&table.KeyValue{TableID: tab.ID, PrimaryKey: primaryKey}, "byScopePrimary")
	requiredPrefix := getPartialKey("byScopePrimary", &table.KeyValue{}, tab.ID)
	iterator := newReverseIterator(s, requiredPrefix, func(b []byte) (*table.KeyValue, error) {
		return s.FindKeyValue(types.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(key)

	// We don't want the KV with the same primary key
	if bytes.Equal(key, iterator.iterator.Item().Key()) {
		iterator.Next()
	}

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}
