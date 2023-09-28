package state

import (
	"context"
	"fmt"
	"time"

	"github.com/MetalBlockchain/antelopevm/chain/block"
	"github.com/MetalBlockchain/antelopevm/chain/entity"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	chainTime "github.com/MetalBlockchain/antelopevm/chain/time"
	"github.com/MetalBlockchain/antelopevm/chain/transaction"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow/choices"
	"github.com/MetalBlockchain/metalgo/snow/consensus/snowman"
	log "github.com/inconshreveable/log15"
)

var (
	_ snowman.Block = &Block{}
	_ entity.Entity = &Block{}
)

type Block struct {
	Index           types.IdType                     `serialize:"true"`
	Hash            block.BlockHash                  `serialize:"true"`
	Header          block.BlockHeader                `serialize:"true"`
	Transactions    []transaction.TransactionReceipt `serialize:"true"`
	BlockExtensions []types.Extension                `serialize:"true"`
	BlockStatus     block.BlockStatus                `serialize:"true"`
	vm              VM
}

func NewBlock(vm VM, timestamp chainTime.TimePoint, parent block.BlockHash, height uint64) *Block {
	return &Block{
		Header: block.BlockHeader{
			Timestamp: block.NewBlockTimeStampFromTimePoint(timestamp),
			Producer:  name.StringToName("eosio"),
			Confirmed: 1,
			Previous:  *crypto.Hash256(parent),
		},
		Transactions: make([]transaction.TransactionReceipt, 0),
		vm:           vm,
		BlockStatus:  block.BlockStatusProcessing,
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
	bytes, err := Codec.Marshal(CodecVersion, b)

	if err != nil {
		log.Error("failed to get bytes", "error", err)
	}

	return bytes
}

// ID returns the ID of this block
func (b *Block) ID() ids.ID {
	return ids.ID(b.Hash)
}

// ParentID returns [b]'s parent's ID
func (b *Block) Parent() ids.ID {
	return ids.ID(b.Header.Previous.Bytes())
}

// Height returns this block's height. The genesis block has height 1.
func (b *Block) Height() uint64 { return uint64(b.Header.BlockNum()) }

// Timestamp returns this block's time. The genesis block has time 0.
func (b *Block) Timestamp() time.Time { return b.Header.Timestamp.ToTimePoint().ToTime() }

// Status returns the status of this block
func (b *Block) Status() choices.Status {
	if b.BlockStatus == block.BlockStatusAccepted {
		return choices.Accepted
	} else if b.BlockStatus == block.BlockStatusRejected {
		return choices.Rejected
	}

	return choices.Processing
}

// SetStatus sets the status of this block
func (b *Block) SetStatus(status choices.Status) {
	if status == choices.Accepted {
		b.BlockStatus = block.BlockStatusAccepted
	} else if status == choices.Rejected {
		b.BlockStatus = block.BlockStatusRejected
	} else {
		b.BlockStatus = block.BlockStatusProcessing
	}
}

func (b *Block) Finalize() {
	digest := crypto.Hash256(b.Header)
	b.Hash = block.BlockHash(digest.FixedBytes())
}

func (b Block) GetId() []byte {
	return b.Index.ToBytes()
}

func (b Block) GetIndexes() map[string]entity.EntityIndex {
	return map[string]entity.EntityIndex{
		"id": {
			Fields: []string{"Index"},
		},
		"byHash": {
			Fields: []string{"Hash"},
		},
	}
}

func (a Block) GetObjectType() uint8 {
	return entity.BlockType
}
