package api

import (
	"encoding/hex"
	"math"
	"strconv"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	antelopeMath "github.com/MetalBlockchain/antelopevm/math"
	log "github.com/inconshreveable/log15"
)

func GetConsoleFunctions(context Context) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["prints"] = prints(context)
	functions["prints_l"] = prints_l(context)
	functions["printi"] = printi(context)
	functions["printui"] = printui(context)
	functions["printi128"] = printi128(context)
	functions["printui128"] = printui128(context)
	functions["printsf"] = printsf(context)
	functions["printdf"] = printdf(context)
	functions["printqf"] = printqf(context)
	functions["printn"] = printn(context)
	functions["printhex"] = printhex(context)

	return functions
}

func prints(context Context) func(uint32) {
	return func(ptr uint32) {
		log.Info("prints", "ptr", ptr)
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

func prints_l(context Context) func(uint32, uint32) {
	return func(ptr uint32, length uint32) {
		log.Info("prints_l", "ptr", ptr, "length", length)
		data := context.ReadMemory(ptr, length)
		text := string(data)
		context.GetApplyContext().ConsoleAppend(text)
	}
}

func printi(context Context) func(int64) {
	return func(value int64) {
		log.Info("printi", "value", value)
		text := strconv.FormatInt(value, 10)
		context.GetApplyContext().ConsoleAppend(text)
	}
}

func printui(context Context) func(uint64) {
	return func(value uint64) {
		log.Info("printui", "value", value)
		text := strconv.FormatUint(value, 10)
		context.GetApplyContext().ConsoleAppend(text)
	}
}

func printi128(context Context) func(uint32, uint32) {
	return func(ptr uint32, length uint32) {
		log.Info("printi128", "ptr", ptr, "length", length)
		data := context.ReadMemory(ptr, length)
		var value antelopeMath.Int128

		if err := rlp.DecodeBytes(data, &value); err != nil {
			panic("could not decode RLP")
		}

		context.GetApplyContext().ConsoleAppend(value.String())
	}
}

func printui128(context Context) func(uint32, uint32) {
	return func(ptr uint32, length uint32) {
		log.Info("printui128", "ptr", ptr, "length", length)
		data := context.ReadMemory(ptr, length)
		var value antelopeMath.Uint128

		if err := rlp.DecodeBytes(data, &value); err != nil {
			panic("could not decode RLP")
		}

		context.GetApplyContext().ConsoleAppend(value.String())
	}
}

func printsf(context Context) func(uint32) {
	return func(value uint32) {
		log.Info("printsf", "value", value)
		val := math.Float32frombits(value)
		text := strconv.FormatFloat(float64(val), 'e', 6, 32)
		context.GetApplyContext().ConsoleAppend(text)
	}
}

func printdf(context Context) func(uint64) {
	return func(value uint64) {
		log.Info("printdf", "value", value)
		val := math.Float64frombits(value)
		text := strconv.FormatFloat(val, 'e', 15, 64)
		context.GetApplyContext().ConsoleAppend(text)
	}
}

func printqf(context Context) func(uint32, uint32) {
	return func(ptr uint32, length uint32) {
		log.Info("printqf", "ptr", ptr, "length", length)
		data := context.ReadMemory(ptr, length)
		var value antelopeMath.Float128

		if err := rlp.DecodeBytes(data, &value); err != nil {
			panic("could not decode RLP")
		}

		context.GetApplyContext().ConsoleAppend(value.String())
	}
}

func printn(context Context) func(core.Name) {
	return func(value core.Name) {
		log.Info("printn", "value", value)
		context.GetApplyContext().ConsoleAppend(value.String())
	}
}

func printhex(context Context) func(uint32, uint32) {
	return func(ptr uint32, length uint32) {
		log.Info("printhex", "ptr", ptr, "length", length)
		data := context.ReadMemory(ptr, length)
		text := hex.EncodeToString(data)
		context.GetApplyContext().ConsoleAppend(text)
	}
}
