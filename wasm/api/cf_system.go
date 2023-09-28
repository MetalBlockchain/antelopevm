package api

import (
	"fmt"
)

func init() {
	Functions["abort"] = abort
	Functions["eosio_assert"] = assert
	Functions["eosio_assert_message"] = assertMessage
	Functions["eosio_assert_code"] = assertCode
	Functions["eosio_exit"] = exit
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
