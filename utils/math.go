package utils

func Min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}

func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func MinUint32(x, y uint32) uint32 {
	if x < y {
		return x
	}
	return y
}

func AbsInt32(value int32) int32 {
	if value < 0 {
		return -value
	}

	return value
}

func Max(x, y uint64) uint64 {
	if x > y {
		return x
	}
	return y
}
