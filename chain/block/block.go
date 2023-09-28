package block

import (
	"encoding/hex"
)

type BlockHash [32]byte

func (b BlockHash) Hex() string {
	return hex.EncodeToString(b[:])
}

type BlockStatus uint8

const (
	BlockStatusProcessing BlockStatus = 0
	BlockStatusAccepted   BlockStatus = 1
	BlockStatusRejected   BlockStatus = 2
)
