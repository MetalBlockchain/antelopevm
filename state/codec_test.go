package state

import (
	"fmt"
	"testing"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/account"
	"github.com/MetalBlockchain/antelopevm/core/name"
)

func BenchmarkMarshal(b *testing.B) {
	account := &account.Account{
		ID:             core.IdType(102001002),
		Name:           name.StringToName("glenn"),
		VmType:         1,
		VmVersion:      1,
		Privileged:     false,
		LastCodeUpdate: core.Now(),
	}

	for i := 0; i < b.N; i++ {
		data, _ := Codec.Marshal(CodecVersion, account)
		fmt.Println(len(data))
	}
}
