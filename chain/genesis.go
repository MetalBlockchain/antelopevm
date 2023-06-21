package chain

import (
	"encoding/json"

	"github.com/MetalBlockchain/antelopevm/config"
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
)

type GenesisFile struct {
	InitialTimeStamp     core.TimePoint     `json:"initial_timestamp"`
	InitialKey           ecc.PublicKey      `json:"initial_key"`
	InitialConfiguration config.ChainConfig `json:"initial_configuration"`
}

func (g *GenesisFile) GetChainId() core.ChainIdType {
	return core.ChainIdType(*crypto.Hash256(g))
}

func (g *GenesisFile) Validate() error {
	return nil
}

func ParseGenesisData(data []byte) *GenesisFile {
	genesis := &GenesisFile{}
	json.Unmarshal(data, genesis)
	return genesis
}
