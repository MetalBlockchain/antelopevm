package mempool

import (
	"container/heap"
	"sync"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/inconshreveable/log15"
)

type Mempool struct {
	mu      sync.RWMutex
	maxSize int
	heap    *txHeap

	// Pending is a channel of length one, which the mempool ensures has an item on
	// it as long as there is an unissued transaction remaining in [txs]
	Pending chan struct{}
	// newTxs is an array of [Tx] that are ready to be gossiped.
	newTxs []*core.PackedTransaction
}

func New(maxSize int) *Mempool {
	return &Mempool{
		maxSize: maxSize,
		heap:    newTxHeap(maxSize),
		Pending: make(chan struct{}, 1),
	}
}

func (th *Mempool) Add(tx *core.PackedTransaction) bool {
	log15.Info("adding tx to mempool", "tx", tx)
	th.mu.Lock()
	defer th.mu.Unlock()

	unpacked, err := tx.GetUnpackedTransaction()

	if err != nil {
		return false
	}

	// Don't add duplicates
	if th.heap.Has(unpacked.ID()) {
		return false
	}

	oldLen := th.heap.Len()

	heap.Push(th.heap, &txEntry{
		id:      unpacked.ID(),
		created: core.Now(),
		tx:      tx,
		index:   oldLen,
	})

	if th.heap.Len() > th.maxSize {
		th.popMin()
	}

	// When adding [tx] to the mempool make sure that there is an item in Pending
	// to signal the VM to produce a block. Note: if the VM's buildStatus has already
	// been set to something other than [dontBuild], this will be ignored and won't be
	// reset until the engine calls BuildBlock. This case is handled in IssueCurrentTx
	// and CancelCurrentTx.
	th.newTxs = append(th.newTxs, tx)
	th.addPending()

	return true
}

func (th *Mempool) Pop() *core.PackedTransaction { // O(log N)
	th.mu.Lock()
	defer th.mu.Unlock()

	item := th.heap.items[0]
	return th.remove(item.id)
}

func (th *Mempool) Peek() *core.PackedTransaction {
	th.mu.RLock()
	defer th.mu.RUnlock()

	txEntry := th.newTxs[0]

	return txEntry
}

func (th *Mempool) Len() int {
	th.mu.RLock()
	defer th.mu.RUnlock()

	return th.heap.Len()
}

// popMin assumes the write lock is held and takes O(log N) time to run.
func (th *Mempool) popMin() *core.PackedTransaction {
	item := th.heap.items[0]

	return th.remove(item.id)
}

// remove assumes the write lock is held and takes O(log N) time to run.
func (th *Mempool) remove(id core.TransactionIdType) *core.PackedTransaction {
	entry, ok := th.heap.Get(id)

	if !ok {
		return nil
	}

	heap.Remove(th.heap, entry.index)

	return entry.tx
}

func (th *Mempool) addPending() {
	select {
	case th.Pending <- struct{}{}:
	default:
	}
}
