package api

import (
	"fmt"
)

func GetContextFreeSystemFunctions(context Context) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["abort"] = abort(context)
	functions["eosio_assert"] = assert(context)
	functions["eosio_assert_message"] = assertMessage(context)
	functions["eosio_assert_code"] = assertCode(context)
	functions["eosio_exit"] = exit(context)

	return functions
}

func abort(context Context) func() {
	return func() {
		panic("abort called")
	}
}

func assert(context Context) func(uint32, uint32) {
	return func(condition uint32, ptr uint32) {
		if condition == 0 {
			data := context.ReadMemory(ptr, 512)
			var size uint32

			for i := 0; i < len(data); i++ {
				if data[i] == 0 {
					break
				}

				size++
			}

			text := string(data[ptr : ptr+size])

			panic(fmt.Sprintf("assertion failure with message: %s", text))
		}
	}
}

func assertMessage(context Context) func(uint32, uint32, uint32) {
	return func(condition uint32, ptr uint32, size uint32) {
		if condition == 0 {
			data := context.ReadMemory(ptr, size)

			panic(fmt.Sprintf("assertion failure with message: %s", string(data)))
		}
	}
}

func assertCode(context Context) func(uint32, uint64) {
	return func(condition uint32, code uint64) {
		if condition == 0 {
			panic(fmt.Sprintf("assertion failure with error code: %v", code))
		}
	}
}

func exit(context Context) func(int32) {
	return func(int32) {
		panic("eosio_exit called")
	}
}
