package state

import (
	"bytes"
	"math"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/dgraph-io/badger/v3"
)

func (s *Session) FindTable(id core.IdType) (*core.Table, error) {
	if obj, found := s.tableCache.Get(id); found {
		return obj, nil
	}

	key := getObjectKeyByIndex(&core.Table{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &core.Table{}

	if _, err := out.UnmarshalMsg(data); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) FindTableByCodeScopeTable(code name.AccountName, scope name.ScopeName, table name.TableName) (*core.Table, error) {
	key := getObjectKeyByIndex(&core.Table{Code: code, Scope: scope, Table: table}, "byCodeScopeTable")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return s.FindTable(core.NewIdType(data))
}

func (s *Session) FindOrCreateTable(code name.AccountName, scope name.ScopeName, tableName name.TableName, payer name.AccountName) (*core.Table, error) {
	table, err := s.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		if err == badger.ErrKeyNotFound {
			table := &core.Table{
				Code:  code,
				Scope: scope,
				Table: tableName,
				Payer: payer,
				Count: 0,
			}

			if err := s.CreateTable(table); err != nil {
				return nil, err
			}

			return table, nil
		}

		return nil, err
	}

	return table, nil
}

func (s *Session) CreateTable(in *core.Table) error {
	return s.create(true, func(id core.IdType) error {
		in.ID = id
		return nil
	}, in)
}

func (s *Session) ModifyTable(in *core.Table, modifyFunc func()) error {
	return s.modify(in, modifyFunc)
}

func (s *Session) RemoveTable(in *core.Table) error {
	return s.remove(in)
}

func (s *Session) FindKeyValue(id core.IdType) (*core.KeyValue, error) {
	if obj, found := s.kvCache.Get(id); found {
		return obj, nil
	}

	key := getObjectKeyByIndex(&core.KeyValue{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &core.KeyValue{}

	if _, err := out.UnmarshalMsg(data); err != nil {
		return nil, err
	}

	s.kvCache.Put(id, out)

	return out, nil
}

func (s *Session) FindKeyValueByScopePrimary(tableId core.IdType, primaryKey uint64) (*core.KeyValue, error) {
	key := getObjectKeyByIndex(&core.KeyValue{TableID: tableId, PrimaryKey: primaryKey}, "byScopePrimary")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return s.FindKeyValue(core.NewIdType(data))
}

func (s *Session) FindKeyValuesByScope(tableId core.IdType) *Iterator[core.KeyValue] {
	key := getPartialKey("byScopePrimary", &core.KeyValue{}, tableId)
	opts := badger.DefaultIteratorOptions
	opts.Prefix = key

	return newIterator(s, key, func(b []byte) (*core.KeyValue, error) {
		return s.FindKeyValue(core.NewIdType(b))
	})
}

func (s *Session) CreateKeyValue(in *core.KeyValue) error {
	err := s.create(true, func(id core.IdType) error {
		in.ID = id
		return nil
	}, in)

	if err != nil {
		return err
	}

	s.kvCache.Put(in.ID, in)

	return nil
}

func (s *Session) ModifyKeyValue(in *core.KeyValue, modifyFunc func()) error {
	if err := s.modify(in, modifyFunc); err != nil {
		return err
	}

	s.kvCache.Put(in.ID, in)

	return nil
}

func (s *Session) RemoveKeyValue(in *core.KeyValue) error {
	if err := s.remove(in); err != nil {
		return err
	}

	s.kvCache.Evict(in.ID)

	return nil
}

func (s *Session) FindNextKeyValue(kv *core.KeyValue) (*core.KeyValue, error) {
	key := getObjectKeyByIndex(&core.KeyValue{TableID: kv.TableID, PrimaryKey: kv.PrimaryKey}, "byScopePrimary")
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

	return s.FindKeyValue(core.NewIdType(data))
}

func (s *Session) FindPreviousKeyValue(kv *core.KeyValue) (*core.KeyValue, error) {
	key := getObjectKeyByIndex(&core.KeyValue{TableID: kv.TableID, PrimaryKey: kv.PrimaryKey}, "byScopePrimary")
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

	return s.FindKeyValue(core.NewIdType(data))
}

// Find KV with primary key greater than or equal to the provided primary key
func (s *Session) LowerboundKeyValueByScopePrimary(table *core.Table, primaryKey uint64) (*core.KeyValue, error) {
	// See if we have the object
	if kv, err := s.FindKeyValueByScopePrimary(table.ID, primaryKey); err == nil {
		return kv, nil
	} else if err != badger.ErrKeyNotFound {
		return nil, err
	}

	key := getObjectKeyByIndex(&core.KeyValue{TableID: table.ID, PrimaryKey: primaryKey}, "byScopePrimary")
	requiredPrefix := getPartialKey("byScopePrimary", &core.KeyValue{}, table.ID)
	iterator := newIterator(s, requiredPrefix, func(b []byte) (*core.KeyValue, error) {
		return s.FindKeyValue(core.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(key)

	if iterator.ValidForPrefix(requiredPrefix) {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

// Find KV with primary key less than the provided primary key
func (s *Session) UpperboundKeyValueByScope(table *core.Table) (*core.KeyValue, error) {
	key := getObjectKeyByIndex(&core.KeyValue{TableID: table.ID, PrimaryKey: math.MaxUint64}, "byScopePrimary")
	requiredPrefix := getPartialKey("byScopePrimary", &core.KeyValue{}, table.ID)
	iterator := newReverseIterator(s, requiredPrefix, func(b []byte) (*core.KeyValue, error) {
		return s.FindKeyValue(core.NewIdType(b))
	})
	defer iterator.Close()
	iterator.Seek(key)

	if iterator.Valid() {
		return iterator.Item()
	}

	return nil, badger.ErrKeyNotFound
}

func (s *Session) UpperboundKeyValueByScopePrimary(table *core.Table, primaryKey uint64) (*core.KeyValue, error) {
	key := getObjectKeyByIndex(&core.KeyValue{TableID: table.ID, PrimaryKey: primaryKey}, "byScopePrimary")
	requiredPrefix := getPartialKey("byScopePrimary", &core.KeyValue{}, table.ID)
	iterator := newReverseIterator(s, requiredPrefix, func(b []byte) (*core.KeyValue, error) {
		return s.FindKeyValue(core.NewIdType(b))
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
