package service

import (
	"context"

	"github.com/MetalBlockchain/antelopevm/chain"
	"github.com/MetalBlockchain/antelopevm/mempool"
	"github.com/MetalBlockchain/antelopevm/state"
	"github.com/MetalBlockchain/metalgo/ids"
)

type VM interface {
	GetState() *state.State
	GetController() *chain.Controller
	GetMempool() *mempool.Mempool
	LastAccepted(ctx context.Context) (ids.ID, error)
}
