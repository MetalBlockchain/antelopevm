package core

import "github.com/MetalBlockchain/metalgo/codec"

type State interface {
	GetCodec() codec.Manager
	GetCodecVersion() int
}
