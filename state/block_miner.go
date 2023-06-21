package state

import (
	"context"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow/consensus/snowman"
	log "github.com/inconshreveable/log15"
)

func BuildBlock(vm VM, preferred ids.ID) (snowman.Block, error) {
	parent, err := vm.GetStoredBlock(context.Background(), preferred)

	if err != nil {
		return nil, err
	}

	session := vm.State().CreateSession(true)
	defer session.Discard()
	mempool := vm.GetMempool()
	block := NewBlock(vm, core.Now(), parent.Hash, parent.Header.Index+1)

	for mempool.Len() > 0 {
		next := mempool.Pop()
		receipt, err := vm.ExecuteTransaction(next, block, session)

		if err != nil {
			log.Error("failed to execute transaction", "id", next.Id, "error", err)
			continue
		}

		block.Transactions = append(block.Transactions, receipt.Receipt)
	}

	// Calculate hash of this block at the end
	block.Finalize()

	return block, nil
}
