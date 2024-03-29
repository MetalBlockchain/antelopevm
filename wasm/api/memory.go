package api

import (
	"bytes"

	"github.com/MetalBlockchain/antelopevm/utils"
)

func init() {
	Functions["memset"] = MemSet
	Functions["memcpy"] = MemCopy
	Functions["memmove"] = MemMove
	Functions["memcmp"] = MemCmp
}

func MemSet(context Context) interface{} {
	return func(dest uint32, value int32, length uint32) uint32 {
		destData := context.ReadMemory(dest, length)
		memset(destData, byte(value), int(length))
		return dest
	}
}

func MemCopy(context Context) interface{} {
	return func(dest uint32, source uint32, length uint32) uint32 {
		if utils.AbsInt32(int32(dest-source)) < int32(length) {
			panic("memcpy can only accept non-aliasing pointers")
		}

		sourceData := context.ReadMemory(source, length)
		destData := context.ReadMemory(dest, length)
		memcpy(destData, sourceData, length)

		return dest
	}
}

func MemMove(context Context) interface{} {
	return func(dest uint32, source uint32, length uint32) uint32 {
		sourceData := context.ReadMemory(source, length)
		destData := context.ReadMemory(dest, length)
		memcpy(destData, sourceData, length)
		return dest
	}
}

func MemCmp(context Context) interface{} {
	return func(dest uint32, source uint32, length uint32) int32 {
		sourceData := context.ReadMemory(source, length)
		destData := context.ReadMemory(dest, length)

		return memcmp(destData, sourceData, length)
	}
}

func memset(dest []byte, value byte, length int) {
	bytes := bytes.Repeat([]byte{value}, length)
	copy(dest[:length], bytes)
}

func memcmp(dest []byte, source []byte, length uint32) int32 {
	return int32(bytes.Compare(dest[:length], source[:length]))
}

func memcpy(dest []byte, source []byte, length uint32) int {
	return copy(dest[:length], source[:length])
}
