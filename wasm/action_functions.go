package wasm

func GetActionFunctions(context *ExecutionContext) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["read_action_data"] = readActionData(context)
	functions["action_data_size"] = actionDataSize(context)

	return functions
}

func readActionData(context *ExecutionContext) func(uint32, uint32) uint32 {
	return func(msg uint32, size uint32) uint32 {
		dataSize := len(actionData)

		if size == 0 {
			return uint32(dataSize)
		}

		memory := context.module.Memory()
		dest, _ := memory.Read(context.context, msg, size)
		memcpy(dest, actionData, size)

		return size
	}
}

func actionDataSize(context *ExecutionContext) func() uint32 {
	return func() uint32 {
		return uint32(len(actionData))
	}
}
