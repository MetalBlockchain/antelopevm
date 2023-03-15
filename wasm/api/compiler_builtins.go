package api

import (
	"math"

	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	eosMath "github.com/MetalBlockchain/antelopevm/math"
)

func GetCompilerBuiltinFunctions(c Context) map[string]interface{} {
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
	functions["__extendsftf2"] = __extendsftf2(c)
	functions["__extenddftf2"] = __extenddftf2(c)
	functions["__trunctfdf2"] = __trunctfdf2(c)
	functions["__trunctfsf2"] = __trunctfsf2(c)
	functions["__fixtfsi"] = __fixtfsi(c)
	functions["__fixtfdi"] = __fixtfdi(c)
	functions["__fixtfti"] = __fixtfti(c)
	functions["__fixunstfsi"] = __fixunstfsi(c)
	functions["__fixunstfdi"] = __fixunstfdi(c)
	functions["__fixunstfti"] = __fixunstfti(c)
	functions["__fixsfti"] = __fixsfti(c)
	functions["__fixdfti"] = __fixdfti(c)
	functions["__fixunssfti"] = __fixunssfti(c)
	functions["__fixunsdfti"] = __fixunsdfti(c)
	functions["__floatsidf"] = __floatsidf(c)
	functions["__floatsitf"] = __floatsitf(c)
	functions["__floatditf"] = __floatditf(c)
	functions["__floatunsitf"] = __floatunsitf(c)
	functions["__floatunditf"] = __floatunditf(c)
	functions["__floattidf"] = __floattidf(c)
	functions["__floatuntidf"] = __floatuntidf(c)
	functions["__eqtf2"] = __eqtf2(c)
	functions["__netf2"] = __netf2(c)
	functions["__getf2"] = __getf2(c)
	functions["__gttf2"] = __gttf2(c)
	functions["__letf2"] = __letf2(c)
	functions["__lttf2"] = __lttf2(c)
	functions["__cmptf2"] = __cmptf2(c)
	functions["__unordtf2"] = __unordtf2(c)

	return functions
}

func __ashlti3(context Context) func(uint32, uint64, uint64, uint32) {
	return func(ptr uint32, low uint64, high uint64, shift uint32) {
		i := eosMath.Int128{Low: low, High: high}
		i.LeftShifts(int(shift))
		data, _ := rlp.EncodeToBytes(i)
		context.WriteMemory(ptr, data)
	}
}

func __ashrti3(context Context) func(uint32, uint64, uint64, uint32) {
	return func(ptr uint32, low uint64, high uint64, shift uint32) {
		i := eosMath.Int128{Low: low, High: high}
		i.RightShifts(int(shift))
		data, _ := rlp.EncodeToBytes(i)
		context.WriteMemory(ptr, data)
	}
}

func __lshlti3(context Context) func(uint32, uint64, uint64, uint32) {
	return func(ptr uint32, low uint64, high uint64, shift uint32) {
		i := eosMath.Int128{Low: low, High: high}
		i.LeftShifts(int(shift))
		data, _ := rlp.EncodeToBytes(i)
		context.WriteMemory(ptr, data)
	}
}

func __lshrti3(context Context) func(uint32, uint64, uint64, uint32) {
	return func(ptr uint32, low uint64, high uint64, shift uint32) {
		i := eosMath.Int128{Low: low, High: high}
		i.RightShifts(int(shift))
		data, _ := rlp.EncodeToBytes(i)
		context.WriteMemory(ptr, data)
	}
}

func __divti3(context Context) func(uint32, uint64, uint64, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		lhs := eosMath.Int128{Low: la, High: ha}
		rhs := eosMath.Int128{Low: lb, High: hb}

		if rhs.IsZero() {
			panic("divide by zero")
		}

		quotient, _ := lhs.Div(rhs)
		data, _ := rlp.EncodeToBytes(quotient)
		context.WriteMemory(ptr, data)
	}
}

func __udivti3(context Context) func(uint32, uint64, uint64, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		lhs := eosMath.Uint128{Low: la, High: ha}
		rhs := eosMath.Uint128{Low: lb, High: hb}
		if rhs.IsZero() {
			panic("divide by zero")
		}
		quotient, _ := lhs.Div(rhs)
		data, _ := rlp.EncodeToBytes(quotient)
		context.WriteMemory(ptr, data)
	}
}

func __multi3(context Context) func(uint32, uint64, uint64, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		lhs := eosMath.Int128{Low: la, High: ha}
		rhs := eosMath.Int128{Low: lb, High: hb}
		data, _ := rlp.EncodeToBytes(lhs.Mul(rhs))
		context.WriteMemory(ptr, data)
	}
}

func __modti3(context Context) func(uint32, uint64, uint64, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		lhs := eosMath.Int128{Low: la, High: ha}
		rhs := eosMath.Int128{Low: lb, High: hb}
		if rhs.IsZero() {
			panic("divide by zero")
		}
		_, remainder := lhs.Div(rhs)
		data, _ := rlp.EncodeToBytes(remainder)
		context.WriteMemory(ptr, data)
	}
}

func __umodti3(context Context) func(uint32, uint64, uint64, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		lhs := eosMath.Uint128{Low: la, High: ha}
		rhs := eosMath.Uint128{Low: lb, High: hb}
		if rhs.IsZero() {
			panic("divide by zero")
		}
		_, remainder := lhs.Div(rhs)
		data, _ := rlp.EncodeToBytes(remainder)
		context.WriteMemory(ptr, data)
	}
}

func __addtf3(context Context) func(uint32, uint64, uint64, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		a := eosMath.Float128{Low: la, High: ha}
		b := eosMath.Float128{Low: lb, High: hb}

		data, _ := rlp.EncodeToBytes(a.Add(b))
		context.WriteMemory(ptr, data)
	}
}

func __subtf3(context Context) func(uint32, uint64, uint64, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		a := eosMath.Float128{Low: la, High: ha}
		b := eosMath.Float128{Low: lb, High: hb}

		data, _ := rlp.EncodeToBytes(a.Sub(b))
		context.WriteMemory(ptr, data)
	}
}

func __multf3(context Context) func(uint32, uint64, uint64, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		a := eosMath.Float128{Low: la, High: ha}
		b := eosMath.Float128{Low: lb, High: hb}

		data, _ := rlp.EncodeToBytes(a.Mul(b))
		context.WriteMemory(ptr, data)
	}
}

func __divtf3(context Context) func(uint32, uint64, uint64, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		a := eosMath.Float128{Low: la, High: ha}
		b := eosMath.Float128{Low: lb, High: hb}

		data, _ := rlp.EncodeToBytes(a.Div(b))
		context.WriteMemory(ptr, data)
	}
}

func __negtf2(context Context) func(uint32, uint64, uint64) {
	return func(ptr uint32, la uint64, ha uint64) {
		high := uint64(ha)
		high ^= uint64(1) << 63
		f128 := eosMath.Float128{Low: uint64(la), High: high}

		data, _ := rlp.EncodeToBytes(f128)
		context.WriteMemory(ptr, data)
	}
}

func __extendsftf2(context Context) func(uint32, uint32) {
	return func(ptr uint32, f uint32) {
		f32 := eosMath.Float32(f)
		f128 := eosMath.F32ToF128(f32)

		data, err := rlp.EncodeToBytes(f128)

		if err != nil {
			panic(err)
		}

		context.WriteMemory(ptr, data)
	}
}

func __extenddftf2(context Context) func(uint32, uint64) {
	return func(ptr uint32, f uint64) {
		f64 := eosMath.Float64(f)
		f128 := eosMath.F64ToF128(f64)

		data, err := rlp.EncodeToBytes(f128)

		if err != nil {
			panic(err)
		}

		context.WriteMemory(ptr, data)
	}
}

func __trunctfdf2(context Context) func(uint64, uint64) int64 {
	return func(low uint64, high uint64) int64 {
		f128 := eosMath.Float128{Low: low, High: high}
		f64 := eosMath.F128ToF64(f128)

		return int64(f64)
	}
}

func __trunctfsf2(context Context) func(uint64, uint64) int64 {
	return func(low uint64, high uint64) int64 {
		f128 := eosMath.Float128{Low: low, High: high}
		f32 := eosMath.F128ToF32(f128)

		return int64(f32)
	}
}

func __fixtfsi(context Context) func(uint64, uint64) int32 {
	return func(low uint64, high uint64) int32 {
		f128 := eosMath.Float128{Low: low, High: high}

		return eosMath.F128ToI32(f128, 0, false)
	}
}

func __fixtfdi(context Context) func(uint64, uint64) int64 {
	return func(low uint64, high uint64) int64 {
		f128 := eosMath.Float128{Low: low, High: high}

		return eosMath.F128ToI64(f128, 0, false)
	}
}

func __fixtfti(context Context) func(uint32, uint64, uint64) {
	return func(ptr uint32, low uint64, high uint64) {
		f128 := eosMath.Float128{Low: low, High: high}
		int128 := eosMath.Fixtfti(f128)
		data, err := rlp.EncodeToBytes(int128)

		if err != nil {
			panic(err)
		}

		context.WriteMemory(ptr, data)
	}
}

func __fixunstfsi(context Context) func(uint64, uint64) uint32 {
	return func(low uint64, high uint64) uint32 {
		f128 := eosMath.Float128{Low: low, High: high}

		return eosMath.F128ToUi32(f128, 0, false)
	}
}

func __fixunstfdi(context Context) func(uint64, uint64) uint64 {
	return func(low uint64, high uint64) uint64 {
		f128 := eosMath.Float128{Low: low, High: high}

		return eosMath.F128ToUi64(f128, 0, false)
	}
}

func __fixunstfti(context Context) func(uint32, uint64, uint64) {
	return func(ptr uint32, low uint64, high uint64) {
		f128 := eosMath.Float128{Low: low, High: high}
		uint128 := eosMath.Fixunstfti(f128)
		data, err := rlp.EncodeToBytes(uint128)

		if err != nil {
			panic(err)
		}

		context.WriteMemory(ptr, data)
	}
}

func __fixsfti(context Context) func(uint32, uint32) {
	return func(ptr uint32, value uint32) {
		int128 := eosMath.Fixsfti(value)
		data, err := rlp.EncodeToBytes(int128)

		if err != nil {
			panic(err)
		}

		context.WriteMemory(ptr, data)
	}
}

func __fixdfti(context Context) func(uint32, uint64) {
	return func(ptr uint32, value uint64) {
		int128 := eosMath.Fixdfti(value)
		data, err := rlp.EncodeToBytes(int128)

		if err != nil {
			panic(err)
		}

		context.WriteMemory(ptr, data)
	}
}

func __fixunssfti(context Context) func(uint32, uint32) {
	return func(ptr uint32, value uint32) {
		uint128 := eosMath.Fixunssfti(value)
		data, err := rlp.EncodeToBytes(uint128)

		if err != nil {
			panic(err)
		}

		context.WriteMemory(ptr, data)
	}
}

func __fixunsdfti(context Context) func(uint32, uint64) {
	return func(ptr uint32, value uint64) {
		uint128 := eosMath.Fixunsdfti(value)
		data, err := rlp.EncodeToBytes(uint128)

		if err != nil {
			panic(err)
		}

		context.WriteMemory(ptr, data)
	}
}

func __floatsidf(context Context) func(int32) int64 {
	return func(value int32) int64 {
		return int64(eosMath.I32ToF64(value))
	}
}

func __floatsitf(context Context) func(uint32, int32) {
	return func(ptr uint32, value int32) {
		data := eosMath.I32ToF128(value).Bytes()
		context.WriteMemory(ptr, data)
	}
}

func __floatditf(context Context) func(uint32, int64) {
	return func(ptr uint32, value int64) {
		data := eosMath.I64ToF128(value).Bytes()
		context.WriteMemory(ptr, data)
	}
}

func __floatunsitf(context Context) func(uint32, uint32) {
	return func(ptr uint32, value uint32) {
		data := eosMath.Ui32ToF128(value).Bytes()
		context.WriteMemory(ptr, data)
	}
}

func __floatunditf(context Context) func(uint32, uint64) {
	return func(ptr uint32, value uint64) {
		data := eosMath.Ui64ToF128(value).Bytes()
		context.WriteMemory(ptr, data)
	}
}

func __floattidf(context Context) func(uint64, uint64) uint64 {
	return func(low uint64, high uint64) uint64 {
		int128 := eosMath.Int128{Low: low, High: high}
		return math.Float64bits(eosMath.Floattidf(int128))
	}
}

func __floatuntidf(context Context) func(uint64, uint64) uint64 {
	return func(low uint64, high uint64) uint64 {
		uint128 := eosMath.Uint128{Low: low, High: high}
		return math.Float64bits(eosMath.Floatuntidf(uint128))
	}
}

func __eqtf2(context Context) func(int64, int64, int64, int64) int64 {
	return func(la int64, ha int64, lb int64, hb int64) int64 {
		return _cmptf2(la, ha, lb, hb, 1)
	}
}

func __netf2(context Context) func(int64, int64, int64, int64) int64 {
	return func(la int64, ha int64, lb int64, hb int64) int64 {
		return _cmptf2(la, ha, lb, hb, 1)
	}
}

func __getf2(context Context) func(int64, int64, int64, int64) int64 {
	return func(la int64, ha int64, lb int64, hb int64) int64 {
		return _cmptf2(la, ha, lb, hb, -1)
	}
}

func __gttf2(context Context) func(int64, int64, int64, int64) int64 {
	return func(la int64, ha int64, lb int64, hb int64) int64 {
		return _cmptf2(la, ha, lb, hb, 0)
	}
}

func __letf2(context Context) func(int64, int64, int64, int64) int64 {
	return func(la int64, ha int64, lb int64, hb int64) int64 {
		return _cmptf2(la, ha, lb, hb, 1)
	}
}

func __lttf2(context Context) func(int64, int64, int64, int64) int64 {
	return func(la int64, ha int64, lb int64, hb int64) int64 {
		return _cmptf2(la, ha, lb, hb, 0)
	}
}

func __cmptf2(context Context) func(int64, int64, int64, int64) int64 {
	return func(la int64, ha int64, lb int64, hb int64) int64 {
		return _cmptf2(la, ha, lb, hb, 1)
	}
}

func __unordtf2(context Context) func(int64, int64, int64, int64) int64 {
	return func(la int64, ha int64, lb int64, hb int64) int64 {
		a := eosMath.Float128{Low: uint64(la), High: uint64(ha)}
		b := eosMath.Float128{Low: uint64(lb), High: uint64(hb)}

		if a.IsNan() || b.IsNan() {
			return 1
		}

		return 0
	}
}

func _cmptf2(la, ha, lb, hb int64, return_value_if_nan int) int64 { //TODO unsame with regist
	a := eosMath.Float128{Low: uint64(la), High: uint64(ha)}
	b := eosMath.Float128{Low: uint64(lb), High: uint64(hb)}

	if _unordtf2(la, ha, lb, hb) != 0 {
		return int64(return_value_if_nan)
	}

	if a.F128Lt(b) {
		return -1
	}

	if a.F128EQ(b) {
		return 0
	}

	return 1
}

func _unordtf2(la, ha, lb, hb int64) int64 {
	a := eosMath.Float128{Low: uint64(la), High: uint64(ha)}
	b := eosMath.Float128{Low: uint64(lb), High: uint64(hb)}

	if a.IsNan() || b.IsNan() {
		return 1
	}

	return 0
}
