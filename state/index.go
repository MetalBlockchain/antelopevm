package state

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/name"
)

func getPartialKey(index string, obj core.Entity, values ...interface{}) []byte {
	key := []byte{obj.GetObjectType()}
	key = append(key, []byte("__"+index)...)

	for _, value := range values {
		key = append(key, []byte("__")...)
		key = append(key, encodeType(value)...)
	}

	return key
}

func getObjectKeyByIndex(obj core.Entity, indexName string) []byte {
	fields := obj.GetIndexes()

	if index, found := fields[indexName]; found {
		indexKey := []byte{obj.GetObjectType()}
		indexKey = append(indexKey, []byte("__"+indexName)...)

		for _, field := range index.Fields {
			r := reflect.ValueOf(obj)
			f := reflect.Indirect(r).FieldByName(field).Interface()
			indexKey = append(indexKey, []byte("__")...)
			indexKey = append(indexKey, encodeType(f)...)
		}

		return indexKey
	}

	panic("invalid index")
}

func getObjectKeys(obj core.Entity) map[string][]byte {
	fields := obj.GetIndexes()
	keys := make(map[string][]byte)

	for index, _ := range fields {
		keys[index] = getObjectKeyByIndex(obj, index)
	}

	return keys
}

func encodeType(obj interface{}) []byte {
	switch v := obj.(type) {
	case core.IdType:
		return v.ToBytes()
	case core.TransactionIdType:
		return v.Bytes()
	case name.Name:
		return v.Pack()
	case core.BlockHash:
		return v[:]
	case uint64:
		return core.IdType(v).ToBytes()
	case bool:
		if v {
			return []byte{1}
		}

		return []byte{0}
	case core.TimePointSec:
		a := make([]byte, 4)
		binary.BigEndian.PutUint32(a, uint32(v))
		return a
	default:
		panic(fmt.Sprintf("type %v not supported", reflect.TypeOf(obj)))
	}
}
