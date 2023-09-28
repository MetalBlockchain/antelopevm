package state

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/MetalBlockchain/antelopevm/chain/block"
	"github.com/MetalBlockchain/antelopevm/chain/entity"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/time"
	"github.com/MetalBlockchain/antelopevm/chain/transaction"
	"github.com/MetalBlockchain/antelopevm/chain/types"
)

func getPartialKey(index string, obj entity.Entity, values ...interface{}) []byte {
	key := []byte{obj.GetObjectType()}
	key = append(key, []byte("__"+index)...)

	for _, value := range values {
		key = append(key, []byte("__")...)
		key = append(key, encodeType(value)...)
	}

	return key
}

func getObjectKeyByIndex(obj entity.Entity, indexName string) []byte {
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

func getObjectKeys(obj entity.Entity) map[string][]byte {
	fields := obj.GetIndexes()
	keys := make(map[string][]byte)

	for index, _ := range fields {
		keys[index] = getObjectKeyByIndex(obj, index)
	}

	return keys
}

func encodeType(obj interface{}) []byte {
	switch v := obj.(type) {
	case types.IdType:
		return v.ToBytes()
	case transaction.TransactionIdType:
		return v.Bytes()
	case name.Name:
		return v.Pack()
	case block.BlockHash:
		return v[:]
	case uint64:
		return types.IdType(v).ToBytes()
	case bool:
		if v {
			return []byte{1}
		}

		return []byte{0}
	case time.TimePointSec:
		a := make([]byte, 4)
		binary.BigEndian.PutUint32(a, uint32(v))
		return a
	default:
		panic(fmt.Sprintf("type %v not supported", reflect.TypeOf(obj)))
	}
}
