package utils

import "encoding/binary"

type Increment uint64

func IncrementFromBytes(data []byte) Increment {
	val := binary.BigEndian.Uint64(data)

	return Increment(val)
}

func (i *Increment) Increment() {
	*i += 1
}

func (i *Increment) ToBytes() []byte {
	a := make([]byte, 8)
	binary.BigEndian.PutUint64(a, uint64(*i))
	return a
}
