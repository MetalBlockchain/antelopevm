package api

import (
	"math"

	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	eosMath "github.com/MetalBlockchain/antelopevm/math"
)

func init() {
	Functions["__ashlti3"] = __ashlti3
	Functions["__ashrti3"] = __ashrti3
	Functions["__lshlti3"] = __lshlti3
	Functions["__lshrti3"] = __lshrti3
	Functions["__divti3"] = __divti3
	Functions["__udivti3"] = __udivti3
	Functions["__multi3"] = __multi3
	Functions["__modti3"] = __modti3
	Functions["__umodti3"] = __umodti3
	Functions["__addtf3"] = __addtf3
	Functions["__subtf3"] = __subtf3
	Functions["__multf3"] = __multf3
	Functions["__divtf3"] = __divtf3
	Functions["__negtf2"] = __negtf2
	Functions["__extendsftf2"] = __extendsftf2
	Functions["__extenddftf2"] = __extenddftf2
	Functions["__trunctfdf2"] = __trunctfdf2
	Functions["__trunctfsf2"] = __trunctfsf2
	Functions["__fixtfsi"] = __fixtfsi
	Functions["__fixtfdi"] = __fixtfdi
	Functions["__fixtfti"] = __fixtfti
	Functions["__fixunstfsi"] = __fixunstfsi
	Functions["__fixunstfdi"] = __fixunstfdi
	Functions["__fixunstfti"] = __fixunstfti
	Functions["__fixsfti"] = __fixsfti
	Functions["__fixdfti"] = __fixdfti
	Functions["__fixunssfti"] = __fixunssfti
	Functions["__fixunsdfti"] = __fixunsdfti
	Functions["__floatsidf"] = __floatsidf
	Functions["__floatsitf"] = __floatsitf
	Functions["__floatditf"] = __floatditf
	Functions["__floatunsitf"] = __floatunsitf
	Functions["__floatunditf"] = __floatunditf
	Functions["__floattidf"] = __floattidf
	Functions["__floatuntidf"] = __floatuntidf
	Functions["__eqtf2"] = __eqtf2
	Functions["__netf2"] = __netf2
	Functions["__getf2"] = __getf2
	Functions["__gttf2"] = __gttf2
	Functions["__letf2"] = __letf2
	Functions["__lttf2"] = __lttf2
	Functions["__cmptf2"] = __cmptf2
	Functions["__unordtf2"] = __unordtf2
}

func __ashlti3(context Context) interface{} {
	return func(ptr uint32, low uint64, high uint64, shift uint32) {
		i := eosMath.Int128{Low: low, High: high}
		i.LeftShifts(int(shift))
		data, _ := rlp.EncodeToBytes(i)
		context.WriteMemory(ptr, data)
	}
}

func __ashrti3(context Context) interface{} {
	return func(ptr uint32, low uint64, high uint64, shift uint32) {
		i := eosMath.Int128{Low: low, High: high}
		i.RightShifts(int(shift))
		data, _ := rlp.EncodeToBytes(i)
		context.WriteMemory(ptr, data)
	}
}

func __lshlti3(context Context) interface{} {
	return func(ptr uint32, low uint64, high uint64, shift uint32) {
		i := eosMath.Int128{Low: low, High: high}
		i.LeftShifts(int(shift))
		data, _ := rlp.EncodeToBytes(i)
		context.WriteMemory(ptr, data)
	}
}

func __lshrti3(context Context) interface{} {
	return func(ptr uint32, low uint64, high uint64, shift uint32) {
		i := eosMath.Int128{Low: low, High: high}
		i.RightShifts(int(shift))
		data, _ := rlp.EncodeToBytes(i)
		context.WriteMemory(ptr, data)
	}
}

func __divti3(context Context) interface{} {
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

func __udivti3(context Context) interface{} {
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

func __multi3(context Context) interface{} {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		lhs := eosMath.Int128{Low: la, High: ha}
		rhs := eosMath.Int128{Low: lb, High: hb}
		data, _ := rlp.EncodeToBytes(lhs.Mul(rhs))
		context.WriteMemory(ptr, data)
	}
}

func __modti3(context Context) interface{} {
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

func __umodti3(context Context) interface{} {
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

func __addtf3(context Context) interface{} {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		a := eosMath.Float128{Low: la, High: ha}
		b := eosMath.Float128{Low: lb, High: hb}

		data, _ := rlp.EncodeToBytes(a.Add(b))
		context.WriteMemory(ptr, data)
	}
}

func __subtf3(context Context) interface{} {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		a := eosMath.Float128{Low: la, High: ha}
		b := eosMath.Float128{Low: lb, High: hb}

		data, _ := rlp.EncodeToBytes(a.Sub(b))
		context.WriteMemory(ptr, data)
	}
}

func __multf3(context Context) interface{} {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		a := eosMath.Float128{Low: la, High: ha}
		b := eosMath.Float128{Low: lb, High: hb}

		data, _ := rlp.EncodeToBytes(a.Mul(b))
		context.WriteMemory(ptr, data)
	}
}

func __divtf3(context Context) interface{} {
	return func(ptr uint32, la uint64, ha uint64, lb uint64, hb uint64) {
		a := eosMath.Float128{Low: la, High: ha}
		b := eosMath.Float128{Low: lb, High: hb}

		data, _ := rlp.EncodeToBytes(a.Div(b))
		context.WriteMemory(ptr, data)
	}
}

func __negtf2(context Context) interface{} {
	return func(ptr uint32, la uint64, ha uint64) {
		high := uint64(ha)
		high ^= uint64(1) << 63
		f128 := eosMath.Float128{Low: uint64(la), High: high}

		data, _ := rlp.EncodeToBytes(f128)
		context.WriteMemory(ptr, data)
	}
}

func __extendsftf2(context Context) interface{} {
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

func __extenddftf2(context Context) interface{} {
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

func __trunctfdf2(context Context) interface{} {
	return func(low uint64, high uint64) int64 {
		f128 := eosMath.Float128{Low: low, High: high}
		f64 := eosMath.F128ToF64(f128)

		return int64(f64)
	}
}

func __trunctfsf2(context Context) interface{} {
	return func(low uint64, high uint64) int64 {
		f128 := eosMath.Float128{Low: low, High: high}
		f32 := eosMath.F128ToF32(f128)

		return int64(f32)
	}
}

func __fixtfsi(context Context) interface{} {
	return func(low uint64, high uint64) int32 {
		f128 := eosMath.Float128{Low: low, High: high}

		return eosMath.F128ToI32(f128, 0, false)
	}
}

func __fixtfdi(context Context) interface{} {
	return func(low uint64, high uint64) int64 {
		f128 := eosMath.Float128{Low: low, High: high}

		return eosMath.F128ToI64(f128, 0, false)
	}
}

func __fixtfti(context Context) interface{} {
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

func __fixunstfsi(context Context) interface{} {
	return func(low uint64, high uint64) uint32 {
		f128 := eosMath.Float128{Low: low, High: high}

		return eosMath.F128ToUi32(f128, 0, false)
	}
}

func __fixunstfdi(context Context) interface{} {
	return func(low uint64, high uint64) uint64 {
		f128 := eosMath.Float128{Low: low, High: high}

		return eosMath.F128ToUi64(f128, 0, false)
	}
}

func __fixunstfti(context Context) interface{} {
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

func __fixsfti(context Context) interface{} {
	return func(ptr uint32, value uint32) {
		int128 := eosMath.Fixsfti(value)
		data, err := rlp.EncodeToBytes(int128)

		if err != nil {
			panic(err)
		}

		context.WriteMemory(ptr, data)
	}
}

func __fixdfti(context Context) interface{} {
	return func(ptr uint32, value uint64) {
		int128 := eosMath.Fixdfti(value)
		data, err := rlp.EncodeToBytes(int128)

		if err != nil {
			panic(err)
		}

		context.WriteMemory(ptr, data)
	}
}

func __fixunssfti(context Context) interface{} {
	return func(ptr uint32, value uint32) {
		uint128 := eosMath.Fixunssfti(value)
		data, err := rlp.EncodeToBytes(uint128)

		if err != nil {
			panic(err)
		}

		context.WriteMemory(ptr, data)
	}
}

func __fixunsdfti(context Context) interface{} {
	return func(ptr uint32, value uint64) {
		uint128 := eosMath.Fixunsdfti(value)
		data, err := rlp.EncodeToBytes(uint128)

		if err != nil {
			panic(err)
		}

		context.WriteMemory(ptr, data)
	}
}

func __floatsidf(context Context) interface{} {
	return func(value int32) int64 {
		return int64(eosMath.I32ToF64(value))
	}
}

func __floatsitf(context Context) interface{} {
	return func(ptr uint32, value int32) {
		data := eosMath.I32ToF128(value).Bytes()
		context.WriteMemory(ptr, data)
	}
}

func __floatditf(context Context) interface{} {
	return func(ptr uint32, value int64) {
		data := eosMath.I64ToF128(value).Bytes()
		context.WriteMemory(ptr, data)
	}
}

func __floatunsitf(context Context) interface{} {
	return func(ptr uint32, value uint32) {
		data := eosMath.Ui32ToF128(value).Bytes()
		context.WriteMemory(ptr, data)
	}
}

func __floatunditf(context Context) interface{} {
	return func(ptr uint32, value uint64) {
		data := eosMath.Ui64ToF128(value).Bytes()
		context.WriteMemory(ptr, data)
	}
}

func __floattidf(context Context) interface{} {
	return func(low uint64, high uint64) uint64 {
		int128 := eosMath.Int128{Low: low, High: high}
		return math.Float64bits(eosMath.Floattidf(int128))
	}
}

func __floatuntidf(context Context) interface{} {
	return func(low uint64, high uint64) uint64 {
		uint128 := eosMath.Uint128{Low: low, High: high}
		return math.Float64bits(eosMath.Floatuntidf(uint128))
	}
}

func __eqtf2(context Context) interface{} {
	return func(la int64, ha int64, lb int64, hb int64) int64 {
		return _cmptf2(la, ha, lb, hb, 1)
	}
}

func __netf2(context Context) interface{} {
	return func(la int64, ha int64, lb int64, hb int64) int64 {
		return _cmptf2(la, ha, lb, hb, 1)
	}
}

func __getf2(context Context) interface{} {
	return func(la int64, ha int64, lb int64, hb int64) int64 {
		return _cmptf2(la, ha, lb, hb, -1)
	}
}

func __gttf2(context Context) interface{} {
	return func(la int64, ha int64, lb int64, hb int64) int64 {
		return _cmptf2(la, ha, lb, hb, 0)
	}
}

func __letf2(context Context) interface{} {
	return func(la int64, ha int64, lb int64, hb int64) int64 {
		return _cmptf2(la, ha, lb, hb, 1)
	}
}

func __lttf2(context Context) interface{} {
	return func(la int64, ha int64, lb int64, hb int64) int64 {
		return _cmptf2(la, ha, lb, hb, 0)
	}
}

func __cmptf2(context Context) interface{} {
	return func(la int64, ha int64, lb int64, hb int64) int64 {
		return _cmptf2(la, ha, lb, hb, 1)
	}
}

func __unordtf2(context Context) interface{} {
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
