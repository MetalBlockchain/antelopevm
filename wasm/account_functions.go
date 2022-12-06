package wasm

import (
	"fmt"

	"github.com/MetalBlockchain/antelopevm/utils"
)

func GetAccountFunctions(context *ExecutionContext) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["require_auth"] = requireAuth(context)
	functions["current_receiver"] = currentReceiver(context)
	functions["is_account"] = isAccount(context)
	functions["require_recipient"] = requireRecipient(context)
	functions["has_auth"] = hasAuth(context)

	return functions
}

func requireAuth(context *ExecutionContext) func(uint64) {
	return func(arg uint64) {
		fmt.Printf("RequireAuth %v\n", utils.NameToString(arg))
	}
}

func currentReceiver(context *ExecutionContext) func() uint64 {
	return func() uint64 {
		fmt.Printf("current_receiver\n")
		return uint64(1)
	}
}

func isAccount(context *ExecutionContext) func(uint64) uint32 {
	return func(arg uint64) uint32 {
		fmt.Printf("IsAccount %v\n", arg)
		return 1
	}
}

func requireRecipient(context *ExecutionContext) func(uint64) {
	return func(arg uint64) {
		fmt.Printf("RequireRecipient %v\n", arg)
	}
}

func hasAuth(context *ExecutionContext) func(uint64) uint32 {
	return func(arg uint64) uint32 {
		fmt.Printf("HasAuth %v\n", arg)
		return 1
	}
}
