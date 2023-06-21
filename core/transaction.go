package core

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"

	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/inconshreveable/log15"
)

//go:generate msgp
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

func (t *Transaction) GetSignatureKeys(signatures []ecc.Signature, chainId *ChainIdType, deadline TimePoint, cfd []HexBytes, allowDuplicateKeys bool) (ecc.PublicKeySet, TimePoint, error) {
	start := Now()
	set := ecc.NewPublicKeySet(len(signatures))
	digest := t.SigDigest(chainId, cfd)

	for i := 0; i < len(signatures); i++ {
		now := Now()

		if now >= deadline {
			return nil, Now() - start, fmt.Errorf("transaction signature verification executed for too long %sus", now-start)
		}

		recov, err := signatures[i].PublicKey(digest.Bytes())

		if err != nil {
			return nil, Now() - start, err
		}

		if set.Contains(recov) {
			if !allowDuplicateKeys {
				return nil, Now() - start, fmt.Errorf("transaction includes more than one signature signed using the same key associated with public key: %s", recov.String())
			}
		} else {
			set.Insert(recov)
		}
	}

	return set, Now() - start, nil
}

func (t *Transaction) TotalActions() uint32 {
	return uint32(len(t.ContextFreeActions) + len(t.Actions))
}

func (tx *Transaction) FirstAuthorizor() name.AccountName {
	for _, a := range tx.Actions {
		for _, auth := range a.Authorization {
			return auth.Actor
		}
	}
	return name.AccountName(0)
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

func (s *SignedTransaction) GetSignatureKeys(chainId *ChainIdType, deadline TimePoint, allowDeplicateKeys bool) (ecc.PublicKeySet, TimePoint, error) {
	return s.Transaction.GetSignatureKeys(s.Signatures, chainId, deadline, s.ContextFreeData, allowDeplicateKeys)
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
	err := p.UnpackTransaction()

	if err != nil {
		return nil, err
	}

	contextFreeData, err := p.GetContextFreeData()

	if err != nil {
		return nil, err
	}

	return NewSignedTransaction(p.UnpackedTrx, p.Signatures, contextFreeData), nil
}

func (p *PackedTransaction) PackedDigest() crypto.Sha256 {
	log15.Info("PackedDigest")
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

func (p *PackedTransaction) GetTransaction() (*Transaction, error) {
	if p.UnpackedTrx != nil {
		return p.UnpackedTrx, nil
	}

	if err := p.UnpackTransaction(); err == nil {
		return p.UnpackedTrx, nil
	} else {
		return nil, err
	}
}

func (p *PackedTransaction) UnpackTransaction() error {
	if p.UnpackedTrx != nil {
		return nil
	}

	if p.Compression == CompressionNone {
		unpacked, err := unpackTransaction(p.PackedTrx)

		if err != nil {
			return err
		}

		p.UnpackedTrx = unpacked
		p.Id = unpacked.ID()

		return nil
	} else if p.Compression == CompressionZlib {
		unpacked, err := zlibDecompressTransaction(&p.PackedTrx)

		if err != nil {
			return err
		}

		p.UnpackedTrx = unpacked
		p.Id = unpacked.ID()

		return nil
	}

	return fmt.Errorf("unknown compression")
}

func (p *PackedTransaction) GetContextFreeData() ([]HexBytes, error) {
	log15.Info("GetContextFreeData")
	if p.Compression == CompressionNone {
		return unpackContextFreeData(&p.PackedContextFreeData)
	} else if p.Compression == CompressionZlib {
		return zlibDecompressContextFreeData(&p.PackedContextFreeData)
	}

	return nil, fmt.Errorf("unknown compression")
}

func (p *PackedTransaction) GetUnprunableSize() uint32 {
	size := uint32(16) // FixedNetOverheadOfPackedTrx
	size += uint32(len(p.PackedTrx))
	return size
}

func (p *PackedTransaction) GetPrunableSize() uint32 {
	size, _ := rlp.EncodeSize(p.Signatures)
	size += len(p.PackedContextFreeData)
	return uint32(size)
}

func (p *PackedTransaction) MarshalJSON() ([]byte, error) {
	log15.Info("MarshalJSON")
	err := p.UnpackTransaction()

	if err != nil {
		return nil, err
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
		Id:                    p.UnpackedTrx.ID(),
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

func unpackTransaction(data HexBytes) (*Transaction, error) {
	tx := &Transaction{}

	if err := rlp.DecodeBytes(data, tx); err != nil {
		return nil, err
	}

	return tx, nil
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

func zlibDecompressTransaction(data *HexBytes) (*Transaction, error) {
	packedTrx, err := zlibDecompress(data)

	if err != nil {
		return nil, err
	}

	return unpackTransaction(packedTrx)
}

func zlibDecompressContextFreeData(data *HexBytes) ([]HexBytes, error) {
	if len(*data) == 0 {
		return []HexBytes{}, nil
	}

	packedData, err := zlibDecompress(data)

	if err != nil {
		return nil, err
	}

	return unpackContextFreeData(&packedData)
}

func zlibDecompress(data *HexBytes) (HexBytes, error) {
	in := bytes.NewReader(*data)
	reader, err := zlib.NewReader(in)

	if err != nil {
		return nil, err
	}

	defer reader.Close()
	result, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	return result, nil
}

type TransactionMetaData struct {
}
