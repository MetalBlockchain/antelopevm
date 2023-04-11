package chain

import (
	"fmt"

	"github.com/MetalBlockchain/antelopevm/core"
)

type pairTableIterator struct {
	table    *core.Table
	iterator int
}

type IteratorCache struct {
	tableCache         map[core.IdType]*pairTableIterator
	endIteratorToTable []*core.Table
	iteratorToObject   []interface{}
	objectToIterator   map[interface{}]int
}

func NewIteratorCache() *IteratorCache {
	i := IteratorCache{
		tableCache:         make(map[core.IdType]*pairTableIterator),
		endIteratorToTable: make([]*core.Table, 4),
		iteratorToObject:   make([]interface{}, 8),
		objectToIterator:   make(map[interface{}]int),
	}

	return &i
}

func (i *IteratorCache) getEndIteratorByTableId(id core.IdType) (int, error) {
	if table, found := i.tableCache[id]; found {
		return table.iterator, nil
	}

	return 0, fmt.Errorf("an invariant was broken, table should be in cache")
}
func (i *IteratorCache) findTableByEndIterator(endIterator int) (*core.Table, error) {
	index := i.endIteratorToIndex(endIterator)

	if index >= len(i.endIteratorToTable) {
		return nil, fmt.Errorf("not a valid end iterator")
	}

	return i.endIteratorToTable[index], nil
}
func (i *IteratorCache) endIteratorToIndex(endIterator int) int { return (-endIterator - 2) }
func (i *IteratorCache) indexToEndIterator(index int) int       { return -(index + 2) }
func (i *IteratorCache) cacheTable(table *core.Table) int {
	if itr, ok := i.tableCache[table.ID]; ok {
		return itr.iterator
	}

	ei := i.indexToEndIterator(len(i.endIteratorToTable))
	i.endIteratorToTable = append(i.endIteratorToTable, table)
	i.tableCache[table.ID] = &pairTableIterator{
		table:    table,
		iterator: ei,
	}

	return ei
}

func (i *IteratorCache) getTable(id core.IdType) (*core.Table, error) {
	if table, ok := i.tableCache[id]; ok {
		return table.table, nil
	}

	return nil, fmt.Errorf("an invariant was broken, table should be in cache")
}

func (i *IteratorCache) add(obj interface{}) int {
	if itr, ok := i.objectToIterator[obj]; ok {
		return itr
	}

	i.iteratorToObject = append(i.iteratorToObject, obj)
	i.objectToIterator[obj] = len(i.iteratorToObject) - 1

	return len(i.iteratorToObject) - 1
}

func (i *IteratorCache) get(iterator int) interface{} {
	if iterator == -1 {
		panic("invalid iterator")
	}

	if iterator < 0 {
		panic("dereference of end iterator")
	}

	if iterator >= len(i.iteratorToObject) {
		panic("iterator out of range")
	}

	obj := i.iteratorToObject[iterator]

	if obj == nil {
		panic("could not find iterator")
	}

	return obj
}

func (i *IteratorCache) remove(iterator int) {
	obj := i.iteratorToObject[iterator]

	if obj == nil {
		return
	}

	i.iteratorToObject = nil
	delete(i.objectToIterator, obj)
}
