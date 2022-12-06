package wasm

import (
	"encoding/hex"
	"fmt"

	"github.com/MetalBlockchain/antelopevm/utils"
)

var (
	actionData, _ = hex.DecodeString("000000594e9ae9ad0014be6a4f9ae9ade550f30e00000000045850520000000008446f6e6174696f6e")
)

func FindI64(context *ExecutionContext) func(uint64, uint64, uint64, uint64) uint32 {
	return func(code uint64, scope uint64, table uint64, id uint64) uint32 {
		fmt.Printf("find_i64: %v %v %v %v\n", utils.NameToString(code), utils.NameToString(scope), utils.NameToString(table), utils.NameToString(id))
		return 100
	}
}

func StoreI64(context *ExecutionContext) func(uint64, uint64, uint64, uint64, uint32, uint32) uint32 {
	return func(scope uint64, table uint64, payer uint64, id uint64, arg5 uint32, arg6 uint32) uint32 {
		fmt.Printf("store_i64: %v\n", scope)
		return 1
	}
}

func UpdateI64(context *ExecutionContext) func(uint32, uint64, uint32, uint32) {
	return func(arg1 uint32, arg2 uint64, arg3 uint32, arg4 uint32) {
		fmt.Printf("update_i64: %v\n", arg1)
	}
}

func NextI64(context *ExecutionContext) func(uint32, uint32) uint32 {
	return func(arg1 uint32, arg2 uint32) uint32 {
		fmt.Printf("next_i64: %v\n", arg1)
		return 42
	}
}

func GetI64(context *ExecutionContext) func(uint32, uint32, uint32) uint32 {
	return func(arg1 uint32, arg2 uint32, arg3 uint32) uint32 {
		fmt.Printf("get_i64: %v %v %v\n", arg1, arg2, arg3)
		return 1
	}
}

func RemoveI64(context *ExecutionContext) func(uint32) {
	return func(arg uint32) {
		fmt.Printf("remove_i64: %v\n", arg)
	}
}
