package block_test

import (
	"testing"

	"github.com/MetalBlockchain/antelopevm/chain/block"
	"github.com/MetalBlockchain/antelopevm/chain/time"
	"github.com/stretchr/testify/assert"
)

func TestBlockTimeStampToTimePoint(t *testing.T) {
	blockTimestamp := block.BlockTimeStamp(100)
	tp := blockTimestamp.ToTimePoint()
	assert.Equal(t, "2000-01-01T00:00:50.000", tp.String())
}

func TestBlockTimeStampFromTimePoint(t *testing.T) {
	time, err := time.FromIsoString("2020-04-22T17:00:00.000")
	assert.NoError(t, err)
	blockTimestamp := block.NewBlockTimeStampFromTimePoint(time)
	assert.Equal(t, uint32(1281780000), (uint32)(blockTimestamp))
}
