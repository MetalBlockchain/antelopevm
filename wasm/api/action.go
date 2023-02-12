package api

import log "github.com/inconshreveable/log15"

func GetActionFunctions(context Context) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["read_action_data"] = readActionData(context)
	functions["action_data_size"] = actionDataSize(context)

	return functions
}

func readActionData(context Context) func(uint32, uint32) uint32 {
	return func(msg uint32, size uint32) uint32 {
		log.Info("read_action_data", "msg", msg, "size", size, "data", context.GetApplyContext().GetAction().Data.HexString())

		if size == 0 {
			return uint32(context.GetApplyContext().GetAction().Data.Size())
		}

		context.WriteMemory(msg, context.GetApplyContext().GetAction().Data[0:size])

		return uint32(size)
	}
}

func actionDataSize(context Context) func() int32 {
	return func() int32 {
		log.Info("action_data_size", "length", len(context.GetApplyContext().GetAction().Data))
		return int32(len(context.GetApplyContext().GetAction().Data))
	}
}
