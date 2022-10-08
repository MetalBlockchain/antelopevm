package wasm

import "github.com/wasmerio/wasmer-go/wasmer"

func MemSet(context *ExecutionContext) *wasmer.Function {
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

func MemCopy(context *ExecutionContext) *wasmer.Function {
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

func MemMove(context *ExecutionContext) *wasmer.Function {
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
