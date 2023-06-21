package chain

import (
	"fmt"

	"github.com/MetalBlockchain/antelopevm/core/contract"
	"github.com/MetalBlockchain/antelopevm/core/name"
	wasmApi "github.com/MetalBlockchain/antelopevm/wasm/api"
	"github.com/inconshreveable/log15"
)

var _ wasmApi.MultiIndex[float64] = &IdxDouble{}

type IdxDouble struct {
	Context *applyContext
}

func (i *IdxDouble) Store(scope name.ScopeName, tableName name.TableName, payer name.AccountName, primaryKey uint64, secondaryKey float64) (int, error) {
	if payer.IsEmpty() {
		return -1, errInvalidTablePayer
	}

	table, err := i.Context.FindOrCreateTable(i.Context.Receiver, scope, tableName, payer)

	if err != nil {
		return -1, err
	}

	obj := &contract.IndexDoubleObject{
		TableID:      table.ID,
		PrimaryKey:   primaryKey,
		Payer:        payer,
		SecondaryKey: secondaryKey,
	}

	if err := i.Context.Session.CreateIdxDoubleObject(obj); err != nil {
		return -1, err
	}

	if err := i.Context.Session.ModifyTable(table, func() {
		table.Count++
	}); err != nil {
		return -1, err
	}

	i.Context.KeyValueCache.cacheTable(table)

	iterator := i.Context.KeyValueCache.add(obj)

	return iterator, nil
}

func (i *IdxDouble) Remove(iterator int) error {
	obj, ok := i.Context.KeyValueCache.get(iterator).(*contract.IndexDoubleObject)

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

	if err := i.Context.Session.RemoveIndexDoubleObject(obj); err != nil {
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

func (i *IdxDouble) Update(iterator int, payer name.AccountName, secondaryKey float64) error {
	obj, ok := i.Context.KeyValueCache.get(iterator).(*contract.IndexDoubleObject)

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

	return i.Context.Session.ModifyIndexDoubleObject(obj, func() {
		obj.Payer = payer
		obj.SecondaryKey = secondaryKey
	})
}

func (i *IdxDouble) FindSecondary(code name.AccountName, scope name.ScopeName, tableName name.TableName, secondaryKey *float64, primaryKey *uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)

	if obj, err := i.Context.Session.FindIdxDoubleObjectBySecondary(table.ID, *secondaryKey); err == nil {
		*primaryKey = obj.PrimaryKey

		return i.Context.KeyValueCache.add(obj)
	}

	return endIterator
}

func (i *IdxDouble) LowerboundSecondary(code name.AccountName, scope name.ScopeName, tableName name.TableName, secondaryKey *float64, primaryKey *uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.LowerboundSecondaryIndexDouble(table.ID, *secondaryKey)

	if err != nil {
		return endIterator
	}

	*primaryKey = obj.PrimaryKey
	*secondaryKey = obj.SecondaryKey

	return i.Context.KeyValueCache.add(obj)
}

func (i *IdxDouble) UpperboundSecondary(code name.AccountName, scope name.ScopeName, tableName name.TableName, secondaryKey *float64, primaryKey *uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.UpperboundSecondaryIndexDouble(table.ID, *secondaryKey)

	if err != nil {
		return endIterator
	}

	*primaryKey = obj.PrimaryKey
	*secondaryKey = obj.SecondaryKey

	return i.Context.KeyValueCache.add(obj)
}

func (i *IdxDouble) EndSecondary(code name.AccountName, scope name.ScopeName, tableName name.TableName) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	return i.Context.KeyValueCache.cacheTable(table)
}

func (i *IdxDouble) NextSecondary(iterator int, primaryKey *uint64) (int, error) {
	if iterator < -1 {
		return -1, nil // cannot increment past end iterator of index
	}

	obj := i.Context.KeyValueCache.get(iterator).(*contract.IndexDoubleObject)
	nextObj, err := i.Context.Session.NextSecondaryIndexDouble(obj)

	if err != nil || nextObj.TableID != obj.TableID {
		return i.Context.KeyValueCache.getEndIteratorByTableId(obj.TableID)
	}

	*primaryKey = nextObj.PrimaryKey

	return i.Context.KeyValueCache.add(nextObj), nil
}

func (i *IdxDouble) PreviousSecondary(iterator int, primaryKey *uint64) (int, error) {
	if iterator < -1 {
		table, err := i.Context.KeyValueCache.findTableByEndIterator(iterator)

		if err != nil {
			return -1, fmt.Errorf("not a valid end iterator")
		}

		obj, err := i.Context.Session.UpperboundSecondaryIndexDouble(table.ID)

		if err != nil || obj.TableID != table.ID {
			return -1, nil // Empty index
		}

		*primaryKey = obj.PrimaryKey

		return i.Context.KeyValueCache.add(obj), nil
	}

	obj := i.Context.KeyValueCache.get(iterator).(*contract.IndexDoubleObject)
	previousObj, err := i.Context.Session.PreviousSecondaryIndexDouble(obj)

	log15.Info("previous", "g", previousObj, "obj", obj)

	if err != nil || previousObj.TableID != obj.TableID || previousObj.SecondaryKey > obj.SecondaryKey {
		return -1, nil
	}

	*primaryKey = previousObj.PrimaryKey

	return i.Context.KeyValueCache.add(previousObj), nil
}

func (i *IdxDouble) FindPrimary(code name.AccountName, scope name.ScopeName, tableName name.TableName, secondaryKey *float64, primaryKey uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.FindIdxDoubleObjectByPrimary(table.ID, primaryKey)

	if err != nil {
		return endIterator
	}

	*secondaryKey = obj.SecondaryKey

	return i.Context.KeyValueCache.add(obj)
}

func (i *IdxDouble) LowerboundPrimary(code name.AccountName, scope name.ScopeName, tableName name.TableName, primaryKey uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.LowerboundPrimaryIndexDouble(table.ID, primaryKey)

	if err != nil || obj.TableID != table.ID {
		return endIterator
	}

	return i.Context.KeyValueCache.add(obj)
}

func (i *IdxDouble) UpperboundPrimary(code name.AccountName, scope name.ScopeName, tableName name.TableName, primaryKey uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.UpperboundPrimaryIndexDouble(table.ID, primaryKey)

	if err != nil || obj.TableID != table.ID {
		return endIterator
	}

	return i.Context.KeyValueCache.add(obj)
}

func (i *IdxDouble) NextPrimary(iterator int, primaryKey *uint64) (int, error) {
	if iterator < -1 {
		return -1, nil // cannot increment past end iterator of index
	}

	obj := i.Context.KeyValueCache.get(iterator).(*contract.IndexDoubleObject)
	nextObj, err := i.Context.Session.NextPrimaryIndexDouble(obj)

	if err != nil || nextObj.TableID != obj.TableID {
		return i.Context.KeyValueCache.getEndIteratorByTableId(obj.TableID)
	}

	*primaryKey = nextObj.PrimaryKey

	return i.Context.KeyValueCache.add(nextObj), nil
}

func (i *IdxDouble) PreviousPrimary(iterator int, primaryKey *uint64) (int, error) {
	if iterator < -1 {
		table, err := i.Context.KeyValueCache.findTableByEndIterator(iterator)

		if err != nil {
			return -1, fmt.Errorf("not a valid end iterator")
		}

		obj, err := i.Context.Session.UpperboundPrimaryIndexDouble(table.ID)

		if err != nil || obj.TableID != table.ID {
			return -1, nil // Empty index
		}

		*primaryKey = obj.PrimaryKey

		return i.Context.KeyValueCache.add(obj), nil
	}

	obj := i.Context.KeyValueCache.get(iterator).(*contract.IndexDoubleObject)
	previousObj, err := i.Context.Session.PreviousPrimaryIndexDouble(obj)

	if err != nil || previousObj.TableID != obj.TableID {
		return -1, nil
	}

	*primaryKey = previousObj.PrimaryKey

	return i.Context.KeyValueCache.add(previousObj), nil
}
