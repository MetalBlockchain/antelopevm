package block

import (
	"math"

	"github.com/MetalBlockchain/antelopevm/chain/time"
	"github.com/MetalBlockchain/antelopevm/config"
)

type BlockTimeStamp uint32

func MaxBlockTimeStamp() BlockTimeStamp {
	return 0xffff
}

func MinBlockTimeStamp() BlockTimeStamp {
	return 0
}

func (b BlockTimeStamp) Next() BlockTimeStamp {
	if math.MaxUint32-b >= 1 {
		panic("block timestamp overflow")
	}

	return b + 1
}

func NewBlockTimeStampFromTimePoint(t time.TimePoint) BlockTimeStamp {
	microSinceEpoch := t.TimeSinceEpoch()
	msecSinceEpoch := microSinceEpoch.Count() / 1000
	return BlockTimeStamp((msecSinceEpoch - config.BlockTimestampEpochMs) / config.BlockIntervalMs)
}

func NewBlockTimeStampFromTimePointSec(t time.TimePointSec) BlockTimeStamp {
	secSinceEpoch := t.SecSinceEpoch()
	return BlockTimeStamp((secSinceEpoch*1000 - uint32(config.BlockTimestampEpochMs)) / uint32(config.BlockIntervalMs))
}

func (b BlockTimeStamp) ToTimePoint() time.TimePoint {
	var msec uint64 = uint64(b) * uint64(config.BlockIntervalMs)
	msec += uint64(config.BlockTimestampEpochMs)
	return time.TimePoint(time.Milliseconds(int64(msec)))
}
