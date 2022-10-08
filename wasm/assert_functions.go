package wasm

import (
	"errors"
	"fmt"

	"github.com/wasmerio/wasmer-go/wasmer"
)

func EosIoAssert(context *ExecutionContext) *wasmer.Function {
	return wasmer.NewFunction(
		context.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32, wasmer.I32),
			wasmer.NewValueTypes(),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			fmt.Println("eosio_assert")
			assertion := args[1].I32()

			if assertion == 0 {
				return nil, errors.New("assertion failed")
			}

			return nil, nil
		},
	)
}

func Abort(context *ExecutionContext) *wasmer.Function {
	return wasmer.NewFunction(
		context.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(),
			wasmer.NewValueTypes(),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			fmt.Println("Abort")
			return nil, errors.New("abort requested")
		},
	)
}

func AssertCode(context *ExecutionContext) *wasmer.Function {
	return wasmer.NewFunction(
		context.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32, wasmer.I64),
			wasmer.NewValueTypes(),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			fmt.Println("AssertCode")
			return nil, errors.New("abort requested")
		},
	)
}
