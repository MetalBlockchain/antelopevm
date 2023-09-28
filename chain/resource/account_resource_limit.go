package resource

import "github.com/MetalBlockchain/antelopevm/chain/block"

type AccountResourceLimit struct {
	Used                int64                ///< quantity used in current window
	Available           int64                ///< quantity available in current window (based upon fractional reserve)
	Max                 int64                ///< max per window under current congestion
	LastUsageUpdateTime block.BlockTimeStamp ///< last usage timestamp
	CurrentUsed         int64
}
