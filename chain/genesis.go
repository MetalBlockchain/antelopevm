package chain

import (
	"encoding/json"

	"github.com/MetalBlockchain/antelopevm/chain/time"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/config"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
)

type GenesisState struct {
	InitialTimeStamp     time.TimePoint     `json:"initial_timestamp"`
	InitialKey           ecc.PublicKey      `json:"initial_key"`
	InitialConfiguration config.ChainConfig `json:"initial_configuration"`
}

func (g *GenesisState) ComputeChainId() (*types.ChainIdType, error) {
	return crypto.Hash256(g), nil
}

func ParseGenesisData(data []byte) (*GenesisState, error) {
	genesis := &GenesisState{}

	if err := json.Unmarshal(data, genesis); err != nil {
		return nil, err
	}

	return genesis, nil
}
