package chain

import (
	"fmt"

	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/table"
	wasmApi "github.com/MetalBlockchain/antelopevm/wasm/api"
)

var (
	_                    wasmApi.MultiIndex[uint64] = &Idx64{}
	errInvalidTablePayer                            = fmt.Errorf("must specify a valid account to pay for new record")
)

type Idx64 struct {
	Context *applyContext
}

func (i *Idx64) Store(scope name.ScopeName, tableName name.TableName, payer name.AccountName, primaryKey uint64, secondaryKey uint64) (int, error) {
	if payer.IsEmpty() {
		return -1, errInvalidTablePayer
	}

	tab, err := i.Context.FindOrCreateTable(i.Context.Receiver, scope, tableName, payer)

	if err != nil {
		return -1, err
	}

	obj := &table.Index64Object{
		TableID:      tab.ID,
		PrimaryKey:   primaryKey,
		Payer:        payer,
		SecondaryKey: secondaryKey,
	}

	if err := i.Context.Session.CreateIdx64Object(obj); err != nil {
		return -1, err
	}

	if err := i.Context.Session.ModifyTable(tab, func() {
		tab.Count++
	}); err != nil {
		return -1, err
	}

	i.Context.KeyValueCache.cacheTable(tab)

	iterator := i.Context.KeyValueCache.add(obj)

	return iterator, nil
}

func (i *Idx64) Remove(iterator int) error {
	obj, ok := i.Context.KeyValueCache.get(iterator).(*table.Index64Object)

	if !ok {
		return fmt.Errorf("could not cast value to")
	}

	table, err := i.Context.KeyValueCache.getTable(obj.TableID)

	if err != nil {
		return err
	}

	if err := i.Context.Session.ModifyTable(table, func() {
		table.Count--
	}); err != nil {
		return err
	}

	if err := i.Context.Session.RemoveIndex64Object(obj); err != nil {
		return err
	}

	if table.Count == 0 {
		if err := i.Context.Session.RemoveTable(table); err != nil {
			return err
		}
	}

	i.Context.KeyValueCache.remove(iterator)

	return nil
}

func (i *Idx64) Update(iterator int, payer name.AccountName, secondaryKey uint64) error {
	obj, ok := i.Context.KeyValueCache.get(iterator).(*table.Index64Object)

	if !ok {
		return fmt.Errorf("could not cast value to")
	}

	table, err := i.Context.KeyValueCache.getTable(obj.TableID)

	if err != nil {
		return err
	}

	if table.Code != i.Context.Receiver {
		return fmt.Errorf("db access violation")
	}

	if payer.IsEmpty() {
		payer = obj.Payer
	}

	if obj.Payer != payer {
		// Update billing size
	}

	return i.Context.Session.ModifyIndex64Object(obj, func() {
		obj.Payer = payer
		obj.SecondaryKey = secondaryKey
	})
}

func (i *Idx64) FindSecondary(code name.AccountName, scope name.ScopeName, tableName name.TableName, secondaryKey *uint64, primaryKey *uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)

	if obj, err := i.Context.Session.FindIdx64ObjectBySecondary(table.ID, *secondaryKey); err == nil {
		*primaryKey = obj.PrimaryKey

		return i.Context.KeyValueCache.add(obj)
	}

	return endIterator
}

func (i *Idx64) LowerboundSecondary(code name.AccountName, scope name.ScopeName, tableName name.TableName, secondaryKey *uint64, primaryKey *uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.LowerboundSecondaryIndex64(table.ID, *secondaryKey)

	if err != nil {
		return endIterator
	}

	*primaryKey = obj.PrimaryKey
	*secondaryKey = obj.SecondaryKey

	return i.Context.KeyValueCache.add(obj)
}

func (i *Idx64) UpperboundSecondary(code name.AccountName, scope name.ScopeName, tableName name.TableName, secondaryKey *uint64, primaryKey *uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.UpperboundSecondaryIndex64(table.ID, *secondaryKey)

	if err != nil {
		return endIterator
	}

	*primaryKey = obj.PrimaryKey
	*secondaryKey = obj.SecondaryKey

	return i.Context.KeyValueCache.add(obj)
}

func (i *Idx64) EndSecondary(code name.AccountName, scope name.ScopeName, tableName name.TableName) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	return i.Context.KeyValueCache.cacheTable(table)
}

func (i *Idx64) NextSecondary(iterator int, primaryKey *uint64) (int, error) {
	if iterator < -1 {
		return -1, nil // cannot increment past end iterator of index
	}

	obj := i.Context.KeyValueCache.get(iterator).(*table.Index64Object)
	nextObj, err := i.Context.Session.NextSecondaryIndex64(obj)

	if err != nil || nextObj.TableID != obj.TableID {
		return i.Context.KeyValueCache.getEndIteratorByTableId(obj.TableID)
	}

	*primaryKey = nextObj.PrimaryKey

	return i.Context.KeyValueCache.add(nextObj), nil
}

func (i *Idx64) PreviousSecondary(iterator int, primaryKey *uint64) (int, error) {
	if iterator < -1 {
		table, err := i.Context.KeyValueCache.findTableByEndIterator(iterator)

		if err != nil {
			return -1, fmt.Errorf("not a valid end iterator")
		}

		obj, err := i.Context.Session.UpperboundSecondaryIndex64(table.ID)

		if err != nil || obj.TableID != table.ID {
			return -1, nil // Empty index
		}

		*primaryKey = obj.PrimaryKey

		return i.Context.KeyValueCache.add(obj), nil
	}

	obj := i.Context.KeyValueCache.get(iterator).(*table.Index64Object)
	previousObj, err := i.Context.Session.PreviousSecondaryIndex64(obj)

	if err != nil || previousObj.TableID != obj.TableID || previousObj.SecondaryKey > obj.SecondaryKey {
		return -1, nil
	}

	*primaryKey = previousObj.PrimaryKey

	return i.Context.KeyValueCache.add(previousObj), nil
}

func (i *Idx64) FindPrimary(code name.AccountName, scope name.ScopeName, tableName name.TableName, secondaryKey *uint64, primaryKey uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.FindIdx64ObjectByPrimary(table.ID, primaryKey)

	if err != nil {
		return endIterator
	}

	*secondaryKey = obj.SecondaryKey

	return i.Context.KeyValueCache.add(obj)
}

func (i *Idx64) LowerboundPrimary(code name.AccountName, scope name.ScopeName, tableName name.TableName, primaryKey uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.LowerboundPrimaryIndex64(table.ID, primaryKey)

	if err != nil || obj.TableID != table.ID {
		return endIterator
	}

	return i.Context.KeyValueCache.add(obj)
}

func (i *Idx64) UpperboundPrimary(code name.AccountName, scope name.ScopeName, tableName name.TableName, primaryKey uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.UpperboundPrimaryIndex64(table.ID, primaryKey)

	if err != nil || obj.TableID != table.ID {
		return endIterator
	}

	return i.Context.KeyValueCache.add(obj)
}

func (i *Idx64) NextPrimary(iterator int, primaryKey *uint64) (int, error) {
	if iterator < -1 {
		return -1, nil // cannot increment past end iterator of index
	}

	obj := i.Context.KeyValueCache.get(iterator).(*table.Index64Object)
	nextObj, err := i.Context.Session.NextPrimaryIndex64(obj)

	if err != nil || nextObj.TableID != obj.TableID {
		return i.Context.KeyValueCache.getEndIteratorByTableId(obj.TableID)
	}

	*primaryKey = nextObj.PrimaryKey

	return i.Context.KeyValueCache.add(nextObj), nil
}

func (i *Idx64) PreviousPrimary(iterator int, primaryKey *uint64) (int, error) {
	if iterator < -1 {
		table, err := i.Context.KeyValueCache.findTableByEndIterator(iterator)

		if err != nil {
			return -1, fmt.Errorf("not a valid end iterator")
		}

		obj, err := i.Context.Session.UpperboundPrimaryIndex64(table.ID)

		if err != nil || obj.TableID != table.ID {
			return -1, nil // Empty index
		}

		*primaryKey = obj.PrimaryKey

		return i.Context.KeyValueCache.add(obj), nil
	}

	obj := i.Context.KeyValueCache.get(iterator).(*table.Index64Object)
	previousObj, err := i.Context.Session.PreviousPrimaryIndex64(obj)

	if err != nil || previousObj.TableID != obj.TableID {
		return -1, nil
	}

	*primaryKey = previousObj.PrimaryKey

	return i.Context.KeyValueCache.add(previousObj), nil
}
