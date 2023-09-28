package transaction

import (
	"github.com/MetalBlockchain/antelopevm/chain/authority"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
)

type Action struct {
	Account       name.AccountName            `serialize:"true" json:"account"`
	Name          name.ActionName             `serialize:"true" json:"name"`
	Authorization []authority.PermissionLevel `serialize:"true" json:"authorization,omitempty"`
	Data          types.HexBytes              `serialize:"true" json:"hex_data"`
	ParsedData    map[string]interface{}      `json:"data,omitempty" eos:"-"`
}

func (a Action) DataAs(t interface{}) error {
	return rlp.DecodeBytes(a.Data, t)
}
