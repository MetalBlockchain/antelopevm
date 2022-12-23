package state

import (
	"context"
	"fmt"
	"time"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/metalgo/database"
	"github.com/MetalBlockchain/metalgo/database/versiondb"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow/choices"
	"github.com/MetalBlockchain/metalgo/snow/consensus/snowman"
	"github.com/ethereum/go-ethereum/log"
)

var (
	_ snowman.Block = &Block{}
)

type Block struct {
	core.BlockHeader `serialize:"true"`
	Transactions     []core.TransactionReceipt `serialize:"true"`
	BlockExtensions  []core.Extension          `serialize:"true"`

	status choices.Status
	hash   [32]byte

	vm         VM
	onAcceptDB *versiondb.Database
}

func NewBlock(vm VM, timestamp core.TimePoint, parentID ids.ID, height uint64) *Block {
	return &Block{
		BlockHeader: core.BlockHeader{
			Created:       timestamp,
			Producer:      core.StringToName("eosio"),
			Confirmed:     1,
			PreviousBlock: parentID,
			Index:         height,
		},
		Transactions: []core.TransactionReceipt{},
		vm:           vm,
		status:       choices.Processing,
	}
}

func (b *Block) Initialize(vm VM, status choices.Status) {
	b.vm = vm
	b.status = status

	copy(b.hash[:], crypto.Hash256(*b).Bytes())
}

// Verify returns nil iff this block is valid.
// To be valid, it must be that:
// b.parent.Timestamp < b.Timestamp <= [local time] + 1 hour
func (b *Block) Verify(ctx context.Context) error {
	log.Debug("verifying block", "block", b)
	parentBlock, err := b.vm.GetStoredBlock(context.Background(), b.Parent())

	if err != nil {
		return err
	}

	parentDB, err := parentBlock.onAccept()

	if err != nil {
		return err
	}

	b.onAcceptDB = versiondb.New(parentDB)

	for _, trx := range b.Transactions {
		if _, err := b.vm.ExecuteTransaction(&trx.Transaction, b.onAcceptDB); err != nil {
			return fmt.Errorf("block contains transaction that failed")
		}
	}

	return b.vm.Verified(b)
}

// Accept sets this block's status to Accepted and sets lastAccepted to this
// block's ID and saves this info to b.vm.DB
func (b *Block) Accept(ctx context.Context) error {
	log.Debug("accepting block", "block", b)
	// Commit changes included in this block
	if b.Index > 0 {
		if err := b.onAcceptDB.Commit(); err != nil {
			return err
		}
	}

	return b.vm.Accepted(b)
}

// Reject sets this block's status to Rejected and saves the status in state
// Recall that b.vm.DB.Commit() must be called to persist to the DB
func (b *Block) Reject(ctx context.Context) error {
	log.Warn("rejected block", "block", b)
	return b.vm.Rejected(b)
}

func (b *Block) Bytes() []byte {
	bytes, _ := Codec.Marshal(uint16(CodecVersion), b)

	return bytes
}

// ID returns the ID of this block
func (b *Block) ID() ids.ID {
	return ids.ID(b.hash)
}

// ParentID returns [b]'s parent's ID
func (b *Block) Parent() ids.ID { return b.PreviousBlock }

// Height returns this block's height. The genesis block has height 1.
func (b *Block) Height() uint64 { return b.Index }

// Timestamp returns this block's time. The genesis block has time 0.
func (b *Block) Timestamp() time.Time { return b.Created.ToTime() }

// Status returns the status of this block
func (b *Block) Status() choices.Status { return b.status }

// SetStatus sets the status of this block
func (b *Block) SetStatus(status choices.Status) { b.status = status }

func (b *Block) onAccept() (database.Database, error) {
	if b.status == choices.Accepted || b.Index == 0 /* genesis */ {
		return b.vm.State(), nil
	}

	if b.onAcceptDB != nil {
		return b.onAcceptDB, nil
	}

	return nil, fmt.Errorf("block not verified")
}

func (b *Block) Finalize() {
	copy(b.hash[:], crypto.Hash256(*b).Bytes())
}
