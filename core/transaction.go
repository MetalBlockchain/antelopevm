package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/inconshreveable/log15"
)

type TransactionIdType = crypto.Sha256

type Extension struct {
	Type uint16   `serialize:"true" json:"type"`
	Data HexBytes `serialize:"true" json:"data"`
}

/**
*  The transaction header contains the fixed-sized data
*  associated with each transaction. It is separated from
*  the transaction body to facilitate partial parsing of
*  transactions without requiring dynamic memory allocation.
*
*  All transactions have an expiration time after which they
*  may no longer be included in the blockchain. Once a block
*  with a block_header::timestamp greater than expiration is
*  deemed irreversible, then a user can safely trust the transaction
*  will never be included.
 */
type TransactionHeader struct {
	Expiration     TimePointSec `serialize:"true" json:"expiration"`
	RefBlockNum    uint16       `serialize:"true" json:"ref_block_num"`
	RefBlockPrefix uint32       `serialize:"true" json:"ref_block_prefix"`

	MaxNetUsageWords Vuint32 `serialize:"true" json:"max_net_usage_words"`
	MaxCpuUsageMS    uint8   `serialize:"true" json:"max_cpu_usage_ms"`
	DelaySec         Vuint32 `serialize:"true" json:"delay_sec"` // number of secs to delay, making it cancellable for that duration
}

func (t TransactionHeader) Validate() {
	if t.MaxNetUsageWords >= math.MaxUint32/8 {
		panic("declared max_net_usage_words overflows when expanded to max net usage")
	}
}

type Transaction struct {
	TransactionHeader     `serialize:"true"`
	ContextFreeActions    []*Action    `serialize:"true" json:"context_free_actions"`
	Actions               []*Action    `serialize:"true" json:"actions"`
	TransactionExtensions []*Extension `serialize:"true" json:"transaction_extensions"`
}

func (t *Transaction) ID() TransactionIdType {
	b, err := rlp.EncodeToBytes(t)
	if err != nil {
		fmt.Println("Transaction ID() is error :", err.Error()) //TODO
	}
	enc := crypto.NewSha256()
	enc.Write(b)
	hashed := enc.Sum(nil)
	return TransactionIdType(*crypto.NewSha256Byte(hashed))
}

func (t *Transaction) SigDigest(chainID *ChainIdType, cfd []HexBytes) *DigestType {
	enc := crypto.NewSha256()
	chainIDByte, err := rlp.EncodeToBytes(chainID)

	if err != nil {
		fmt.Println(err)
	}

	thByte, err := rlp.EncodeToBytes(t)

	if err != nil {
		fmt.Println(err)
	}

	enc.Write(chainIDByte)
	enc.Write(thByte)

	if len(cfd) > 0 {
		enc.Write(crypto.Hash256(cfd).Bytes())
	} else {
		enc.Write(crypto.NewSha256Nil().Bytes())
	}

	hashed := enc.Sum(nil)

	return crypto.NewSha256Byte(hashed)
}

func (t *Transaction) GetSignatureKeys(signatures []ecc.Signature, chainID *ChainIdType, cfd []HexBytes, allowDuplicateKeys bool, useCache bool) ([]ecc.PublicKey, error) {
	digest := t.SigDigest(chainID, cfd)
	recovered := make(map[string]ecc.PublicKey)

	for _, sig := range signatures {
		recov, _ := sig.PublicKey(digest.Bytes())

		if _, found := recovered[recov.String()]; found {
			if !allowDuplicateKeys {
				return nil, errors.New("transaction includes more than one signature signed using the same key associated with public key")
			}
		}

		recovered[recov.String()] = recov
	}

	list := make([]ecc.PublicKey, len(recovered))

	for _, value := range recovered {
		list = append(list, value)
	}

	return list, nil
}

func (t *Transaction) TotalActions() uint32 {
	return uint32(len(t.ContextFreeActions) + len(t.Actions))
}

func (tx *Transaction) FirstAuthorizor() AccountName {
	for _, a := range tx.Actions {
		for _, auth := range a.Authorization {
			return auth.Actor
		}
	}
	return AccountName(0)
}

type SignedTransaction struct {
	Transaction     `serialize:"true"`
	Signatures      []ecc.Signature `serialize:"true" json:"signatures"`
	ContextFreeData []HexBytes      `serialize:"true" json:"context_free_data"`
}

func NewSignedTransaction(tx *Transaction, signature []ecc.Signature, contextFreeData []HexBytes) *SignedTransaction {
	return &SignedTransaction{
		Transaction:     *tx,
		Signatures:      signature,
		ContextFreeData: contextFreeData,
	}
}

func (s *SignedTransaction) Sign(key *ecc.PrivateKey, chainID *ChainIdType) ecc.Signature {
	signature, err := key.Sign(s.Transaction.SigDigest(chainID, s.ContextFreeData).Bytes())

	if err != nil {
		fmt.Println(err) // TODO: Handle this
	}

	s.Signatures = append(s.Signatures, signature)

	return signature
}

func (s *SignedTransaction) GetSignatureKeys(chainID *ChainIdType, allowDeplicateKeys bool, useCache bool) ([]ecc.PublicKey, error) {
	return s.Transaction.GetSignatureKeys(s.Signatures, chainID, s.ContextFreeData, allowDeplicateKeys, useCache)
}

type TransactionStatus uint8

const (
	TransactionStatusExecuted TransactionStatus = iota ///< succeed, no error handler executed
	TransactionStatusSoftFail                          ///< objectively failed (not executed), error handler executed
	TransactionStatusHardFail                          ///< objectively failed and error handler objectively failed thus no state change
	TransactionStatusDelayed                           ///< transaction delayed
	TransactionStatusExpired
	TransactionStatusUnknown = TransactionStatus(255)
)

func (s *TransactionStatus) UnmarshalJSON(data []byte) error {
	var decoded string
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	switch decoded {
	case "executed":
		*s = TransactionStatusExecuted
	case "soft_fail":
		*s = TransactionStatusSoftFail

	case "hard_fail":
		*s = TransactionStatusHardFail
	case "delayed":
		*s = TransactionStatusDelayed
	default:
		*s = TransactionStatusUnknown
	}
	return nil
}

func (s TransactionStatus) MarshalJSON() (data []byte, err error) {
	out := "unknown"
	switch s {
	case TransactionStatusExecuted:
		out = "executed"
	case TransactionStatusSoftFail:
		out = "soft_fail"
	case TransactionStatusHardFail:
		out = "hard_fail"
	case TransactionStatusDelayed:
		out = "delayed"
	}
	return json.Marshal(out)
}
func (s TransactionStatus) String() string {
	switch s {
	case TransactionStatusExecuted:
		return "executed"
	case TransactionStatusSoftFail:
		return "soft fail"
	case TransactionStatusHardFail:
		return "hard fail"
	case TransactionStatusDelayed:
		return "delayed"
	default:
		return "unknown"
	}
}

// PackedTransaction represents a fully packed transaction, with
// signatures, and all. They circulate like that on the P2P net, and
// that's how they are stored.
type PackedTransaction struct {
	Signatures            []ecc.Signature   `serialize:"true" json:"signatures"`
	Compression           CompressionType   `serialize:"true" json:"compression"` // in C++, it's an enum, not sure how it Binary-marshals..
	PackedContextFreeData HexBytes          `serialize:"true" json:"packed_context_free_data"`
	PackedTrx             HexBytes          `serialize:"true" json:"packed_trx"`
	UnpackedTrx           *Transaction      `json:"transaction"`
	Id                    TransactionIdType `json:"id"`
}

func NewPackedTransactionFromSignedTransaction(signedTrx SignedTransaction, compressionType CompressionType) (*PackedTransaction, error) {
	packedBytes, err := packTransaction(&signedTrx.Transaction)

	if err != nil {
		return nil, err
	}

	packedContextFreeData, err := packContextFreeData(&signedTrx.ContextFreeData)

	if err != nil {
		return nil, err
	}

	packedTrx := &PackedTransaction{
		Signatures:            signedTrx.Signatures,
		PackedTrx:             packedBytes,
		PackedContextFreeData: packedContextFreeData,
		Compression:           compressionType,
	}

	return packedTrx, nil
}

func (p *PackedTransaction) GetSignedTransaction() (*SignedTransaction, error) {
	if p.Compression == CompressionNone {
		unpackContextFreeData, err := unpackContextFreeData(&p.PackedContextFreeData)

		if err != nil {
			log15.Error("failed to get unpackContextFreeData", "error", err)
			return nil, err
		}

		unpackedTransaction, err := p.GetUnpackedTransaction()

		if err != nil {
			log15.Error("failed to get unpackedTransaction", "error", err)
			return nil, err
		}

		return NewSignedTransaction(unpackedTransaction, p.Signatures, unpackContextFreeData), nil
	}

	return nil, fmt.Errorf("unknown compression")
}

func (p *PackedTransaction) PackedDigest() crypto.Sha256 {
	prunable := crypto.NewSha256()
	result, _ := rlp.EncodeToBytes(p.Signatures)
	prunable.Write(result)
	result, _ = rlp.EncodeToBytes(p.PackedContextFreeData)
	prunable.Write(result)
	prunableResult := *crypto.NewSha256Byte(prunable.Sum(nil))

	enc := crypto.NewSha256()
	result, _ = rlp.EncodeToBytes(p.Compression)
	enc.Write(result)
	result, _ = rlp.EncodeToBytes(p.PackedTrx)
	enc.Write(result)
	result, _ = rlp.EncodeToBytes(prunableResult)
	enc.Write(result)

	return *crypto.NewSha256Byte(enc.Sum(nil))
}

func (p *PackedTransaction) GetUnpackedTransaction() (*Transaction, error) {
	if p.Compression == CompressionNone {
		return p.unpackTransaction()
	}

	return nil, fmt.Errorf("unknown compression")
}

func (p *PackedTransaction) unpackTransaction() (*Transaction, error) {
	transaction := &Transaction{}

	if err := rlp.DecodeBytes(p.PackedTrx, transaction); err != nil {
		return nil, err
	}

	p.UnpackedTrx = transaction
	p.Id = transaction.ID()

	return transaction, nil
}

func (p *PackedTransaction) MarshalJSON() ([]byte, error) {
	if p.UnpackedTrx == nil {
		p.unpackTransaction()
	}

	return json.Marshal(&struct {
		Signatures            []ecc.Signature   `json:"signatures"`
		Compression           CompressionType   `json:"compression"`
		PackedContextFreeData HexBytes          `json:"packed_context_free_data"`
		PackedTrx             HexBytes          `json:"packed_trx"`
		UnpackedTrx           *Transaction      `json:"transaction"`
		Id                    TransactionIdType `json:"id"`
	}{
		Signatures:            p.Signatures,
		Compression:           p.Compression,
		PackedContextFreeData: p.PackedContextFreeData,
		PackedTrx:             p.PackedTrx,
		UnpackedTrx:           p.UnpackedTrx,
		Id:                    p.Id,
	})
}

func unpackContextFreeData(data *HexBytes) ([]HexBytes, error) {
	out := make([]HexBytes, 0)

	if len(*data) == 0 {
		return out, nil
	}

	err := rlp.DecodeBytes([]byte(*data), &out)

	if err != nil {
		return nil, err
	}

	return out, nil
}

func packTransaction(t *Transaction) ([]byte, error) { //Bytes
	out, err := rlp.EncodeToBytes(t)

	if err != nil {
		return nil, err
	}

	return out, nil
}

func packContextFreeData(cfd *[]HexBytes) ([]byte, error) {
	if len(*cfd) == 0 {
		return []byte{}, nil
	}

	out, err := rlp.EncodeToBytes(cfd)

	if err != nil {
		return nil, err
	}

	return out, nil
}
