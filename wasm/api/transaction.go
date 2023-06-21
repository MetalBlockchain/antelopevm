package api

import (
	"github.com/MetalBlockchain/antelopevm/config"
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
)

func GetTransactionFunctions(context Context) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["send_inline"] = sendInline(context)
	functions["send_context_free_inline"] = sendContextFreeInline(context)
	functions["send_deferred"] = sendDeferred(context)
	functions["cancel_deferred"] = cancelDeferred(context)

	return functions
}

func sendInline(context Context) func(uint32, uint32) {
	return func(ptr uint32, length uint32) {
		if length >= uint32(config.MaxInlineActionSize) {
			panic("inline action too big")
		}

		data := context.ReadMemory(ptr, length)
		action := &core.Action{}

		if err := rlp.DecodeBytes(data, action); err != nil {
			panic("failed to decode action")
		}

		if err := context.GetApplyContext().ExecuteInline(*action); err != nil {
			panic(err)
		}
	}
}

func sendContextFreeInline(context Context) func(uint32, uint32) {
	return func(ptr uint32, length uint32) {
		if length >= uint32(config.MaxInlineActionSize) {
			panic("inline action too big")
		}

		data := context.ReadMemory(ptr, length)
		action := &core.Action{}

		if err := rlp.DecodeBytes(data, action); err != nil {
			panic("failed to decode action")
		}

		if err := context.GetApplyContext().ExecuteInline(*action); err != nil {
			panic(err)
		}
	}
}

func sendDeferred(context Context) func(uint32, name.AccountName, uint32, uint32, uint32) {
	return func(ptrSender uint32, payer name.AccountName, ptrData, ptrLength, replaceExisting uint32) {
		panic("not implemented")
	}
}

func cancelDeferred(context Context) func(uint32) int32 {
	return func(ptr uint32) int32 {
		return 0
	}
}
