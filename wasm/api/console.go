package api

import (
	"encoding/hex"
	"math"
	"strconv"

	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	antelopeMath "github.com/MetalBlockchain/antelopevm/math"
)

func init() {
	Functions["prints"] = prints
	Functions["prints_l"] = prints_l
	Functions["printi"] = printi
	Functions["printui"] = printui
	Functions["printi128"] = printi128
	Functions["printui128"] = printui128
	Functions["printsf"] = printsf
	Functions["printdf"] = printdf
	Functions["printqf"] = printqf
	Functions["printn"] = printn
	Functions["printhex"] = printhex
}

func prints(context Context) interface{} {
	return func(ptr uint32) {
		data := context.ReadMemory(ptr, 512)
		var size uint32

		for i := 0; i < len(data); i++ {
			if data[i] == 0 {
				break
			}

			size++
		}

		text := string(data[ptr : ptr+size])
		context.GetApplyContext().ConsoleAppend(text)
	}
}

func prints_l(context Context) interface{} {
	return func(ptr uint32, length uint32) {
		data := context.ReadMemory(ptr, length)
		text := string(data)
		context.GetApplyContext().ConsoleAppend(text)
	}
}

func printi(context Context) interface{} {
	return func(value int64) {
		text := strconv.FormatInt(value, 10)
		context.GetApplyContext().ConsoleAppend(text)
	}
}

func printui(context Context) interface{} {
	return func(value uint64) {
		text := strconv.FormatUint(value, 10)
		context.GetApplyContext().ConsoleAppend(text)
	}
}

func printi128(context Context) interface{} {
	return func(ptr uint32, length uint32) {
		data := context.ReadMemory(ptr, length)
		var value antelopeMath.Int128

		if err := rlp.DecodeBytes(data, &value); err != nil {
			panic("could not decode RLP")
		}

		context.GetApplyContext().ConsoleAppend(value.String())
	}
}

func printui128(context Context) interface{} {
	return func(ptr uint32, length uint32) {
		data := context.ReadMemory(ptr, length)
		var value antelopeMath.Uint128

		if err := rlp.DecodeBytes(data, &value); err != nil {
			panic("could not decode RLP")
		}

		context.GetApplyContext().ConsoleAppend(value.String())
	}
}

func printsf(context Context) interface{} {
	return func(value uint32) {
		val := math.Float32frombits(value)
		text := strconv.FormatFloat(float64(val), 'e', 6, 32)
		context.GetApplyContext().ConsoleAppend(text)
	}
}

func printdf(context Context) interface{} {
	return func(value uint64) {
		val := math.Float64frombits(value)
		text := strconv.FormatFloat(val, 'e', 15, 64)
		context.GetApplyContext().ConsoleAppend(text)
	}
}

func printqf(context Context) interface{} {
	return func(ptr uint32, length uint32) {
		data := context.ReadMemory(ptr, length)
		var value antelopeMath.Float128

		if err := rlp.DecodeBytes(data, &value); err != nil {
			panic("could not decode RLP")
		}

		context.GetApplyContext().ConsoleAppend(value.String())
	}
}

func printn(context Context) interface{} {
	return func(value name.Name) {
		context.GetApplyContext().ConsoleAppend(value.String())
	}
}

func printhex(context Context) interface{} {
	return func(ptr uint32, length uint32) {
		data := context.ReadMemory(ptr, length)
		text := hex.EncodeToString(data)
		context.GetApplyContext().ConsoleAppend(text)
	}
}
