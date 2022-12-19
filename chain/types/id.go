package types

import "encoding/binary"

type IdType uint64

func (p *IdType) ToBytes() []byte {
	a := make([]byte, 8)
	binary.BigEndian.PutUint64(a, uint64(*p))
	return a
}
