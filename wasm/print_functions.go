package wasm

import "fmt"

func GetConsoleFunctions(context *ExecutionContext) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["prints"] = prints(context)
	functions["prints_l"] = prints_l(context)
	functions["printi"] = printi(context)
	functions["printui"] = printui(context)
	functions["printi128"] = printi128(context)
	functions["printui128"] = printui128(context)
	functions["printsf"] = printsf(context)
	functions["printdf"] = printdf(context)
	functions["printqf"] = printqf(context)
	functions["printn"] = printn(context)
	functions["printhex"] = printhex(context)

	return functions
}

func prints(context *ExecutionContext) func(uint32) {
	return func(arg1 uint32) {
		fmt.Printf("prints %v\n", arg1)
	}
}

func prints_l(context *ExecutionContext) func(uint32, uint32) {
	return func(arg1 uint32, arg2 uint32) {
		fmt.Printf("prints_l %v\n", arg1)
	}
}

func printi(context *ExecutionContext) func(int64) {
	return func(arg1 int64) {
		fmt.Printf("printi %v\n", arg1)
	}
}

func printui(context *ExecutionContext) func(uint64) {
	return func(arg1 uint64) {
		fmt.Printf("printui %v\n", arg1)
	}
}

func printi128(context *ExecutionContext) func(uint32, uint32) {
	return func(arg1 uint32, arg2 uint32) {
		fmt.Printf("printi128 %v\n", arg1)
	}
}

func printui128(context *ExecutionContext) func(uint32, uint32) {
	return func(arg1 uint32, arg2 uint32) {
		fmt.Printf("printi128 %v\n", arg1)
	}
}

func printsf(context *ExecutionContext) func(float32) {
	return func(arg1 float32) {
		fmt.Printf("printsf %v\n", arg1)
	}
}

func printdf(context *ExecutionContext) func(float64) {
	return func(arg1 float64) {
		fmt.Printf("printdf %v\n", arg1)
	}
}

func printqf(context *ExecutionContext) func(uint32, uint32) {
	return func(arg1 uint32, arg2 uint32) {
		fmt.Printf("printqf %v\n", arg1)
	}
}

func printn(context *ExecutionContext) func(uint64) {
	return func(arg1 uint64) {
		fmt.Printf("printn %v\n", arg1)
	}
}

func printhex(context *ExecutionContext) func(uint32, uint32) {
	return func(arg1 uint32, arg2 uint32) {
		fmt.Printf("printhex %v\n", arg1)
	}
}
