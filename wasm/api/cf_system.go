package api

import (
	"fmt"

	"github.com/inconshreveable/log15"
)

func init() {
	Functions["abort"] = abort
	Functions["eosio_assert"] = assert
	Functions["eosio_assert_message"] = assertMessage
	Functions["eosio_assert_code"] = assertCode
	Functions["eosio_exit"] = exit
}

func GetContextFreeSystemFunctions(context Context) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["abort"] = abort(context)
	functions["eosio_assert"] = assert(context)
	functions["eosio_assert_message"] = assertMessage(context)
	functions["eosio_assert_code"] = assertCode(context)
	functions["eosio_exit"] = exit(context)

	return functions
}

func abort(context Context) interface{} {
	return func() {
		panic("abort called")
	}
}

func assert(context Context) interface{} {
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

			log15.Info("f", "data", data, "size", size, "ptr", ptr)

			text := string(data[:size])

			panic(fmt.Sprintf("assertion failure with message: %s", text))
		}
	}
}

func assertMessage(context Context) interface{} {
	return func(condition uint32, ptr uint32, size uint32) {
		if condition == 0 {
			data := context.ReadMemory(ptr, size)

			panic(fmt.Sprintf("assertion failure with message: %s", string(data)))
		}
	}
}

func assertCode(context Context) interface{} {
	return func(condition uint32, code uint64) {
		if condition == 0 {
			panic(fmt.Sprintf("assertion failure with error code: %v", code))
		}
	}
}

func exit(context Context) interface{} {
	return func(int32) {
		panic("eosio_exit called")
	}
}
