package wasm

import (
	"bytes"
	"fmt"
)

func EosIoAssert(context *ExecutionContext) func(uint32, uint32) {
	return func(arg1 uint32, arg2 uint32) {
		fmt.Printf("eosio_assert %v %v\n", arg1, arg2)

		if arg1 == 0 {
			memory := context.module.Memory()
			data, found := memory.Read(context.context, arg2, memory.Size(context.context)-arg2)

			if !found {
				panic("not foundsss")
			}

			n := bytes.IndexByte(data, 0)

			if n > 0 {
				panic(string(data[:n]))
			} else {
				panic("assertion failed")
			}
		}
	}
}

func Abort(context *ExecutionContext) func() {
	return func() {
		panic("abort called")
	}
}

func AssertCode(context *ExecutionContext) func(uint32, uint64) {
	return func(arg1 uint32, arg2 uint64) {
		fmt.Println("AssertCode")
	}
}
