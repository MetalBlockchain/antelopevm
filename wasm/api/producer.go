package api

import (
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/MetalBlockchain/antelopevm/utils"
)

func init() {
	Functions["get_active_producers"] = getActiveProducers
}

func getActiveProducers(context Context) interface{} {
	return func(ptr uint32, length uint32) int32 {
		producers, err := context.GetController().GetActiveProducers()
		if err != nil {
			panic(err)
		} else if len(producers) == 0 {
			return 0
		}

		data, err := rlp.EncodeToBytes(producers)
		if err != nil {
			panic(err)
		}
		s := len(data)
		if length == 0 {
			return int32(s)
		}
		copySize := utils.MinInt(s, int(length))
		context.WriteMemory(ptr, data[0:copySize])
		return int32(copySize)
	}
}
