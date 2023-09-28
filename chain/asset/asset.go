package asset

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
)

type SymbolCode = uint64
type Symbol struct {
	Precision uint8
	Symbol    string
}

func (s Symbol) Pack() (re []byte, err error) {
	symbol := make([]byte, 7, 7)
	copy(symbol[:], []byte(s.Symbol))

	re = append(re, byte(s.Precision))
	re = append(re, symbol...)
	return re, nil
}

func (s *Symbol) Unpack(in []byte) (int, error) {
	if len(in) < 8 {
		return 0, fmt.Errorf("asset symbol required [%d] bytes, remaining [%d]", 7, len(in))
	}
	s.Precision = uint8(in[0])
	s.Symbol = strings.TrimRight(string(in[1:8]), "\x00")
	return 8, nil
}

func (sym *Symbol) Name() string {
	return sym.Symbol
}

func StringToSymbol(precision uint8, str string) (uint64, error) {
	var result uint64
	len := uint32(len(str))

	for i := uint32(0); i < len; i++ {
		// All characters must be upper case alphabets
		if str[i] < 'A' || str[i] > 'Z' {
			return 0, fmt.Errorf("invalid character in symbol name")
		}

		result |= uint64(str[i]) << (8 * (i + 1))
	}

	result |= uint64(precision)

	return result, nil
}

type Asset struct {
	Amount int64 `eos:"asset"`
	Symbol
}

func (s *Asset) Unpack(in []byte) (int, error) {
	decoder := rlp.NewDecoder(in)
	a, err := decoder.ReadInt64()

	if err != nil {
		return 0, err
	}

	s.Amount = a
	l, err := s.Symbol.Unpack(decoder.GetData()[decoder.GetPos():])

	return l + decoder.GetPos(), err
}

func (a Asset) String() string {
	sign := ""
	abs := a.Amount

	if a.Amount < 0 {
		sign = "-"
		abs = -1 * a.Amount
	}

	strInt := fmt.Sprintf("%d", abs)

	if len(strInt) < int(a.Symbol.Precision+1) {
		// prepend `0` for the difference:
		strInt = strings.Repeat("0", int(a.Symbol.Precision+uint8(1))-len(strInt)) + strInt
	}

	var result string

	if a.Symbol.Precision == 0 {
		result = strInt
	} else {
		result = strInt[:len(strInt)-int(a.Symbol.Precision)] + "." + strInt[len(strInt)-int(a.Symbol.Precision):]
	}

	return fmt.Sprintf("%s %s", sign+result, a.Symbol.Symbol)
}

func (a Asset) MarshalJSON() (data []byte, err error) {
	return json.Marshal(a.String())
}

type ExtendedAsset struct {
	Asset    Asset `json:"asset"`
	Contract name.AccountName
}
