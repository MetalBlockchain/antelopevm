package transaction

import (
	"github.com/MetalBlockchain/antelopevm/chain/time"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
)

type TransactionType uint32

const (
	Input     TransactionType = 0
	Implicit  TransactionType = 1
	Scheduled TransactionType = 2
	DryRun    TransactionType = 3
)

type TransactionMetaData struct {
	packedTransaction   *PackedTransaction
	signatureCpuUsage   time.Microseconds
	recoveredPublicKeys ecc.PublicKeySet
	transactionType     TransactionType
	Accepted            bool
	BilledCpuTimeUs     uint32
}

func RecoverKeys(trx *PackedTransaction, chainId types.ChainIdType, timeLimit time.Microseconds, trxType TransactionType, maxVariableSignatureSize uint32) (*TransactionMetaData, error) {
	deadline := time.MaxTimePoint()

	if timeLimit != time.MaxMicroseconds() {
		deadline = time.Now() + time.TimePoint(timeLimit)
	}

	// TODO: check_variable_sig_size
	signedTrx, err := trx.GetSignedTransaction()

	if err != nil {
		return nil, err
	}

	recoveredPublicKeys, cpuUsage, err := signedTrx.GetSignatureKeys(&chainId, deadline, false)

	if err != nil {
		return nil, err
	}

	return &TransactionMetaData{
		packedTransaction:   trx,
		signatureCpuUsage:   time.Microseconds(cpuUsage),
		recoveredPublicKeys: recoveredPublicKeys,
		transactionType:     trxType,
	}, nil
}

func (m *TransactionMetaData) Implicit() bool {
	return m.transactionType == Implicit
}

func (m *TransactionMetaData) Scheduled() bool {
	return m.transactionType == Scheduled
}

func (m *TransactionMetaData) IsDryRun() bool {
	return m.transactionType == DryRun
}

func (m *TransactionMetaData) Id() *TransactionIdType {
	id, _ := m.packedTransaction.ID()
	return id
}

func (m *TransactionMetaData) RecoveredKeys() ecc.PublicKeySet {
	return m.recoveredPublicKeys
}

func (m *TransactionMetaData) PackedTrx() *PackedTransaction {
	return m.packedTransaction
}

func (m *TransactionMetaData) GetTrxType() TransactionType {
	return m.transactionType
}
