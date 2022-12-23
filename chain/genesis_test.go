package chain

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeChainId(t *testing.T) {
	genesisFile, _ := os.ReadFile("./genesis_test.json")
	genesis := &GenesisFile{}
	json.Unmarshal(genesisFile, genesis)
	hash := genesis.GetChainId()
	assert.Equal(t, hash.String(), "e771944f7015cfb2fafcb687181d373051ebdbf7396045bd0199cf783b2397e6")
}

func TestInitialKey(t *testing.T) {
	genesisFile, _ := os.ReadFile("./genesis_test.json")
	genesis := &GenesisFile{}
	json.Unmarshal(genesisFile, genesis)
	fmt.Printf("genesis.InitialKey: %v\n", genesis.InitialKey.Valid())
}
