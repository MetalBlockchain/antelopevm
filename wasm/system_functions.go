package wasm

import (
	"github.com/MetalBlockchain/antelopevm/core"
)

const (
	maxAssertMessage = 1024
)

func GetSystemFunctions(context *ExecutionContext) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["current_time"] = currentTime(context)
	functions["is_feature_activated"] = isFeatureActivated(context)
	functions["get_sender"] = getSender(context)

	return functions
}

func currentTime(context *ExecutionContext) func() uint64 {
	return func() uint64 {
		currentTime := core.Now().TimeSinceEpoch().Count()

		return uint64(currentTime)
	}
}

func isFeatureActivated(context *ExecutionContext) func(uint32) uint32 {
	return func(ptr uint32) uint32 {
		return 0
	}
}

func getSender(context *ExecutionContext) func() uint64 {
	return func() uint64 {
		if sender, err := context.applyContext.GetSender(); err == nil {
			if sender == nil {
				return 0
			}

			return uint64(*sender)
		}

		return 0
	}
}
