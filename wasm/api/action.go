package api

import (
	"github.com/MetalBlockchain/antelopevm/core"
)

func GetActionFunctions(context Context) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["read_action_data"] = readActionData(context)
	functions["action_data_size"] = actionDataSize(context)
	functions["current_receiver"] = currentReceiver(context)
	functions["set_action_return_value"] = setActionReturnValue(context)

	return functions
}

func readActionData(context Context) func(uint32, uint32) uint32 {
	return func(msg uint32, size uint32) uint32 {
		if size == 0 {
			return uint32(context.GetApplyContext().GetAction().Data.Size())
		}

		context.WriteMemory(msg, context.GetApplyContext().GetAction().Data[0:size])

		return uint32(size)
	}
}

func actionDataSize(context Context) func() uint32 {
	return func() uint32 {
		return uint32(len(context.GetApplyContext().GetAction().Data))
	}
}

func currentReceiver(context Context) func() core.AccountName {
	return func() core.AccountName {
		return context.GetApplyContext().GetReceiver()
	}
}

func setActionReturnValue(context Context) func(uint32, uint32) {
	return func(ptr uint32, length uint32) {
		data := context.ReadMemory(ptr, length)
		context.GetApplyContext().SetActionReturnValue(data)
	}
}
