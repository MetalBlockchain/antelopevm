package chain

import (
	"fmt"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/contract"
	"github.com/MetalBlockchain/antelopevm/math"
	"github.com/inconshreveable/log15"
)

type Idx256 struct {
	Context *applyContext
}

func (i *Idx256) Store(scope core.ScopeName, tableName core.TableName, payer core.AccountName, primaryKey uint64, secondaryKey math.Uint256) (int, error) {
	if payer.IsEmpty() {
		return -1, errInvalidTablePayer
	}

	table, err := i.Context.FindOrCreateTable(i.Context.Receiver, scope, tableName, payer)

	if err != nil {
		return -1, err
	}

	obj := &contract.Index256Object{
		TableID:      table.ID,
		PrimaryKey:   primaryKey,
		Payer:        payer,
		SecondaryKey: secondaryKey,
	}

	if err := i.Context.Session.CreateIdx256Object(obj); err != nil {
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

func (i *Idx256) Remove(iterator int) error {
	obj, ok := i.Context.KeyValueCache.get(iterator).(*contract.Index256Object)

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

	if err := i.Context.Session.RemoveIndex256Object(obj); err != nil {
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

func (i *Idx256) Update(iterator int, payer core.AccountName, secondaryKey math.Uint256) error {
	obj, ok := i.Context.KeyValueCache.get(iterator).(*contract.Index256Object)

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

	return i.Context.Session.ModifyIndex256Object(obj, func() {
		obj.Payer = payer
		obj.SecondaryKey = secondaryKey
	})
}

func (i *Idx256) FindSecondary(code core.AccountName, scope core.ScopeName, tableName core.TableName, secondaryKey *math.Uint256, primaryKey *uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)

	if obj, err := i.Context.Session.FindIdx256ObjectBySecondary(table.ID, *secondaryKey); err == nil {
		*primaryKey = obj.PrimaryKey

		return i.Context.KeyValueCache.add(obj)
	}

	return endIterator
}

func (i *Idx256) LowerboundSecondary(code core.AccountName, scope core.ScopeName, tableName core.TableName, secondaryKey *math.Uint256, primaryKey *uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.LowerboundSecondaryIndex256(table.ID, *secondaryKey)

	if err != nil {
		return endIterator
	}

	*primaryKey = obj.PrimaryKey
	*secondaryKey = obj.SecondaryKey

	return i.Context.KeyValueCache.add(obj)
}

func (i *Idx256) UpperboundSecondary(code core.AccountName, scope core.ScopeName, tableName core.TableName, secondaryKey *math.Uint256, primaryKey *uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.UpperboundSecondaryIndex256(table.ID, *secondaryKey)

	if err != nil {
		return endIterator
	}

	*primaryKey = obj.PrimaryKey
	*secondaryKey = obj.SecondaryKey

	return i.Context.KeyValueCache.add(obj)
}

func (i *Idx256) EndSecondary(code core.AccountName, scope core.ScopeName, tableName core.TableName) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	return i.Context.KeyValueCache.cacheTable(table)
}

func (i *Idx256) NextSecondary(iterator int, primaryKey *uint64) (int, error) {
	if iterator < -1 {
		return -1, nil // cannot increment past end iterator of index
	}

	obj := i.Context.KeyValueCache.get(iterator).(*contract.Index256Object)
	nextObj, err := i.Context.Session.NextSecondaryIndex256(obj)

	if err != nil || nextObj.TableID != obj.TableID {
		return i.Context.KeyValueCache.getEndIteratorByTableId(obj.TableID)
	}

	*primaryKey = nextObj.PrimaryKey

	return i.Context.KeyValueCache.add(nextObj), nil
}

func (i *Idx256) PreviousSecondary(iterator int, primaryKey *uint64) (int, error) {
	if iterator < -1 {
		table, err := i.Context.KeyValueCache.findTableByEndIterator(iterator)

		if err != nil {
			return -1, fmt.Errorf("not a valid end iterator")
		}

		obj, err := i.Context.Session.UpperboundSecondaryIndex256(table.ID)

		if err != nil || obj.TableID != table.ID {
			return -1, nil // Empty index
		}

		*primaryKey = obj.PrimaryKey

		return i.Context.KeyValueCache.add(obj), nil
	}

	obj := i.Context.KeyValueCache.get(iterator).(*contract.Index256Object)
	previousObj, err := i.Context.Session.PreviousSecondaryIndex256(obj)

	log15.Info("previous", "g", previousObj, "obj", obj)

	if err != nil || previousObj.TableID != obj.TableID || previousObj.SecondaryKey.Compare(obj.SecondaryKey) > 0 {
		return -1, nil
	}

	*primaryKey = previousObj.PrimaryKey

	return i.Context.KeyValueCache.add(previousObj), nil
}

func (i *Idx256) FindPrimary(code core.AccountName, scope core.ScopeName, tableName core.TableName, secondaryKey *math.Uint256, primaryKey uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.FindIdx256ObjectByPrimary(table.ID, primaryKey)

	if err != nil {
		return endIterator
	}

	*secondaryKey = obj.SecondaryKey

	return i.Context.KeyValueCache.add(obj)
}

func (i *Idx256) LowerboundPrimary(code core.AccountName, scope core.ScopeName, tableName core.TableName, primaryKey uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.LowerboundPrimaryIndex256(table.ID, primaryKey)

	if err != nil || obj.TableID != table.ID {
		return endIterator
	}

	return i.Context.KeyValueCache.add(obj)
}

func (i *Idx256) UpperboundPrimary(code core.AccountName, scope core.ScopeName, tableName core.TableName, primaryKey uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.UpperboundPrimaryIndex256(table.ID, primaryKey)

	if err != nil || obj.TableID != table.ID {
		return endIterator
	}

	return i.Context.KeyValueCache.add(obj)
}

func (i *Idx256) NextPrimary(iterator int, primaryKey *uint64) (int, error) {
	if iterator < -1 {
		return -1, nil // cannot increment past end iterator of index
	}

	obj := i.Context.KeyValueCache.get(iterator).(*contract.Index256Object)
	nextObj, err := i.Context.Session.NextPrimaryIndex256(obj)

	if err != nil || nextObj.TableID != obj.TableID {
		return i.Context.KeyValueCache.getEndIteratorByTableId(obj.TableID)
	}

	*primaryKey = nextObj.PrimaryKey

	return i.Context.KeyValueCache.add(nextObj), nil
}

func (i *Idx256) PreviousPrimary(iterator int, primaryKey *uint64) (int, error) {
	if iterator < -1 {
		table, err := i.Context.KeyValueCache.findTableByEndIterator(iterator)

		if err != nil {
			return -1, fmt.Errorf("not a valid end iterator")
		}

		obj, err := i.Context.Session.UpperboundPrimaryIndex256(table.ID)

		if err != nil || obj.TableID != table.ID {
			return -1, nil // Empty index
		}

		*primaryKey = obj.PrimaryKey

		return i.Context.KeyValueCache.add(obj), nil
	}

	obj := i.Context.KeyValueCache.get(iterator).(*contract.Index256Object)
	previousObj, err := i.Context.Session.PreviousPrimaryIndex256(obj)

	if err != nil || previousObj.TableID != obj.TableID {
		return -1, nil
	}

	*primaryKey = previousObj.PrimaryKey

	return i.Context.KeyValueCache.add(previousObj), nil
}
