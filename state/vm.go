package state

import (
	"context"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/mempool"
	"github.com/MetalBlockchain/metalgo/database"
	"github.com/MetalBlockchain/metalgo/database/versiondb"
	"github.com/MetalBlockchain/metalgo/ids"
)

type VM interface {
	Accepted(*Block) error
	Rejected(*Block) error
	Verified(*Block) error
	State() database.Database
	GetStoredBlock(context.Context, ids.ID) (*Block, error)
	GetMempool() *mempool.Mempool
	ExecuteTransaction(*core.PackedTransaction, *versiondb.Database) (*core.TransactionReceipt, error)
}
