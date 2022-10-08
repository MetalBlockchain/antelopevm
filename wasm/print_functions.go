package wasm

import "github.com/wasmerio/wasmer-go/wasmer"

func PrintsL(context *ExecutionContext) *wasmer.Function {
	return wasmer.NewFunction(
		context.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32, wasmer.I32),
			wasmer.NewValueTypes(),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			return []wasmer.Value{wasmer.NewI32(42)}, nil
		},
	)
}
