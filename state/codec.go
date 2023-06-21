package state

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/account"
	"github.com/MetalBlockchain/metalgo/codec"
	"github.com/MetalBlockchain/metalgo/codec/hierarchycodec"
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
	c := hierarchycodec.NewDefault()
	Codec = codec.NewDefaultManager()

	errs := wrappers.Errs{}
	errs.Add(
		c.RegisterType(&account.Account{}),
		c.RegisterType(&core.Permission{}),
		c.RegisterType(&core.PermissionLink{}),
		c.RegisterType(&Block{}),
		c.RegisterType(&core.TransactionTrace{}),
		c.RegisterType(&core.Table{}),
		c.RegisterType(&core.KeyValue{}),
		Codec.RegisterCodec(CodecVersion, c),
	)

	if errs.Errored() {
		panic(errs.Err)
	}
}
