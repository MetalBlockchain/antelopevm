package wasm

import "github.com/wasmerio/wasmer-go/wasmer"

func RequireAuth(context *ExecutionContext) *wasmer.Function {
	return wasmer.NewFunction(
		context.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I64),
			wasmer.NewValueTypes(),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			return []wasmer.Value{wasmer.NewI32(42)}, nil
		},
	)
}

func CurrentReceiver(context *ExecutionContext) *wasmer.Function {
	return wasmer.NewFunction(
		context.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(),
			wasmer.NewValueTypes(wasmer.I64),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			return []wasmer.Value{wasmer.NewI32(42)}, nil
		},
	)
}

func IsAccount(context *ExecutionContext) *wasmer.Function {
	return wasmer.NewFunction(
		context.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I64),
			wasmer.NewValueTypes(wasmer.I32),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			return []wasmer.Value{wasmer.NewI32(42)}, nil
		},
	)
}

func RequireRecipient(context *ExecutionContext) *wasmer.Function {
	return wasmer.NewFunction(
		context.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I64),
			wasmer.NewValueTypes(),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			return []wasmer.Value{wasmer.NewI32(42)}, nil
		},
	)
}

func HasAuth(context *ExecutionContext) *wasmer.Function {
	return wasmer.NewFunction(
		context.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I64),
			wasmer.NewValueTypes(wasmer.I32),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			return []wasmer.Value{wasmer.NewI32(42)}, nil
		},
	)
}
