package api

import (
	"github.com/MetalBlockchain/antelopevm/utils"
)

func init() {
	Functions["get_context_free_data"] = getContextFreeData
}

func GetContextFreeFunctions(context Context) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["get_context_free_data"] = getContextFreeData(context)

	return functions
}

func getContextFreeData(context Context) interface{} {
	return func(index uint32, ptr uint32, length uint32) int32 {
		trx, err := context.GetApplyContext().GetPackedTransaction().GetSignedTransaction()

		if err != nil {
			panic("could not get signed transaction")
		}

		if int(index) >= len(trx.ContextFreeData) {
			return -1
		}

		size := trx.ContextFreeData[index].Size()

		if length == 0 {
			return int32(size)
		}

		copySize := utils.MinInt(int(length), size)
		context.WriteMemory(ptr, trx.ContextFreeData[index][:copySize])

		return int32(copySize)
	}
}
