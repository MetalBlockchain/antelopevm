package state

import (
	"errors"
	"time"

	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow/choices"
	"github.com/MetalBlockchain/metalgo/snow/consensus/snowman"
	"github.com/MetalBlockchain/metalgo/utils/hashing"
)

var (
	errTimestampTooEarly = errors.New("block's timestamp is earlier than its parent's timestamp")
	errDatabaseGet       = errors.New("error while retrieving data from database")
	errTimestampTooLate  = errors.New("block's timestamp is more than 1 hour ahead of local time")

	_ snowman.Block = &Block{}
)

const (
	dataLen = 32
)

type Block struct {
	PreviousBlock ids.ID `serialize:"true" json:"previous"`
	Hght          uint64 `serialize:"true" json:"block_num"`
	id            ids.ID `serialize:"true" json:"id"`
	Tmstmp        int64
	Dt            [dataLen]byte
	bytes         []byte
	status        choices.Status
}

// Verify returns nil iff this block is valid.
// To be valid, it must be that:
// b.parent.Timestamp < b.Timestamp <= [local time] + 1 hour
func (b *Block) Verify() error {
	return nil
}

// Initialize sets [b.bytes] to [bytes], [b.id] to hash([b.bytes]),
// [b.status] to [status] and [b.vm] to [vm]
func (b *Block) Initialize(bytes []byte, status choices.Status) {
	b.bytes = bytes
	b.id = hashing.ComputeHash256Array(b.bytes)
	b.status = status
}

// Accept sets this block's status to Accepted and sets lastAccepted to this
// block's ID and saves this info to b.vm.DB
func (b *Block) Accept() error {
	return nil
}

// Reject sets this block's status to Rejected and saves the status in state
// Recall that b.vm.DB.Commit() must be called to persist to the DB
func (b *Block) Reject() error {
	return nil
}

// ID returns the ID of this block
func (b *Block) ID() ids.ID { return b.id }

// ParentID returns [b]'s parent's ID
func (b *Block) Parent() ids.ID { return b.PreviousBlock }

// Height returns this block's height. The genesis block has height 0.
func (b *Block) Height() uint64 { return b.Hght }

// Timestamp returns this block's time. The genesis block has time 0.
func (b *Block) Timestamp() time.Time { return time.Unix(b.Tmstmp, 0) }

// Status returns the status of this block
func (b *Block) Status() choices.Status { return b.status }

// Bytes returns the byte repr. of this block
func (b *Block) Bytes() []byte { return b.bytes }

// Data returns the data of this block
func (b *Block) Data() [dataLen]byte { return b.Dt }

// SetStatus sets the status of this block
func (b *Block) SetStatus(status choices.Status) { b.status = status }
