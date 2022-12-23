package state

import (
	"context"
	"fmt"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/metalgo/database/versiondb"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow/consensus/snowman"
	log "github.com/inconshreveable/log15"
)

func BuildBlock(vm VM, preferred ids.ID) (snowman.Block, error) {
	parent, err := vm.GetStoredBlock(context.Background(), preferred)

	if err != nil {
		return nil, err
	}

	parentDB, err := parent.onAccept()

	if err != nil {
		log.Error("failed to get parent block DB")
		return nil, err
	}

	mempool := vm.GetMempool()
	block := NewBlock(vm, core.Now(), parent.ID(), parent.Index+1)

	for mempool.Len() > 0 {
		vdb := versiondb.New(parentDB)
		next := mempool.Pop()
		receipt, err := vm.ExecuteTransaction(next, vdb)

		if err != nil {
			log.Error("failed to execute transaction", "id", next.Id, "error", err)
			continue
		}

		block.Transactions = append(block.Transactions, *receipt)
	}

	if len(block.Transactions) == 0 {
		return nil, fmt.Errorf("block has no successful transaction")
	}

	// Calculate hash of this block at the end
	block.Finalize()

	return block, nil
}
