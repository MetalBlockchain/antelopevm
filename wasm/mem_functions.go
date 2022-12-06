package wasm

import (
	//#include <string.h>
	"C"
)
import (
	"bytes"
	"math"
)

func GetMemoryFunctions(context *ExecutionContext) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["memset"] = MemSet(context)
	functions["memcpy"] = MemCopy(context)
	functions["memmove"] = MemMove(context)
	functions["memcmp"] = MemCmp(context)

	return functions
}

func MemSet(context *ExecutionContext) func(uint32, uint32, uint32) uint32 {
	return func(dest uint32, value uint32, length uint32) uint32 {
		memory := context.module.Memory()
		destData, _ := memory.Read(context.context, dest, length)
		memset(destData, byte(value), length)
		return dest
	}
}

func MemCopy(context *ExecutionContext) func(uint32, uint32, uint32) uint32 {
	return func(dest uint32, source uint32, length uint32) uint32 {
		if math.Abs(float64(dest)-float64(source)) < float64(length) {
			panic("memcpy can only accept non-aliasing pointers")
		}

		memory := context.module.Memory()
		sourceData, _ := memory.Read(context.context, source, length)
		destData, _ := memory.Read(context.context, dest, length)
		memcpy(destData, sourceData, length)

		return dest
	}
}

func MemMove(context *ExecutionContext) func(uint32, uint32, uint32) uint32 {
	return func(dest uint32, source uint32, length uint32) uint32 {
		memory := context.module.Memory()
		sourceData, _ := memory.Read(context.context, source, length)
		destData, _ := memory.Read(context.context, dest, length)
		memcpy(destData, sourceData, length)
		return dest
	}
}

func MemCmp(context *ExecutionContext) func(uint32, uint32, uint32) int32 {
	return func(dest uint32, source uint32, length uint32) int32 {
		memory := context.module.Memory()
		sourceData, _ := memory.Read(context.context, source, length)
		destData, _ := memory.Read(context.context, dest, length)

		return memcmp(destData, sourceData, length)
	}
}

func memset(dest []byte, value byte, length uint32) {
	for i := uint32(0); i < length; i++ {
		dest[i] = value
	}
}

func memcmp(dest []byte, source []byte, length uint32) int32 {
	result := bytes.Compare(dest[:length], source[:length])

	if result < 0 {
		return -1
	} else if result > 0 {
		return 1
	}

	return 0
}

func memcpy(dest []byte, source []byte, length uint32) int {
	return copy(dest, source[:length])
}
