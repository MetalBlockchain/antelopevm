package transaction

import (
	"github.com/MetalBlockchain/antelopevm/chain/fc"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
)

type TransactionReceiptHeader struct {
	Status        TransactionStatus `serialize:"true" json:"status"`
	CpuUsageUs    uint32            `serialize:"true" json:"cpu_usage_us"`
	NetUsageWords fc.UnsignedInt    `serialize:"true" json:"net_usage_words"`
}

type TransactionReceipt struct {
	TransactionReceiptHeader `serialize:"true"`
	Transaction              PackedTransaction `serialize:"true" json:"trx" eos:"-"`
}

func (t *TransactionReceipt) Digest() (*crypto.Sha256, error) {
	enc := crypto.NewSha256()

	if data, err := rlp.EncodeMultipleToBytes(t.Status, t.CpuUsageUs, t.NetUsageWords); err != nil {
		return nil, err
	} else {
		enc.Write(data)
	}

	packedDigest, err := t.Transaction.PackedDigest()

	if err != nil {
		return nil, err
	}

	if packedTrx, err := rlp.EncodeToBytes(packedDigest); err != nil {
		return nil, err
	} else {
		enc.Write(packedTrx)
	}

	return crypto.NewSha256Byte(enc.Sum(nil)), nil
}
