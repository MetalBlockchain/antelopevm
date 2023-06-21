package resource

import "github.com/MetalBlockchain/antelopevm/core"

//go:generate msgp
type AccountResourceLimit struct {
	Used                int64               ///< quantity used in current window
	Available           int64               ///< quantity available in current window (based upon fractional reserve)
	Max                 int64               ///< max per window under current congestion
	LastUsageUpdateTime core.BlockTimeStamp ///< last usage timestamp
	CurrentUsed         int64
}
