package rlp_test

import (
	"crypto/rand"
	"testing"

	"github.com/MetalBlockchain/antelopevm/chain"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/stretchr/testify/assert"
)

func BenchmarkXxx(b *testing.B) {
	buf := make([]byte, 2048)
	rand.Read(buf)
	obj := chain.SetCode{
		Account:   name.StringToName("eosio.token"),
		VmType:    0,
		VmVersion: 0,
		Code:      buf,
	}

	for i := 0; i < b.N; i++ {
		rlp.EncodeToBytes(obj)
	}
}

func BenchmarkDecode(b *testing.B) {
	buf := make([]byte, 2048)
	rand.Read(buf)
	obj := chain.SetCode{
		Account:   name.StringToName("eosio.token"),
		VmType:    0,
		VmVersion: 0,
		Code:      buf,
	}
	data, _ := rlp.EncodeToBytes(obj)

	for i := 0; i < b.N; i++ {
		o := chain.SetCode{}
		rlp.DecodeBytes(data, &o)
	}
}

func TestXXX(b *testing.T) {
	buf := make([]byte, 2048)
	rand.Read(buf)
	obj := chain.SetCode{
		Account:   name.StringToName("eosio.token"),
		VmType:    0,
		VmVersion: 0,
		Code:      buf,
	}

	data, _ := rlp.EncodeToBytes(obj)

	obj2 := chain.SetCode{}
	err := rlp.DecodeBytes(data, &obj2)
	assert.NoError(b, err)
	assert.Equal(b, obj.Code, obj2.Code)
}
