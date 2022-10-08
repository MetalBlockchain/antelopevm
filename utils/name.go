package utils

type Name struct {
	Value string
}

func Char(c byte) uint64 {
	if c >= 'a' && c <= 'z' {
		return uint64((c - 'a') + 6)
	}
	if c >= '1' && c <= '5' {
		return uint64((c - '1') + 1)
	}

	return 0
}

func (n *Name) ToUint64() uint64 {
	var v uint64

	for i := len(n.Value) - 1; i >= 0; i-- {
		v |= Char(n.Value[i]) << (64 - 5*(i+1))
	}

	return v
}

func (n *Name) ToInt64() int64 {
	return int64(n.ToUint64())
}
