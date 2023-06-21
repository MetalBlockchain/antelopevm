package utils

import "encoding/binary"

func Uint64ToLittleEndian(value uint64) []byte {
	a := make([]byte, 8)
	binary.LittleEndian.PutUint64(a, uint64(value))
	return a
}
