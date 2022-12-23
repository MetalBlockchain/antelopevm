package core

import (
	"github.com/MetalBlockchain/antelopevm/crypto"
)

type Account struct {
	ID             IdType         `serialize:"true"`
	Name           AccountName    `serialize:"true"`
	VmType         uint8          `serialize:"true"`
	VmVersion      uint8          `serialize:"true"`
	Privileged     bool           `serialize:"true"`
	LastCodeUpdate TimePoint      `serialize:"true"`
	CodeVersion    crypto.Sha256  `serialize:"true"`
	CreationDate   BlockTimeStamp `serialize:"true"`
	Code           HexBytes       `serialize:"true"`
	Abi            HexBytes       `serialize:"true"`
	RecvSequence   uint64         `serialize:"true"`
	AuthSequence   uint64         `serialize:"true"`
	CodeSequence   uint64         `serialize:"true"`
	AbiSequence    uint64         `serialize:"true"`
}
