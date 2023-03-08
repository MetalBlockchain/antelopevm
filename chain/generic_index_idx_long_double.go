package chain

import (
	"fmt"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/contract"
	"github.com/MetalBlockchain/antelopevm/math"
)

type IdxLongDouble struct {
	Context *applyContext
}

func (i *IdxLongDouble) Store(scope core.ScopeName, tableName core.TableName, payer core.AccountName, primaryKey uint64, secondaryKey math.Float128) (int, error) {
	if payer.IsEmpty() {
		return -1, errInvalidTablePayer
	}

	table, err := i.Context.FindOrCreateTable(i.Context.Receiver, scope, tableName, payer)

	if err != nil {
		return -1, err
	}

	obj := &contract.IndexLongDoubleObject{
		TableID:      table.ID,
		PrimaryKey:   primaryKey,
		Payer:        payer,
		SecondaryKey: secondaryKey,
	}

	if err := i.Context.Session.CreateIdxLongDoubleObject(obj); err != nil {
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

func (i *IdxLongDouble) Remove(iterator int) error {
	obj, ok := i.Context.KeyValueCache.get(iterator).(*contract.IndexLongDoubleObject)

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

	if err := i.Context.Session.RemoveIndexLongDoubleObject(obj); err != nil {
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

func (i *IdxLongDouble) Update(iterator int, payer core.AccountName, secondaryKey math.Float128) error {
	obj, ok := i.Context.KeyValueCache.get(iterator).(*contract.IndexLongDoubleObject)

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

	return i.Context.Session.ModifyIndexLongDoubleObject(obj, func() {
		obj.Payer = payer
		obj.SecondaryKey = secondaryKey
	})
}

func (i *IdxLongDouble) FindSecondary(code core.AccountName, scope core.ScopeName, tableName core.TableName, secondaryKey *math.Float128, primaryKey *uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)

	if obj, err := i.Context.Session.FindIdxLongDoubleObjectBySecondary(table.ID, *secondaryKey); err == nil {
		*primaryKey = obj.PrimaryKey

		return i.Context.KeyValueCache.add(obj)
	}

	return endIterator
}

func (i *IdxLongDouble) LowerboundSecondary(code core.AccountName, scope core.ScopeName, tableName core.TableName, secondaryKey *math.Float128, primaryKey *uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.LowerboundSecondaryIndexLongDouble(table.ID, *secondaryKey)

	if err != nil {
		return endIterator
	}

	*primaryKey = obj.PrimaryKey
	*secondaryKey = obj.SecondaryKey

	return i.Context.KeyValueCache.add(obj)
}

func (i *IdxLongDouble) UpperboundSecondary(code core.AccountName, scope core.ScopeName, tableName core.TableName, secondaryKey *math.Float128, primaryKey *uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.UpperboundSecondaryIndexLongDouble(table.ID, *secondaryKey)

	if err != nil {
		return endIterator
	}

	*primaryKey = obj.PrimaryKey
	*secondaryKey = obj.SecondaryKey

	return i.Context.KeyValueCache.add(obj)
}

func (i *IdxLongDouble) EndSecondary(code core.AccountName, scope core.ScopeName, tableName core.TableName) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	return i.Context.KeyValueCache.cacheTable(table)
}

func (i *IdxLongDouble) NextSecondary(iterator int, primaryKey *uint64) (int, error) {
	if iterator < -1 {
		return -1, nil // cannot increment past end iterator of index
	}

	obj := i.Context.KeyValueCache.get(iterator).(*contract.IndexLongDoubleObject)
	nextObj, err := i.Context.Session.NextSecondaryIndexLongDouble(obj)

	if err != nil || nextObj.TableID != obj.TableID {
		return i.Context.KeyValueCache.getEndIteratorByTableId(obj.TableID)
	}

	*primaryKey = nextObj.PrimaryKey

	return i.Context.KeyValueCache.add(nextObj), nil
}

func (i *IdxLongDouble) PreviousSecondary(iterator int, primaryKey *uint64) (int, error) {
	if iterator < -1 {
		table, err := i.Context.KeyValueCache.findTableByEndIterator(iterator)

		if err != nil {
			return -1, fmt.Errorf("not a valid end iterator")
		}

		obj, err := i.Context.Session.UpperboundSecondaryIndexLongDouble(table.ID)

		if err != nil || obj.TableID != table.ID {
			return -1, nil // Empty index
		}

		*primaryKey = obj.PrimaryKey

		return i.Context.KeyValueCache.add(obj), nil
	}

	obj := i.Context.KeyValueCache.get(iterator).(*contract.IndexLongDoubleObject)
	previousObj, err := i.Context.Session.PreviousSecondaryIndexLongDouble(obj)

	if err != nil || previousObj.TableID != obj.TableID || !previousObj.SecondaryKey.F128Lt(obj.SecondaryKey) {
		return -1, nil
	}

	*primaryKey = previousObj.PrimaryKey

	return i.Context.KeyValueCache.add(previousObj), nil
}

func (i *IdxLongDouble) FindPrimary(code core.AccountName, scope core.ScopeName, tableName core.TableName, secondaryKey *math.Float128, primaryKey uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.FindIdxLongDoubleObjectByPrimary(table.ID, primaryKey)

	if err != nil {
		return endIterator
	}

	*secondaryKey = obj.SecondaryKey

	return i.Context.KeyValueCache.add(obj)
}

func (i *IdxLongDouble) LowerboundPrimary(code core.AccountName, scope core.ScopeName, tableName core.TableName, primaryKey uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.LowerboundPrimaryIndexLongDouble(table.ID, primaryKey)

	if err != nil || obj.TableID != table.ID {
		return endIterator
	}

	return i.Context.KeyValueCache.add(obj)
}

func (i *IdxLongDouble) UpperboundPrimary(code core.AccountName, scope core.ScopeName, tableName core.TableName, primaryKey uint64) int {
	table, err := i.Context.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		return -1
	}

	endIterator := i.Context.KeyValueCache.cacheTable(table)
	obj, err := i.Context.Session.UpperboundPrimaryIndexLongDouble(table.ID, primaryKey)

	if err != nil || obj.TableID != table.ID {
		return endIterator
	}

	return i.Context.KeyValueCache.add(obj)
}

func (i *IdxLongDouble) NextPrimary(iterator int, primaryKey *uint64) (int, error) {
	if iterator < -1 {
		return -1, nil // cannot increment past end iterator of index
	}

	obj := i.Context.KeyValueCache.get(iterator).(*contract.IndexLongDoubleObject)
	nextObj, err := i.Context.Session.NextPrimaryIndexLongDouble(obj)

	if err != nil || nextObj.TableID != obj.TableID {
		return i.Context.KeyValueCache.getEndIteratorByTableId(obj.TableID)
	}

	*primaryKey = nextObj.PrimaryKey

	return i.Context.KeyValueCache.add(nextObj), nil
}

func (i *IdxLongDouble) PreviousPrimary(iterator int, primaryKey *uint64) (int, error) {
	if iterator < -1 {
		table, err := i.Context.KeyValueCache.findTableByEndIterator(iterator)

		if err != nil {
			return -1, fmt.Errorf("not a valid end iterator")
		}

		obj, err := i.Context.Session.UpperboundPrimaryIndexLongDouble(table.ID)

		if err != nil || obj.TableID != table.ID {
			return -1, nil // Empty index
		}

		*primaryKey = obj.PrimaryKey

		return i.Context.KeyValueCache.add(obj), nil
	}

	obj := i.Context.KeyValueCache.get(iterator).(*contract.IndexLongDoubleObject)
	previousObj, err := i.Context.Session.PreviousPrimaryIndexLongDouble(obj)

	if err != nil || previousObj.TableID != obj.TableID {
		return -1, nil
	}

	*primaryKey = previousObj.PrimaryKey

	return i.Context.KeyValueCache.add(previousObj), nil
}
