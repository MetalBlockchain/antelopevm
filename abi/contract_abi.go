package abi

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
)

type ContractAbi struct {
	Version          string            `json:"version"`
	Types            []TypeDef         `json:"types,omitempty"`
	Structs          []StructDef       `json:"structs,omitempty"`
	Actions          []ActionDef       `json:"actions,omitempty"`
	Tables           []TableDef        `json:"tables,omitempty"`
	RicardianClauses []ClausePair      `json:"ricardian_clauses,omitempty"`
	ErrorMessages    []ErrorMessage    `json:"error_messages,omitempty"`
	Extensions       []core.Extension  `json:"abi_extensions,omitempty"`
	Variants         []VariantDef      `json:"variants,omitempty"`
	ActionResults    []ActionResultDef `json:"-"`
}

func NewABI(data []byte) (*ContractAbi, error) {
	abi := &ContractAbi{}
	decoder := rlp.NewDecoder(data)
	err := decoder.Decode(abi)

	if err != nil {
		return nil, fmt.Errorf("failed to parse abi: %s", err)
	}

	return abi, nil
}

func (a *ContractAbi) ActionForName(name core.ActionName) *ActionDef {
	for _, a := range a.Actions {
		if a.Name == name {
			return &a
		}
	}

	return nil
}

func (a *ContractAbi) StructForName(name string) *StructDef {
	for _, s := range a.Structs {
		if s.Name == name {
			return &s
		}
	}

	return nil
}

func (a *ContractAbi) TableForName(name core.TableName) *TableDef {
	for _, s := range a.Tables {
		if s.Name == name {
			return &s
		}
	}

	return nil
}

func (a *ContractAbi) TypeNameForNewTypeName(typeName string) string {
	for _, t := range a.Types {
		if t.NewTypeName == typeName {
			return t.Type
		}
	}

	return typeName
}

type TypeDef struct {
	NewTypeName string `json:"new_type_name"`
	Type        string `json:"type"`
}

type StructDef struct {
	Name   string     `json:"name"`
	Base   string     `json:"base"`
	Fields []FieldDef `json:"fields,omitempty"`
}

type FieldDef struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ActionDef struct {
	Name              core.ActionName `json:"name"`
	Type              string          `json:"type"`
	RicardianContract string          `json:"ricardian_contract"`
}

type TableDef struct {
	Name      core.TableName `json:"name"`
	IndexType string         `json:"index_type"`
	KeyNames  []string       `json:"key_names,omitempty"`
	KeyTypes  []string       `json:"key_types,omitempty"`
	Type      string         `json:"type"`
}

type ClausePair struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

type ErrorMessage struct {
	Code    uint64 `json:"error_code"`
	Message string `json:"error_msg"`
}

type VariantDef struct {
	Name  string   `json:"name"`
	Types []string `json:"types"`
}

type ActionResultDef struct {
	Name       string `json:"name"`
	ResultType string `json:"types"`
}

type Int64 int64

func (i Int64) MarshalJSON() (data []byte, err error) {
	if i > 0xffffffff || i < -0xffffffff {
		encodedInt, err := json.Marshal(int64(i))
		if err != nil {
			return nil, err
		}
		data = append([]byte{'"'}, encodedInt...)
		data = append(data, '"')
		return data, nil
	}
	return json.Marshal(int64(i))
}

func (i *Int64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty value")
	}

	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}

		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}

		*i = Int64(val)

		return nil
	}

	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = Int64(v)

	return nil
}

type Uint64 uint64

func (i Uint64) MarshalJSON() (data []byte, err error) {
	if i > 0xffffffff {
		encodedInt, err := json.Marshal(uint64(i))
		if err != nil {
			return nil, err
		}
		data = append([]byte{'"'}, encodedInt...)
		data = append(data, '"')
		return data, nil
	}
	return json.Marshal(uint64(i))
}

func (i *Uint64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty value")
	}

	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}

		val, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}

		*i = Uint64(val)

		return nil
	}

	var v uint64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = Uint64(v)

	return nil
}

type Uint128 struct {
	Lo uint64
	Hi uint64
}

type Int128 Uint128

type Float128 Uint128

func (i Uint128) MarshalJSON() (data []byte, err error) {
	return json.Marshal(i.String())
}

func (i Int128) MarshalJSON() (data []byte, err error) {
	return json.Marshal(Uint128(i).String())
}

func (i Float128) MarshalJSON() (data []byte, err error) {
	return json.Marshal(Uint128(i).String())
}

func (i Uint128) String() string {
	// Same for Int128, Float128
	number := make([]byte, 16)
	binary.LittleEndian.PutUint64(number[:], i.Lo)
	binary.LittleEndian.PutUint64(number[8:], i.Hi)
	return fmt.Sprintf("0x%s%s", hex.EncodeToString(number[:8]), hex.EncodeToString(number[8:]))
}

func (i *Int128) UnmarshalJSON(data []byte) error {
	var el Uint128
	if err := json.Unmarshal(data, &el); err != nil {
		return err
	}

	out := Int128(el)
	*i = out

	return nil
}

func (i *Float128) UnmarshalJSON(data []byte) error {
	var el Uint128
	if err := json.Unmarshal(data, &el); err != nil {
		return err
	}

	out := Float128(el)
	*i = out

	return nil
}

func (i *Uint128) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if !strings.HasPrefix(s, "0x") && !strings.HasPrefix(s, "0X") {
		return fmt.Errorf("int128 expects 0x prefix")
	}

	truncatedVal := s[2:]
	if len(truncatedVal) != 32 {
		return fmt.Errorf("int128 expects 32 characters after 0x, had %d", len(truncatedVal))
	}

	loHex := truncatedVal[:16]
	hiHex := truncatedVal[16:]

	lo, err := hex.DecodeString(loHex)
	if err != nil {
		return err
	}

	hi, err := hex.DecodeString(hiHex)
	if err != nil {
		return err
	}

	loUint := binary.LittleEndian.Uint64(lo)
	hiUint := binary.LittleEndian.Uint64(hi)

	i.Lo = loUint
	i.Hi = hiUint

	return nil
}
