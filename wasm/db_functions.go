package wasm

import (
	"fmt"

	"github.com/wasmerio/wasmer-go/wasmer"
)

func FindI64(context *ExecutionContext) *wasmer.Function {
	return wasmer.NewFunction(
		context.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I64, wasmer.I64, wasmer.I64, wasmer.I64),
			wasmer.NewValueTypes(wasmer.I32),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			return []wasmer.Value{wasmer.NewI32(42)}, nil
		},
	)
}

func StoreI64(context *ExecutionContext) *wasmer.Function {
	return wasmer.NewFunction(
		context.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I64, wasmer.I64, wasmer.I64, wasmer.I64, wasmer.I32, wasmer.I32),
			wasmer.NewValueTypes(wasmer.I32),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			return []wasmer.Value{wasmer.NewI32(42)}, nil
		},
	)
}

func UpdateI64(context *ExecutionContext) *wasmer.Function {
	return wasmer.NewFunction(
		context.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32, wasmer.I64, wasmer.I32, wasmer.I32),
			wasmer.NewValueTypes(),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			return []wasmer.Value{wasmer.NewI32(42)}, nil
		},
	)
}

func NextI64(context *ExecutionContext) *wasmer.Function {
	return wasmer.NewFunction(
		context.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32, wasmer.I32),
			wasmer.NewValueTypes(wasmer.I32),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			return []wasmer.Value{wasmer.NewI32(42)}, nil
		},
	)
}

func GetI64(context *ExecutionContext) *wasmer.Function {
	return wasmer.NewFunction(
		context.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32),
			wasmer.NewValueTypes(wasmer.I32),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			return []wasmer.Value{wasmer.NewI32(42)}, nil
		},
	)
}

func RemoveI64(context *ExecutionContext) *wasmer.Function {
	return wasmer.NewFunction(
		context.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32),
			wasmer.NewValueTypes(),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			return []wasmer.Value{wasmer.NewI32(42)}, nil
		},
	)
}

func ReacActionData(context *ExecutionContext) *wasmer.Function {
	return wasmer.NewFunction(
		context.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32, wasmer.I32),
			wasmer.NewValueTypes(wasmer.I32),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			fmt.Println("ReacActionData")
			msg := args[0].I32()
			len := args[1].I32()
			fmt.Printf("%v %v", msg, len)
			return []wasmer.Value{wasmer.NewI32(42)}, nil
		},
	)
}

func ActionDataSize(context *ExecutionContext) *wasmer.Function {
	return wasmer.NewFunction(
		context.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(),
			wasmer.NewValueTypes(wasmer.I32),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			return []wasmer.Value{wasmer.NewI32(42)}, nil
		},
	)
}
