package vm

import (
	"fmt"
	"sync"

	"github.com/MetalBlockchain/metalgo/snow/engine/common"
)

type buildingBlockStatus uint8

const (
	dontBuild buildingBlockStatus = iota
	mayBuild
	building
)

type BlockBuilder struct {
	vm     *VM
	status buildingBlockStatus
	// [l] must be held when accessing [buildStatus]
	l sync.Mutex

	stop        chan struct{}
	builderStop chan struct{}
	doneBuild   chan struct{}
}

func (vm *VM) NewBlockBuilder() BlockBuilder {
	return BlockBuilder{
		vm:          vm,
		status:      dontBuild,
		builderStop: vm.builderStop,
		stop:        vm.stop,
		doneBuild:   vm.doneBuild,
	}
}

func (b *BlockBuilder) Build() {
	defer close(b.doneBuild)

	for {
		select {
		case <-b.vm.mempool.Pending:
			b.signalTxsReady()
		case <-b.builderStop:
			return
		case <-b.stop:
			return
		}
	}
}

func (b *BlockBuilder) signalTxsReady() {
	b.l.Lock()
	defer b.l.Unlock()

	if b.status != dontBuild {
		return
	}

	b.markBuilding()
}

// signal the metalgo engine
// to build a block from pending transactions
func (b *BlockBuilder) markBuilding() {
	select {
	case b.vm.toEngine <- common.PendingTxs:
		b.status = building
	default:
		fmt.Println("dropping message to consensus engine")
	}
}
