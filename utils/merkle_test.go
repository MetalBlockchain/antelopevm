package utils_test

import (
	"encoding/hex"
	"testing"

	"github.com/MetalBlockchain/antelopevm/chain/transaction"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/utils"
	"github.com/stretchr/testify/assert"
)

func TestTransactionMerkleRootEven(t *testing.T) {
	packed1, err := hex.DecodeString("f6bc0e654154d1cf301b000000000190113253419a7bd5000000572d3ccdcd01a0649a2656ed4dac000000c067175dd621a0649a2656ed4dac5062cd0355efe97b444201000000000002555058000000000000")
	assert.NoError(t, err)
	packed2, err := hex.DecodeString("f6bc0e654154d1cf301b0000000001a0649a2656ed4dac000000000000089901a0649a2656ed4dac000000c067175dd61ba0d8da5bb852ba79c027090000000000025550580000000002626e00")
	assert.NoError(t, err)
	signature1, err := ecc.NewSignature("SIG_K1_KAfFoHJpwZsTSKoyEuHE5c4i8J5wK95r6LxqMzLddk5D9obsiaaTRYrcre9gKDKJE6y4Ckd7uLM62FsbboJ34iioM8cEsV")
	assert.NoError(t, err)
	signature2, err := ecc.NewSignature("SIG_K1_JyThXTKQji3xAKEgK7v99vSQypFN9mJPuQcxy2A3hkCEw1ABpopNYV43Tqn8rAZRAYXZVP3WiXnEtv2rEzEzCcrHemxpgB")
	assert.NoError(t, err)
	receipt1 := transaction.TransactionReceipt{
		TransactionReceiptHeader: transaction.TransactionReceiptHeader{
			Status:        transaction.TransactionStatusExecuted,
			CpuUsageUs:    268,
			NetUsageWords: 16,
		},
		Transaction: transaction.PackedTransaction{
			PackedTrx:             packed1,
			Compression:           transaction.CompressionNone,
			Signatures:            []ecc.Signature{signature1},
			PackedContextFreeData: make(types.HexBytes, 0),
		},
	}
	receipt2 := transaction.TransactionReceipt{
		TransactionReceiptHeader: transaction.TransactionReceiptHeader{
			Status:        transaction.TransactionStatusExecuted,
			CpuUsageUs:    213,
			NetUsageWords: 15,
		},
		Transaction: transaction.PackedTransaction{
			PackedTrx:             packed2,
			Compression:           transaction.CompressionNone,
			Signatures:            []ecc.Signature{signature2},
			PackedContextFreeData: make(types.HexBytes, 0),
		},
	}

	err = receipt1.Transaction.UnpackTransaction()
	assert.NoError(t, err)
	err = receipt2.Transaction.UnpackTransaction()
	assert.NoError(t, err)
	digest1, err := receipt1.Digest()
	assert.NoError(t, err)
	digest2, err := receipt2.Digest()
	assert.NoError(t, err)
	hashes := []crypto.Sha256{*digest1, *digest2}
	merkle := utils.Merkle(hashes)
	assert.Equal(t, merkle.String(), "98db7a674310484bb1444e4981a0a6c79cdaacb94ea3909431c7ad8a4134e421")
}

func TestTransactionMerkleRootOdd(t *testing.T) {
	packed1, err := hex.DecodeString("09bb0e65a05477af3d970000000002a0649a2656ed4dac0000d0155dbabca901a0649a2656ed4dac0000d0155dbabca908a0845467bd413991a0649a2656ed4dac000000000000c29801a0845467bd41399100000000a46962d59903a0845467bd4139913207dd26d1b54a000097d026b3b74a0000c0db2610b64a0000c8d1267ab74a00005ddc2669b54a000083dc2656b54a0000ebd1262eb74a000023d3268ab74a0000a4dc269ab54a000066d22610b74a000091dc2648b54a000060dd26f5b54a00001ddd26d2b54a0000bddb263fb64a0000c2d12671b74a000015d3268ab74a0000bfd22688b74a00007cdc2669b54a00005ddc26a2b54a000095dc2648b54a000092dc2604b54a000099dc2648b54a00002bdd26b0b54a0000cadb2610b64a00008ddc2678b54a0000d4dd26b3b54a00005fdc26bfb44a00001cd3268ab74a0000c6db2610b64a00007edc2699b54a000087e4268ab24a000071022709924a000054da2698b94a000012dd26afb54a0000b7db260fb64a00003ddc2655b54a00005ad22643b74a0000b1db26edb54a000085dc2604b54a0000a0dc2656b54a00009edc2626b54a0000e9dd26b3b54a00007adc2612b54a000053d026c8b74a000070e026d2b44a000087dc2612b54a00005ce026c3b44a00004ddc26b8b34a00008edc269ab54a000085dc2678b54a000000")
	assert.NoError(t, err)
	signature1, err := ecc.NewSignature("SIG_K1_K6Xig9bAh4xg9LJVnGutDEAo8awxb35JAvjnark9urS5ZsZqtiY9FgsMrGa479JS9E9biicHbeKSAtXaFup5VWPYyQyjyU")
	assert.NoError(t, err)
	signature2, err := ecc.NewSignature("SIG_K1_KaMx12aj799w33eq8pyQoi5iwdUuGV8BJ96xYtmUTnn2PsuuXZ9qdbfpiuMS3K8em6pRCD8AHSLHU4aJi22WPV52ErmKNk")
	assert.NoError(t, err)
	receipt1 := transaction.TransactionReceipt{
		TransactionReceiptHeader: transaction.TransactionReceiptHeader{
			Status:        transaction.TransactionStatusExecuted,
			CpuUsageUs:    925,
			NetUsageWords: 70,
		},
		Transaction: transaction.PackedTransaction{
			PackedTrx:             packed1,
			Compression:           transaction.CompressionNone,
			Signatures:            []ecc.Signature{signature1, signature2},
			PackedContextFreeData: make(types.HexBytes, 0),
		},
	}
	err = receipt1.Transaction.UnpackTransaction()
	assert.NoError(t, err)
	digest1, err := receipt1.Digest()
	assert.NoError(t, err)
	hashes := []crypto.Sha256{*digest1}
	merkle := utils.Merkle(hashes)
	assert.Equal(t, merkle.String(), "9b11ff83d03fe4fe8efff148a930f6b5dff788995ae8f89259bd423af2e3ce64")
}
