package api

import (
	"bytes"
	"math"

	log "github.com/inconshreveable/log15"
)

func GetMemoryFunctions(context Context) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["memset"] = MemSet(context)
	functions["memcpy"] = MemCopy(context)
	functions["memmove"] = MemMove(context)
	functions["memcmp"] = MemCmp(context)

	return functions
}

func MemSet(context Context) func(uint32, uint32, uint32) uint32 {
	return func(dest uint32, value uint32, length uint32) uint32 {
		log.Info("memset", "dest", dest, "value", value, "length", length)
		destData := context.ReadMemory(dest, length)
		memset(destData, byte(value), length)
		return dest
	}
}

func MemCopy(context Context) func(uint32, uint32, uint32) uint32 {
	return func(dest uint32, source uint32, length uint32) uint32 {
		log.Info("memcopy", "dest", dest, "source", source, "length", length)
		if math.Abs(float64(dest)-float64(source)) < float64(length) {
			panic("memcpy can only accept non-aliasing pointers")
		}

		sourceData := context.ReadMemory(source, length)
		context.WriteMemory(dest, sourceData)

		return dest
	}
}

func MemMove(context Context) func(uint32, uint32, uint32) uint32 {
	return func(dest uint32, source uint32, length uint32) uint32 {
		log.Info("memmove", "dest", dest, "source", source, "length", length)
		sourceData := context.ReadMemory(source, length)
		destData := context.ReadMemory(dest, length)
		memcpy(destData, sourceData, length)
		return dest
	}
}

func MemCmp(context Context) func(uint32, uint32, uint32) int32 {
	return func(dest uint32, source uint32, length uint32) int32 {
		log.Info("memcmp", "dest", dest, "source", source, "length", length)
		sourceData := context.ReadMemory(source, length)
		destData := context.ReadMemory(dest, length)

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
