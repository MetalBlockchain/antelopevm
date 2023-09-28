package api

import (
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/transaction"
	"github.com/MetalBlockchain/antelopevm/config"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
)

func init() {
	Functions["send_inline"] = sendInline
	Functions["send_context_free_inline"] = sendContextFreeInline
	Functions["send_deferred"] = sendDeferred
	Functions["cancel_deferred"] = cancelDeferred
}

func sendInline(context Context) interface{} {
	return func(ptr uint32, length uint32) {
		if length >= uint32(config.MaxInlineActionSize) {
			panic("inline action too big")
		}

		data := context.ReadMemory(ptr, length)
		action := &transaction.Action{}

		if err := rlp.DecodeBytes(data, action); err != nil {
			panic("failed to decode action")
		}

		if err := context.GetApplyContext().ExecuteInline(*action); err != nil {
			panic(err)
		}
	}
}

func sendContextFreeInline(context Context) interface{} {
	return func(ptr uint32, length uint32) {
		if length >= uint32(config.MaxInlineActionSize) {
			panic("inline action too big")
		}

		data := context.ReadMemory(ptr, length)
		action := &transaction.Action{}

		if err := rlp.DecodeBytes(data, action); err != nil {
			panic("failed to decode action")
		}

		if err := context.GetApplyContext().ExecuteInline(*action); err != nil {
			panic(err)
		}
	}
}

func sendDeferred(context Context) interface{} {
	return func(ptrSender uint32, payer name.AccountName, ptrData, ptrLength, replaceExisting uint32) {
		panic("not implemented")
	}
}

func cancelDeferred(context Context) interface{} {
	return func(ptr uint32) int32 {
		return 0
	}
}
