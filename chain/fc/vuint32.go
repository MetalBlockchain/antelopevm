package fc

import (
	"encoding/binary"
	"fmt"
)

//go:generate msgp
type UnsignedInt uint32

func (v UnsignedInt) Pack() ([]byte, error) {
	return WriteUVarInt(int(v)), nil
}
func (v *UnsignedInt) Unpack(in []byte) (l int, err error) {
	re, l, err := ReadUvarint64(in)
	if err != nil {
		return 0, nil
	}
	*v = UnsignedInt(re)
	return l, nil
}

func WriteUVarInt(v int) []byte {
	buf := make([]byte, 8)
	l := binary.PutUvarint(buf, uint64(v))
	return buf[:l]
}

func ReadUvarint64(in []byte) (uint64, int, error) {
	l, read := binary.Uvarint(in)
	if read < 0 {
		return l, 0, fmt.Errorf("too short")
	}

	return l, read, nil
}
