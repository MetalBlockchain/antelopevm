package name

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/hashicorp/go-set"
)

type Name uint64
type AccountName = Name
type PermissionName = Name
type ActionName = Name
type TableName = Name
type ScopeName = Name

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
	*n = Name(binary.LittleEndian.Uint64(data))
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

func (n Name) Hash() string {
	return n.String()
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
func StringToName(str string) Name {
	var n uint64
	for i := 0; i < len(str) && i < 12; i++ {
		n |= (charToSymbol(str[i]) & 0x1F) << uint(64-5*(i+1))
	}

	if len(str) > 12 {
		n |= charToSymbol(str[12]) & 0x0F
	}

	return Name(n)
}

var charmap = ".12345abcdefghijklmnopqrstuvwxyz"

// S converts a uint64 to a base32 string. String representation of the name.
func NameToString(value uint64) string {
	str := strings.Repeat(".", 13)

	tmp := value
	for i := uint32(0); i <= 12; i++ {
		var c byte
		if i == 0 {
			c = charmap[tmp&0x0F]
		} else {
			c = charmap[tmp&0x1F]
		}
		str = setCharAtIndex(str, 12-int(i), c)
		tmp >>= func() uint64 {
			if i == 0 {
				return 4
			}
			return 5
		}()
	}

	str = strings.TrimRight(str, ".")
	return str
}

func setCharAtIndex(s string, index int, c byte) string {
	if index < 0 || index >= len(s) {
		return s
	}
	chars := []byte(s)
	chars[index] = c
	return string(chars)
}

// charToSymbol converts a base32 symbol into its binary representation, used by N()
func charToSymbol(c byte) uint64 {
	if c >= 'a' && c <= 'z' {
		return uint64(c-'a') + 6
	}
	if c >= '1' && c <= '5' {
		return uint64(c-'1') + 1
	}
	return 0
}

var base32Alphabet = []byte(".12345abcdefghijklmnopqrstuvwxyz")

type NameSet = *set.HashSet[Name, string]

func NewNameSet(capacity int) NameSet {
	return set.NewHashSet[Name, string](capacity)
}
