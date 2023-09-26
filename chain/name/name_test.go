package name_test

import (
	"testing"

	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/stretchr/testify/assert"
)

func TestStringToName(t *testing.T) {
	name := name.StringToName("eosio")
	assert.Equal(t, uint64(6138663577826885632), uint64(name))
}

func TestNameToString(t *testing.T) {
	name := name.Name(6138663577826885632)
	assert.Equal(t, "eosio", name.String())
}
