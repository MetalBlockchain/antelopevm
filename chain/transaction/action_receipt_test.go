package transaction_test

import (
	"encoding/hex"
	"testing"

	"github.com/MetalBlockchain/antelopevm/chain/authority"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/transaction"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/stretchr/testify/assert"
)

func TestActionReceiptDigest(t *testing.T) {
	authSequence := authority.NewAuthSequenceSet()
	authSequence.Set(name.StringToName("glenn"), 2)
	receipt := transaction.ActionReceipt{
		Receiver:       name.StringToName("glenn"),
		AuthSequence:   authSequence,
		GlobalSequence: 1,
		RecvSequence:   3,
	}
	digest := receipt.Digest()
	assert.Equal(t, "3f285f90ea0e7925af98249d4e875bbfc54e936c06d351229bb7dc4e176e4562", digest.String())
	encoded, err := rlp.EncodeToBytes(receipt)
	assert.NoError(t, err)
	assert.Equal(t, "000000008039556400000000000000000000000000000000000000000000000000000000000000000100000000000000030000000000000001000000008039556402000000000000000000", hex.EncodeToString(encoded))
}
