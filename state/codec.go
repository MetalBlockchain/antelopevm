package state

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/metalgo/codec"
	"github.com/MetalBlockchain/metalgo/codec/linearcodec"
	"github.com/MetalBlockchain/metalgo/utils/wrappers"
)

const (
	// CodecVersion is the current default codec version
	CodecVersion = 0
)

// Codecs do serialization and deserialization
var (
	Codec codec.Manager
)

func init() {
	// Create default codec and manager
	c := linearcodec.NewDefault()
	Codec = codec.NewDefaultManager()

	errs := wrappers.Errs{}
	errs.Add(
		c.RegisterType(&core.Account{}),
		c.RegisterType(&core.Permission{}),
		c.RegisterType(&core.PermissionLink{}),
		c.RegisterType(&Block{}),
		c.RegisterType(&core.TransactionTrace{}),
		Codec.RegisterCodec(CodecVersion, c),
	)

	if errs.Errored() {
		panic(errs.Err)
	}
}
