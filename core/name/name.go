package name

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
)

//go:generate msgp
type Name uint64

func (n Name) String() string {
	return NameToString(uint64(n))
}

func (n Name) IsEmpty() bool {
	return n == 0
}

func (n Name) Pack() []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(n))
	return buf
}

func (n *Name) Unpack(in []byte) (rlp.Unpack, error) {
	if len(in) < 8 {
		return nil, fmt.Errorf("rlp: uint64 required a number of [%d] bytes, remaining [%d]", 8, len(in))
	}

	data := in[:8]
	out := binary.LittleEndian.Uint64(data)
	fmt.Println(Name(out))
	return nil, nil
}

var TypeName = reflect.TypeOf(Name(0))

func CompareName(first interface{}, second interface{}) int {
	if first.(Name) == second.(Name) {
		return 0
	}
	if first.(Name) < second.(Name) {
		return -1
	}
	return 1
}

func (n Name) Empty() bool {
	return n == 0
}

func (n Name) MarshalJSON() ([]byte, error) {
	return json.Marshal(NameToString(uint64(n)))
}

func (n *Name) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*n = StringToName(s)
	return nil
}

// N converts a base32 string to a uint64. 64-bit unsigned integer representation of the name.
func StringToName(s string) Name {
	var i uint32
	var val uint64
	sLen := uint32(len(s))
	if sLen > 13 {
		return Name(0)
	}
	for ; i <= 12; i++ {
		var c uint64
		if i < sLen {
			c = uint64(charToSymbol(s[i]))
		}

		if i < 12 {
			c &= 0x1f
			c <<= 64 - 5*(i+1)
		} else {
			c &= 0x0f
		}

		val |= c
	}

	return Name(val)
}

// S converts a uint64 to a base32 string. String representation of the name.
func NameToString(value uint64) string {
	a := []byte{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.'}

	tmp := value
	for i := 0; i <= 12; i++ {
		bit := 0x1f

		if i == 0 {
			bit = 0x0f
		}

		c := base32Alphabet[tmp&uint64(bit)]
		a[12-i] = c

		shift := uint(5)

		if i == 0 {
			shift = 4
		}

		tmp >>= shift
	}

	return strings.TrimRight(string(a), ".")
}

// charToSymbol converts a base32 symbol into its binary representation, used by N()
func charToSymbol(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - 'a' + 6
	}
	if c >= '1' && c <= '5' {
		return c - '1' + 1
	}
	return 0
}

var base32Alphabet = []byte(".12345abcdefghijklmnopqrstuvwxyz")
