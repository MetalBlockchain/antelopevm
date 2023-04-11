package abi

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/MetalBlockchain/antelopevm/math"
)

func (a *ContractAbi) DecodeAction(actionName name.ActionName, data []byte) ([]byte, error) {
	binaryDecoder := rlp.NewDecoder(data)
	action := a.ActionForName(actionName)

	if action == nil {
		return []byte{}, fmt.Errorf("action %s not found in abi", actionName)
	}

	return a.decode(binaryDecoder, action.Type)
}

func (a *ContractAbi) DecodeStruct(structType string, data []byte) ([]byte, error) {
	binaryDecoder := rlp.NewDecoder(data)
	return a.decode(binaryDecoder, structType)
}

func (a *ContractAbi) decode(binaryDecoder *rlp.Decoder, structName string) ([]byte, error) {
	structure := a.StructForName(structName)

	if structure == nil {
		return []byte{}, fmt.Errorf("structure [%s] not found in abi", structName)
	}

	resultingJson := make([]byte, 0)

	if structure.Base != "" {
		var err error
		resultingJson, err = a.decode(binaryDecoder, structure.Base)
		if err != nil {
			return resultingJson, fmt.Errorf("decode base [%s]: %s", structName, err)
		}
	}

	return a.decodeFields(binaryDecoder, structure.Fields, resultingJson)
}

func (a *ContractAbi) decodeFields(binaryDecoder *rlp.Decoder, fields []FieldDef, json []byte) ([]byte, error) {
	resultingJson := json
	for _, field := range fields {

		fieldType, isOptional, isArray := analyzeFieldType(field.Type)
		typeName := a.TypeNameForNewTypeName(fieldType)

		var err error
		resultingJson, err = a.decodeField(binaryDecoder, field.Name, typeName, isOptional, isArray, resultingJson)
		if err != nil {
			return []byte{}, fmt.Errorf("decoding fields: %s", err)
		}
	}

	return resultingJson, nil
}

func (a *ContractAbi) decodeField(binaryDecoder *rlp.Decoder, fieldName string, fieldType string, isOptional bool, isArray bool, json []byte) ([]byte, error) {
	resultingJson := json
	if isOptional {
		b, err := binaryDecoder.ReadByte()
		if err != nil {
			return resultingJson, fmt.Errorf("decoding field [%s] optional flag: %s", fieldName, err)
		}

		if b == 0 {
			return resultingJson, nil
		}
	}

	if isArray {
		length, err := binaryDecoder.ReadUvarint64()
		if err != nil {
			return resultingJson, fmt.Errorf("reading field [%s] array length: %s", fieldName, err)
		}

		if length == 0 {
			resultingJson, _ = SetBytes(resultingJson, fieldName, []interface{}{})
			//ignoring err because there is a bug in sjson. sjson shadow the err in case of a default type ...
			//if err != nil {
			//	return resultingJson, fmt.Errorf("reading field [%s] setting empty array: %s", fieldName, err)
			//}
		}

		for i := uint64(0); i < length; i++ {
			indexedFieldName := fmt.Sprintf("%s.%d", fieldName, i)
			resultingJson, err = a.read(binaryDecoder, indexedFieldName, fieldType, resultingJson)
			if err != nil {
				return resultingJson, fmt.Errorf("reading field [%s] index [%d]: %s", fieldName, i, err)
			}
		}

		return resultingJson, nil
	}

	resultingJson, err := a.read(binaryDecoder, fieldName, fieldType, resultingJson)

	if err != nil {
		return resultingJson, fmt.Errorf("decoding field [%s] of type [%s]: %s", fieldName, fieldType, err)
	}

	return resultingJson, nil
}

func (a *ContractAbi) read(binaryDecoder *rlp.Decoder, fieldName string, fieldType string, json []byte) ([]byte, error) {
	structure := a.StructForName(fieldType)

	if structure != nil {
		structureJson, err := a.decode(binaryDecoder, structure.Name)
		if err != nil {
			return []byte{}, err
		}
		return SetRawBytes(json, fieldName, structureJson)
	}

	var value interface{}
	var err error
	switch fieldType {
	case "int8":
		value, err = binaryDecoder.ReadInt8()
	case "uint8":
		value, err = binaryDecoder.ReadUint8()
	case "int16":
		value, err = binaryDecoder.ReadInt16()
	case "uint16":
		value, err = binaryDecoder.ReadUint16()
	case "int32":
		value, err = binaryDecoder.ReadInt32()
	case "uint32":
		value, err = binaryDecoder.ReadUint32()
	case "int64":
		var val int64
		val, err = binaryDecoder.ReadInt64()
		value = Int64(val)
	case "uint64":
		var val uint64
		val, err = binaryDecoder.ReadUint64()
		value = Uint64(val)
	case "int128":
		var data []byte
		data, err = binaryDecoder.ReadUint128("int128")
		int128 := math.Int128{
			Low:  binary.LittleEndian.Uint64(data),
			High: binary.LittleEndian.Uint64(data[8:]),
		}
		value = int128.String()
	case "uint128":
		var data []byte
		data, err = binaryDecoder.ReadUint128("uint128")
		uint128 := math.Uint128{
			Low:  binary.LittleEndian.Uint64(data),
			High: binary.LittleEndian.Uint64(data[8:]),
		}
		value = uint128.String()
	case "varint32":
		value, err = binaryDecoder.ReadVarint32()
	case "varuint32":
		value, err = binaryDecoder.ReadUvarint32()
	case "float32":
		value, err = binaryDecoder.ReadFloat32()
	case "float64":
		value, err = binaryDecoder.ReadFloat64()
	case "float128":
		var data []byte
		data, err = binaryDecoder.ReadUint128("float128")
		float128 := math.Float128{
			Low:  binary.LittleEndian.Uint64(data),
			High: binary.LittleEndian.Uint64(data[8:]),
		}
		value = float128.String()
	case "bool":
		value, err = binaryDecoder.ReadBool()
	case "time_point":
		var timePoint int64
		timePoint, err = binaryDecoder.ReadInt64()
		if err == nil {
			t := time.Unix(0, int64(timePoint*1000))
			value = t.UTC().Format("2006-01-02T15:04:05.999")
		}
	case "time_point_sec":
		var timePointSec uint32
		timePointSec, err = binaryDecoder.ReadUint32()
		if err == nil {
			t := time.Unix(int64(timePointSec), 0)
			value = t.UTC().Format("2006-01-02T15:04:05")
		}
	case "block_timestamp_type":
		var slot uint32
		slot, err = binaryDecoder.ReadUint32()
		if err == nil {
			value = core.BlockTimeStamp(slot).ToTimePoint(500, 0).String()
		}
	case "name":
		var val uint64
		val, err = binaryDecoder.ReadName() //uint64
		value = name.NameToString(val)
	case "bytes":
		value, err = binaryDecoder.ReadByteArray()
		if err == nil {
			value = hex.EncodeToString(value.([]byte))
		}
	case "string":
		value, err = binaryDecoder.ReadString()
	case "checksum160":
		val, e := binaryDecoder.ReadChecksum160() //[]byte
		if e == nil {
			value = crypto.NewRipemd160Byte(val)
		}
		err = e
	case "checksum256":
		val, e := binaryDecoder.ReadChecksum256() //[]byte
		if e == nil {
			value = crypto.NewSha256Byte(val)
		}
		err = e
	case "checksum512":
		val, e := binaryDecoder.ReadChecksum512() //[]byte
		if e == nil {
			value = crypto.NewSha512Byte(val)
		}
		err = e
	case "public_key":
		var pubKey ecc.PublicKey
		err = binaryDecoder.Decode(&pubKey)
		if err == nil {
			value = pubKey
		}
	case "signature":
		var sig ecc.Signature
		err = binaryDecoder.Decode(&sig)
		if err == nil {
			value = sig
		}
	case "symbol":
		s := core.Symbol{}
		err := binaryDecoder.Decode(&s)
		if err == nil {
			value = fmt.Sprintf("%d,%s", s.Precision, s.Symbol)
		}
	case "symbol_code":
		var data uint64
		data, err = binaryDecoder.ReadUint64()
		value = core.SymbolCode(data)
	case "asset":
		a := core.Asset{}
		err = binaryDecoder.Decode(&a)
		if err == nil {
			value = a
		}
	case "extended_asset":
		e := core.ExtendedAsset{}
		err = binaryDecoder.Decode(&e)
		if err == nil {
			value = e
		}
	default:
		return nil, fmt.Errorf("read field of type [%s]: unknown type", fieldType)
	}

	if err != nil {
		return []byte{}, fmt.Errorf("read: %s", err)
	}

	return SetBytes(json, fieldName, value)
}

func analyzeFieldType(fieldType string) (typeName string, isOptional bool, isArray bool) {
	if strings.HasSuffix(fieldType, "?") {
		return fieldType[0 : len(fieldType)-1], true, false
	}

	if strings.HasSuffix(fieldType, "$") {
		return fieldType[0 : len(fieldType)-1], false, false
	}

	if strings.HasSuffix(fieldType, "[]") {
		return fieldType[0 : len(fieldType)-2], false, true
	}

	return fieldType, false, false
}
