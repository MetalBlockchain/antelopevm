package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var symbolRegex = regexp.MustCompile("^[0-9]{1,2},[A-Z]{1,7}$")
var symbolCodeRegex = regexp.MustCompile("^[A-Z]{1,7}$")

type Symbol struct {
	Precision uint8
	Symbol    string

	// Caching of symbol code if it was computed once
	symbolCode uint64
}

func StringToSymbol(str string) (Symbol, error) {
	symbol := Symbol{}
	if !symbolRegex.MatchString(str) {
		return symbol, fmt.Errorf("%s is not a valid symbol", str)
	}
	arrs := strings.Split(str, ",")
	precision, _ := strconv.ParseUint(string(arrs[0]), 10, 8)

	symbol.Precision = uint8(precision)
	symbol.Symbol = arrs[1]

	return symbol, nil
}

func (s Symbol) SymbolCode() (SymbolCode, error) {
	if s.symbolCode != 0 {
		return SymbolCode(s.symbolCode), nil
	}

	symbolCode, err := StringToSymbolCode(s.Symbol)
	if err != nil {
		return 0, err
	}

	return SymbolCode(symbolCode), nil
}

func (s Symbol) MustSymbolCode() SymbolCode {
	symbolCode, err := StringToSymbolCode(s.Symbol)
	if err != nil {
		panic("invalid symbol code " + s.Symbol)
	}

	return symbolCode
}

func (s Symbol) ToUint64() (uint64, error) {
	symbolCode, err := s.SymbolCode()
	if err != nil {
		return 0, fmt.Errorf("symbol %s is not a valid symbol code: %w", s.Symbol, err)
	}

	return uint64(symbolCode)<<8 | uint64(s.Precision), nil
}

func (s Symbol) ToName() (string, error) {
	u, err := s.ToUint64()
	if err != nil {
		return "", err
	}
	return NameToString(u), nil
}

func (s Symbol) String() string {
	return fmt.Sprintf("%d,%s", s.Precision, s.Symbol)
}

type SymbolCode uint64

func NameToSymbolCode(name Name) (SymbolCode, error) {
	value, err := StringToName(string(name))
	if err != nil {
		return 0, fmt.Errorf("name %s is invalid: %w", name, err)
	}

	return SymbolCode(value), nil
}

func StringToSymbolCode(str string) (SymbolCode, error) {
	if len(str) > 7 {
		return 0, fmt.Errorf("string is too long to be a valid symbol_code")
	}

	var symbolCode uint64
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] < 'A' || str[i] > 'Z' {
			return 0, fmt.Errorf("only uppercase letters allowed in symbol_code string")
		}

		symbolCode <<= 8
		symbolCode = symbolCode | uint64(str[i])
	}

	return SymbolCode(symbolCode), nil
}

func (sc SymbolCode) ToName() string {
	return NameToString(uint64(sc))
}

func (sc SymbolCode) String() string {
	builder := strings.Builder{}

	symbolCode := uint64(sc)
	for i := 0; i < 7; i++ {
		if symbolCode == 0 {
			return builder.String()
		}

		builder.WriteByte(byte(symbolCode & 0xFF))
		symbolCode >>= 8
	}

	return builder.String()
}

func (sc SymbolCode) MarshalJSON() (data []byte, err error) {
	return []byte(`"` + sc.String() + `"`), nil
}
