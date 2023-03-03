package api

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/MetalBlockchain/antelopevm/utils"
)

func GetContextFreeTransactionFunctions(context Context) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["read_transaction"] = readTransaction(context)
	functions["transaction_size"] = transactionSize(context)
	functions["expiration"] = expiration(context)
	functions["tapos_block_num"] = taposBlockNum(context)
	functions["tapos_block_prefix"] = taposBlockPrefix(context)
	functions["get_action"] = getAction(context)

	return functions
}

func readTransaction(context Context) func(uint32, uint32) uint32 {
	return func(ptr uint32, length uint32) uint32 {
		trx := context.GetApplyContext().GetPackedTransaction()
		trxBytes, err := rlp.EncodeToBytes(trx)

		if err != nil {
			panic("could not encode transaction to RLP")
		}

		trxSize := len(trxBytes)

		if length == 0 {
			return uint32(trxSize)
		}

		copySize := utils.MinInt(trxSize, int(length))
		context.WriteMemory(ptr, trxBytes[:copySize])

		return uint32(copySize)
	}
}

func transactionSize(context Context) func() uint32 {
	return func() uint32 {
		trx := context.GetApplyContext().GetPackedTransaction()
		trxBytes, err := rlp.EncodeToBytes(trx)

		if err != nil {
			panic("could not encode transaction to RLP")
		}

		return uint32(len(trxBytes))
	}
}

func expiration(context Context) func() uint32 {
	return func() uint32 {
		trx, err := context.GetApplyContext().GetPackedTransaction().GetTransaction()

		if err != nil {
			panic("could not unpack transaction")
		}

		return trx.Expiration.SecSinceEpoch()
	}
}

func taposBlockNum(context Context) func() uint32 {
	return func() uint32 {
		trx, err := context.GetApplyContext().GetPackedTransaction().GetTransaction()

		if err != nil {
			panic("could not unpack transaction")
		}

		return uint32(trx.RefBlockNum)
	}
}

func taposBlockPrefix(context Context) func() uint32 {
	return func() uint32 {
		trx, err := context.GetApplyContext().GetPackedTransaction().GetTransaction()

		if err != nil {
			panic("could not unpack transaction")
		}

		return trx.RefBlockPrefix
	}
}

func getAction(context Context) func(uint32, uint32, uint32, uint32) int32 {
	return func(actionType uint32, index uint32, ptr uint32, length uint32) int32 {
		trx, err := context.GetApplyContext().GetPackedTransaction().GetTransaction()

		if err != nil {
			panic("could not unpack transaction")
		}

		var action *core.Action

		if actionType == 0 {
			if int(index) >= len(trx.ContextFreeActions) {
				return -1
			}

			action = trx.ContextFreeActions[index]
		} else if actionType == 1 {
			if int(index) >= len(trx.Actions) {
				return -1
			}

			action = trx.ContextFreeActions[index]
		}

		if action == nil {
			panic("action is not found")
		}

		ps, err := rlp.EncodeSize(action)

		if err != nil {
			panic(err)
		}

		if ps <= int(length) {
			data, _ := rlp.EncodeToBytes(action)
			context.WriteMemory(ptr, data)
		}

		return int32(ps)
	}
}