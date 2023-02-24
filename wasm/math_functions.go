package wasm

import (
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/MetalBlockchain/antelopevm/math"
	"github.com/inconshreveable/log15"
)

func GetMathFunctions(c *ExecutionContext) map[string]interface{} {
	functions := make(map[string]interface{})
	functions["__ashlti3"] = __ashlti3(c)
	functions["__ashrti3"] = __ashrti3(c)
	functions["__lshlti3"] = __lshlti3(c)
	functions["__lshrti3"] = __lshrti3(c)
	functions["__divti3"] = __divti3(c)
	functions["__udivti3"] = __udivti3(c)
	functions["__multi3"] = __multi3(c)
	functions["__modti3"] = __modti3(c)
	functions["__umodti3"] = __umodti3(c)
	functions["__addtf3"] = __addtf3(c)
	functions["__subtf3"] = __subtf3(c)
	functions["__multf3"] = __multf3(c)
	functions["__divtf3"] = __divtf3(c)
	functions["__negtf2"] = __negtf2(c)

	return functions
}

func __ashlti3(context *ExecutionContext) func(uint32, uint64, uint64, uint32) {
	return func(ptr uint32, low uint64, high uint64, shift uint32) {
		i := math.Int128{Low: low, High: high}
		i.LeftShifts(int(shift))
		data, _ := rlp.EncodeToBytes(i)
		context.WriteMemory(ptr, data)
	}
}

func __ashrti3(context *ExecutionContext) func(uint32, uint64, uint64, uint32) {
	return func(ptr uint32, low uint64, high uint64, shift uint32) {
		i := math.Int128{Low: low, High: high}
		i.RightShifts(int(shift))
		data, _ := rlp.EncodeToBytes(i)
		context.WriteMemory(ptr, data)
	}
}

func __lshlti3(context *ExecutionContext) func(uint32, uint64, uint64, uint32) {
	return func(ptr uint32, low uint64, high uint64, shift uint32) {
		i := math.Int128{Low: low, High: high}
		i.LeftShifts(int(shift))
		data, _ := rlp.EncodeToBytes(i)
		context.WriteMemory(ptr, data)
	}
}

func __lshrti3(context *ExecutionContext) func(uint32, uint64, uint64, uint32) {
	return func(ptr uint32, low uint64, high uint64, shift uint32) {
		i := math.Int128{Low: low, High: high}
		i.RightShifts(int(shift))
		data, _ := rlp.EncodeToBytes(i)
		context.WriteMemory(ptr, data)
	}
}

func __divti3(context *ExecutionContext) func(uint32, uint64, uint64, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		lhs := math.Int128{Low: la, High: ha}
		rhs := math.Int128{Low: lb, High: hb}

		if rhs.IsZero() {
			panic("divide by zero")
		}

		quotient, _ := lhs.Div(rhs)
		data, _ := rlp.EncodeToBytes(quotient)
		context.WriteMemory(ptr, data)
	}
}

func __udivti3(context *ExecutionContext) func(uint32, uint64, uint64, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		lhs := math.Uint128{Low: la, High: ha}
		rhs := math.Uint128{Low: lb, High: hb}
		if rhs.IsZero() {
			panic("divide by zero")
		}
		quotient, _ := lhs.Div(rhs)
		data, _ := rlp.EncodeToBytes(quotient)
		context.WriteMemory(ptr, data)
	}
}

func __multi3(context *ExecutionContext) func(uint32, uint64, uint64, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		lhs := math.Int128{Low: la, High: ha}
		rhs := math.Int128{Low: lb, High: hb}
		data, _ := rlp.EncodeToBytes(lhs.Mul(rhs))
		context.WriteMemory(ptr, data)
	}
}

func __modti3(context *ExecutionContext) func(uint32, uint64, uint64, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		lhs := math.Int128{Low: la, High: ha}
		rhs := math.Int128{Low: lb, High: hb}
		if rhs.IsZero() {
			panic("divide by zero")
		}
		_, remainder := lhs.Div(rhs)
		data, _ := rlp.EncodeToBytes(remainder)
		context.WriteMemory(ptr, data)
	}
}

func __umodti3(context *ExecutionContext) func(uint32, uint64, uint64, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		lhs := math.Uint128{Low: la, High: ha}
		rhs := math.Uint128{Low: lb, High: hb}
		if rhs.IsZero() {
			panic("divide by zero")
		}
		_, remainder := lhs.Div(rhs)
		data, _ := rlp.EncodeToBytes(remainder)
		context.WriteMemory(ptr, data)
	}
}

func __addtf3(context *ExecutionContext) func(uint32, uint64, uint64, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		a := math.Float128{Low: la, High: ha}
		b := math.Float128{Low: lb, High: hb}

		data, _ := rlp.EncodeToBytes(a.Add(b))
		context.WriteMemory(ptr, data)
	}
}

func __subtf3(context *ExecutionContext) func(uint32, uint64, uint64, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		a := math.Float128{Low: la, High: ha}
		b := math.Float128{Low: lb, High: hb}

		data, _ := rlp.EncodeToBytes(a.Sub(b))
		context.WriteMemory(ptr, data)
	}
}

func __multf3(context *ExecutionContext) func(uint32, uint64, uint64, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		a := math.Float128{Low: la, High: ha}
		b := math.Float128{Low: lb, High: hb}

		data, _ := rlp.EncodeToBytes(a.Mul(b))
		context.WriteMemory(ptr, data)
	}
}

func __divtf3(context *ExecutionContext) func(uint32, uint64, uint64, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		a := math.Float128{Low: la, High: ha}
		b := math.Float128{Low: lb, High: hb}

		data, _ := rlp.EncodeToBytes(a.Div(b))
		context.WriteMemory(ptr, data)
	}
}

func __negtf2(context *ExecutionContext) func(uint32, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64) {
		high := uint64(ha)
		high ^= uint64(1) << 63
		f128 := math.Float128{Low: uint64(la), High: high}

		data, _ := rlp.EncodeToBytes(f128)
		context.WriteMemory(ptr, data)
	}
}

func Extendsftf2(context *ExecutionContext) func(uint32, float32) {
	return func(ptr uint32, f float32) {
		f32 := math.Float32(f)
		f128 := math.F32ToF128(f32)
		data, _ := rlp.EncodeToBytes(f128)
		context.WriteMemory(ptr, data)
	}
}

func Floatsitf(context *ExecutionContext) func(uint32, uint32) {
	return func(arg1 uint32, arg2 uint32) {
		log15.Info("Floatsitf")
	}
}

func Multf3(context *ExecutionContext) func(uint32, uint64, uint64, uint64, uint64) {
	return func(arg1 uint32, arg2 uint64, arg3 uint64, arg4 uint64, arg5 uint64) {
		log15.Info("Multf3")
	}
}

func Floatunsitf(context *ExecutionContext) func(uint32, uint32) {
	return func(arg1 uint32, arg2 uint32) {
		log15.Info("Floatunsitf")
	}
}

func Divtf3(context *ExecutionContext) func(uint32, uint64, uint64, uint64, uint64) {
	return func(arg1 uint32, arg2 uint64, arg3 uint64, arg4 uint64, arg5 uint64) {
		log15.Info("Divtf3")
	}
}

func Addtf3(context *ExecutionContext) func(uint32, uint64, uint64, uint64, uint64) {
	return func(arg1 uint32, arg2 uint64, arg3 uint64, arg4 uint64, arg5 uint64) {
		log15.Info("Addtf3")
	}
}

func Extenddftf2(context *ExecutionContext) func(uint32, float64) {
	return func(arg1 uint32, arg2 float64) {
		log15.Info("Extenddftf2")
	}
}

func Eqtf2(context *ExecutionContext) func(uint64, uint64, uint64, uint64) uint32 {
	return func(arg1 uint64, arg2 uint64, arg3 uint64, arg4 uint64) uint32 {
		log15.Info("Eqtf2")
		return 42
	}
}

func Letf2(context *ExecutionContext) func(uint64, uint64, uint64, uint64) uint32 {
	return func(arg1 uint64, arg2 uint64, arg3 uint64, arg4 uint64) uint32 {
		log15.Info("Letf2")
		return 42
	}
}

func Netf2(context *ExecutionContext) func(uint64, uint64, uint64, uint64) uint32 {
	return func(arg1 uint64, arg2 uint64, arg3 uint64, arg4 uint64) uint32 {
		log15.Info("Netf2")
		return 42
	}
}

func Subtf3(context *ExecutionContext) func(uint32, uint64, uint64, uint64, uint64) {
	return func(arg1 uint32, arg2 uint64, arg3 uint64, arg4 uint64, arg5 uint64) {
		log15.Info("Subtf3")
	}
}

func Trunctfdf2(context *ExecutionContext) func(uint64, uint64) float64 {
	return func(arg1 uint64, arg2 uint64) float64 {
		log15.Info("Trunctfdf2")
		return float64(1)
	}
}

func Getf2(context *ExecutionContext) func(uint64, uint64, uint64, uint64) uint32 {
	return func(arg1 uint64, arg2 uint64, arg3 uint64, arg4 uint64) uint32 {
		log15.Info("Getf2")
		return 1
	}
}

func Trunctfsf2(context *ExecutionContext) func(uint64, uint64) float32 {
	return func(arg1 uint64, arg2 uint64) float32 {
		log15.Info("Trunctfsf2")
		return float32(1)
	}
}

func Unordtf2(context *ExecutionContext) func(uint64, uint64, uint64, uint64) uint32 {
	return func(arg1 uint64, arg2 uint64, arg3 uint64, arg4 uint64) uint32 {
		log15.Info("Unordtf2")
		return 1
	}
}

func Fixunstfsi(context *ExecutionContext) func(uint64, uint64) uint32 {
	return func(arg1 uint64, arg2 uint64) uint32 {
		log15.Info("Fixunstfsi")
		return 1
	}
}

func Fixtfsi(context *ExecutionContext) func(uint64, uint64) uint32 {
	return func(arg1 uint64, arg2 uint64) uint32 {
		log15.Info("Fixtfsi")
		return 1
	}
}
