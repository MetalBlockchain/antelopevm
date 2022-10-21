package state

import (
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto"
)

type Account struct {
	ID             types.IdType         `serialize:"true"`
	Name           types.AccountName    `serialize:"true"`
	VmType         uint8                `serialize:"true"`
	VmVersion      uint8                `serialize:"true"`
	Privileged     bool                 `serialize:"true"`
	LastCodeUpdate types.TimePoint      `serialize:"true"`
	CodeVersion    crypto.Sha256        `serialize:"true"`
	CreationDate   types.BlockTimeStamp `serialize:"true"`
	Code           types.HexBytes       `serialize:"true"`
	Abi            types.HexBytes       `serialize:"true"`
	RecvSequence   uint64               `serialize:"true"`
	AuthSequence   uint64               `serialize:"true"`
	CodeSequence   uint64               `serialize:"true"`
	AbiSequence    uint64               `serialize:"true"`
}
