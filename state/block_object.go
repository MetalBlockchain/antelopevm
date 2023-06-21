package state

import (
	"context"
	"fmt"
	"time"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow/choices"
	"github.com/MetalBlockchain/metalgo/snow/consensus/snowman"
	"github.com/ethereum/go-ethereum/log"
	"github.com/inconshreveable/log15"
)

var (
	_ snowman.Block = &Block{}
	_ core.Entity   = &Block{}
)

//go:generate msgp
type Block struct {
	Index           core.IdType               `serialize:"true"`
	Hash            core.BlockHash            `serialize:"true"`
	Header          core.BlockHeader          `serialize:"true"`
	Transactions    []core.TransactionReceipt `serialize:"true"`
	BlockExtensions []core.Extension          `serialize:"true"`
	BlockStatus     core.BlockStatus          `serialize:"true"`
	vm              VM
}

func NewBlock(vm VM, timestamp core.TimePoint, parent core.BlockHash, height uint64) *Block {
	return &Block{
		Header: core.BlockHeader{
			Created:           timestamp,
			Producer:          name.StringToName("eosio"),
			Confirmed:         1,
			PreviousBlockHash: parent,
			Index:             height,
		},
		Transactions: make([]core.TransactionReceipt, 0),
		vm:           vm,
		BlockStatus:  core.BlockStatusProcessing,
	}
}

func (b *Block) Initialize(vm VM) {
	b.vm = vm
}

// Verify returns nil iff this block is valid.
// To be valid, it must be that:
// b.parent.Timestamp < b.Timestamp <= [local time] + 1 hour
func (b *Block) Verify(ctx context.Context) error {
	log.Debug("verifying block", "block", b)

	session := b.vm.State().CreateSession(true)
	defer session.Discard()

	for _, trx := range b.Transactions {
		if trace, err := b.vm.ExecuteTransaction(&trx.Transaction, b, session); err != nil {
			return fmt.Errorf("block contains transaction that failed")
		} else {
			if err := session.CreateTransaction(trace); err != nil {
				return err
			}
		}
	}

	if err := session.Commit(); err != nil {
		return err
	}

	return b.vm.Verified(b)
}

// Accept sets this block's status to Accepted and sets lastAccepted to this
// block's ID and saves this info to b.vm.DB
func (b *Block) Accept(ctx context.Context) error {
	log.Debug("accepting block", "block", b)

	return b.vm.Accepted(b)
}

// Reject sets this block's status to Rejected and saves the status in state
// Recall that b.vm.DB.Commit() must be called to persist to the DB
func (b *Block) Reject(ctx context.Context) error {
	return b.vm.Rejected(b)
}

func (b *Block) Bytes() []byte {
	bytes, err := b.MarshalMsg(nil)

	if err != nil {
		log15.Error("failed to get bytes", "error", err)
	}

	return bytes
}

// ID returns the ID of this block
func (b *Block) ID() ids.ID {
	return ids.ID(b.Hash)
}

// ParentID returns [b]'s parent's ID
func (b *Block) Parent() ids.ID {
	return ids.ID(b.Header.PreviousBlockHash)
}

// Height returns this block's height. The genesis block has height 1.
func (b *Block) Height() uint64 { return b.Header.Index }

// Timestamp returns this block's time. The genesis block has time 0.
func (b *Block) Timestamp() time.Time { return b.Header.Created.ToTime() }

// Status returns the status of this block
func (b *Block) Status() choices.Status {
	if b.BlockStatus == core.BlockStatusAccepted {
		return choices.Accepted
	} else if b.BlockStatus == core.BlockStatusRejected {
		return choices.Rejected
	}

	return choices.Processing
}

// SetStatus sets the status of this block
func (b *Block) SetStatus(status choices.Status) {
	if status == choices.Accepted {
		b.BlockStatus = core.BlockStatusAccepted
	} else if status == choices.Rejected {
		b.BlockStatus = core.BlockStatusRejected
	} else {
		b.BlockStatus = core.BlockStatusProcessing
	}
}

func (b *Block) Finalize() {
	digest := crypto.Hash256(b.Header)
	b.Hash = core.BlockHash(digest.FixedBytes())
}

func (b Block) GetId() []byte {
	return b.Index.ToBytes()
}

func (b Block) GetIndexes() map[string]core.EntityIndex {
	return map[string]core.EntityIndex{
		"id": {
			Fields: []string{"Index"},
		},
		"byHash": {
			Fields: []string{"Hash"},
		},
	}
}

func (a Block) GetObjectType() uint8 {
	return core.BlockType
}
