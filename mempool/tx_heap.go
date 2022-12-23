package mempool

import (
	"fmt"

	"github.com/MetalBlockchain/antelopevm/core"
)

type txEntry struct {
	id      core.TransactionIdType
	tx      *core.PackedTransaction
	created core.TimePoint
	index   int
}

// txHeap is used to track pending transactions by [price]
type txHeap struct {
	items  []*txEntry
	lookup map[core.TransactionIdType]*txEntry
}

func newTxHeap(items int) *txHeap {
	return &txHeap{
		items:  make([]*txEntry, 0, items),
		lookup: make(map[core.TransactionIdType]*txEntry, items),
	}
}

func (th txHeap) Len() int { return len(th.items) }

func (th txHeap) Less(i, j int) bool {
	return th.items[i].created < th.items[j].created
}

func (th txHeap) Swap(i, j int) {
	th.items[i], th.items[j] = th.items[j], th.items[i]
	th.items[i].index = i
	th.items[j].index = j
}

func (th *txHeap) Push(x interface{}) {
	entry, ok := x.(*txEntry)
	if !ok {
		panic(fmt.Errorf("unexpected %T, expected *txEntry", x))
	}
	if th.Has(entry.id) {
		return
	}
	th.items = append(th.items, entry)
	th.lookup[entry.id] = entry
}

func (th *txHeap) Pop() interface{} {
	n := len(th.items)
	item := th.items[n-1]
	th.items[n-1] = nil // avoid memory leak
	th.items = th.items[0 : n-1]
	delete(th.lookup, item.id)
	return item
}

func (th *txHeap) Get(id core.TransactionIdType) (*txEntry, bool) {
	entry, ok := th.lookup[id]
	return entry, ok
}

func (th *txHeap) Has(id core.TransactionIdType) bool {
	_, has := th.Get(id)
	return has
}
