package transaction

import (
	"github.com/MetalBlockchain/antelopevm/core"
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
	packedTransaction   *core.PackedTransaction
	signatureCpuUsage   core.Microseconds
	recoveredPublicKeys ecc.PublicKeySet
	transactionType     TransactionType
	Accepted            bool
	BilledCpuTimeUs     uint32
}

func RecoverKeys(trx *core.PackedTransaction, chainId core.ChainIdType, timeLimit core.Microseconds, trxType TransactionType, maxVariableSignatureSize uint32) (*TransactionMetaData, error) {
	deadline := core.MaxTimePoint()

	if timeLimit != core.MaxMicroseconds() {
		deadline = core.Now() + core.TimePoint(timeLimit)
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
		signatureCpuUsage:   core.Microseconds(cpuUsage),
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

func (m *TransactionMetaData) Id() core.TransactionIdType {
	return m.packedTransaction.Id
}

func (m *TransactionMetaData) RecoveredKeys() ecc.PublicKeySet {
	return m.recoveredPublicKeys
}

func (m *TransactionMetaData) PackedTrx() *core.PackedTransaction {
	return m.packedTransaction
}

func (m *TransactionMetaData) GetTrxType() TransactionType {
	return m.transactionType
}
