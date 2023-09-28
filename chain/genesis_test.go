package chain_test

import (
	"os"
	"testing"

	"github.com/MetalBlockchain/antelopevm/chain"
	"github.com/stretchr/testify/assert"
)

func TestComputeChainId(t *testing.T) {
	genesisFile, err := os.ReadFile("./genesis_test.json")
	assert.NoError(t, err)
	genesis, err := chain.ParseGenesisData(genesisFile)
	assert.NoError(t, err)
	hash, err := genesis.ComputeChainId()
	assert.NoError(t, err)
	assert.Equal(t, hash.String(), "384da888112027f0321850a169f737c33e53b388aad48b5adace4bab97f437e0") // XPR Network ID
}
