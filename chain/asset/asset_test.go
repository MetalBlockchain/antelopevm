package asset_test

import (
	"testing"

	"github.com/MetalBlockchain/antelopevm/chain/asset"
	"github.com/stretchr/testify/assert"
)

func TestAsset(t *testing.T) {
	asset := asset.Asset{
		Amount: 10000,
		Symbol: asset.Symbol{
			Precision: 4,
			Symbol:    "XPR",
		},
	}
	assert.Equal(t, asset.String(), "1.0000 XPR")
}
