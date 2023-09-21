package api

import (
	"github.com/MetalBlockchain/antelopevm/core/name"
)

func init() {
	Functions["read_action_data"] = readActionData
	Functions["action_data_size"] = actionDataSize
	Functions["current_receiver"] = currentReceiver
	Functions["set_action_return_value"] = setActionReturnValue
}

func readActionData(context Context) interface{} {
	return func(msg uint32, size uint32) uint32 {
		if size == 0 {
			return uint32(context.GetApplyContext().GetAction().Data.Size())
		}

		context.WriteMemory(msg, context.GetApplyContext().GetAction().Data[0:size])

		return uint32(size)
	}
}

func actionDataSize(context Context) interface{} {
	return func() uint32 {
		return uint32(len(context.GetApplyContext().GetAction().Data))
	}
}

func currentReceiver(context Context) interface{} {
	return func() name.AccountName {
		return context.GetApplyContext().GetReceiver()
	}
}

func setActionReturnValue(context Context) interface{} {
	return func(ptr uint32, length uint32) {
		data := context.ReadMemory(ptr, length)
		context.GetApplyContext().SetActionReturnValue(data)
	}
}
