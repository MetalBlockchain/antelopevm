package state

import (
	"errors"
	"time"

	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow/choices"
	"github.com/MetalBlockchain/metalgo/snow/consensus/snowman"
)

var (
	errTimestampTooEarly = errors.New("block's timestamp is earlier than its parent's timestamp")
	errDatabaseGet       = errors.New("error while retrieving data from database")
	errTimestampTooLate  = errors.New("block's timestamp is more than 1 hour ahead of local time")

	_ snowman.Block = &Block{}
)

type BlockHeader struct {
	Created               types.TimePoint   `serialize:"true"`
	Producer              types.AccountName `serialize:"true"`
	Confirmed             uint16            `serialize:"true"`
	PreviousBlock         ids.ID            `serialize:"true"`
	Index                 uint64            `serialize:"true"`
	TransactionMerkleRoot crypto.Sha256     `serialize:"true"`
	ActionMerkleRoot      crypto.Sha256     `serialize:"true"`
	ScheduleVersion       uint32            `serialize:"true"`
	HeaderExtensions      []types.Extension `serialize:"true"`
	status                choices.Status
	hash                  [32]byte
}

type SignedBlockHeader struct {
	BlockHeader       `serialize:"true"`
	ProducerSignature ecc.Signature `serialize:"true"`
}

type Block struct {
	BlockHeader     `serialize:"true"`
	Transactions    []types.TransactionReceipt `serialize:"true"`
	BlockExtensions []types.Extension          `serialize:"true"`

	vm VM
}

// Verify returns nil iff this block is valid.
// To be valid, it must be that:
// b.parent.Timestamp < b.Timestamp <= [local time] + 1 hour
func (b *Block) Verify() error {
	return b.vm.Verified(b)
}

// Initialize sets [b.bytes] to [bytes], [b.id] to hash([b.bytes]),
// [b.status] to [status] and [b.vm] to [vm]
func (b *Block) Initialize(vm VM, status choices.Status) {
	b.vm = vm
	b.status = status

	copy(b.hash[:], crypto.Hash256(*b).Bytes())
}

// Accept sets this block's status to Accepted and sets lastAccepted to this
// block's ID and saves this info to b.vm.DB
func (b *Block) Accept() error {
	return b.vm.Accepted(b)
}

// Reject sets this block's status to Rejected and saves the status in state
// Recall that b.vm.DB.Commit() must be called to persist to the DB
func (b *Block) Reject() error {
	return nil
}

func (b *Block) Bytes() []byte {
	bytes, _ := Codec.Marshal(CodecVersion, b)

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
