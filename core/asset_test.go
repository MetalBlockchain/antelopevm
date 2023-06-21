package core_test

import (
	"testing"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/stretchr/testify/assert"
)

func TestAsset(t *testing.T) {
	asset := core.Asset{
		Amount: 10000,
		Symbol: core.Symbol{
			Precision: 4,
			Symbol:    "XPR",
		},
	}
	assert.Equal(t, asset.String(), "1.0000 XPR")
}
