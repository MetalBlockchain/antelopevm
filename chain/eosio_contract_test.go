package chain_test

import (
	"fmt"
	"testing"

	"github.com/MetalBlockchain/antelopevm/chain"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	applyContext, err := chain.NewApplyContext(nil, 0, 0)
	assert.NoError(t, err)
	fmt.Printf("applyContext: %v\n", applyContext)
}
