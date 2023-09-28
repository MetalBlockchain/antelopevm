package state

import (
	"github.com/MetalBlockchain/antelopevm/chain/account"
	"github.com/MetalBlockchain/antelopevm/chain/authority"
	"github.com/MetalBlockchain/antelopevm/chain/table"
	"github.com/MetalBlockchain/antelopevm/chain/transaction"
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
		c.RegisterType(&authority.Permission{}),
		c.RegisterType(&authority.PermissionLink{}),
		c.RegisterType(&Block{}),
		c.RegisterType(&transaction.TransactionTrace{}),
		c.RegisterType(&table.Table{}),
		c.RegisterType(&table.KeyValue{}),
		Codec.RegisterCodec(CodecVersion, c),
	)

	if errs.Errored() {
		panic(errs.Err)
	}
}
