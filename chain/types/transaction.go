package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
)

//var recoveryCache = make(map[string]CachedPubKey)

type CachedPubKey struct {
	TrxID  TransactionIdType `json:"trx_id"`
	PubKey ecc.PublicKey     `json:"pub_key"`
	Sig    ecc.Signature     `json:"sig"`
}

type TransactionIdType = crypto.Sha256

type Extension struct {
	Type uint16   `json:"type"`
	Data HexBytes `json:"data"`
}

type TransactionHeader struct {
	Expiration     TimePointSec `json:"expiration"`
	RefBlockNum    uint16       `json:"ref_block_num"`
	RefBlockPrefix uint32       `json:"ref_block_prefix"`

	MaxNetUsageWords Vuint32 `json:"max_net_usage_words" eos:"vuint32"`
	MaxCpuUsageMS    uint8   `json:"max_cpu_usage_ms"`
	DelaySec         Vuint32 `json:"delay_sec" eos:"vuint32"` // number of secs to delay, making it cancellable for that duration
}

func (t TransactionHeader) IsEmpty() bool {
	return t.Expiration == 0 && t.RefBlockNum == 0 && t.RefBlockPrefix == 0 && t.MaxNetUsageWords == 0 && t.MaxCpuUsageMS == 0 && t.DelaySec == 0
}
func (t TransactionHeader) GetRefBlocknum(headBlocknum uint32) uint32 {
	return headBlocknum/0xffff*0xffff + headBlocknum%0xffff
}

func (t TransactionHeader) VerifyReferenceBlock(referenceBlock *BlockIdType) bool {
	return t.RefBlockNum == uint16(EndianReverseU32(uint32(referenceBlock.Hash[0]))) && t.RefBlockPrefix == uint32(referenceBlock.Hash[1])
}

func (t TransactionHeader) Validate() {
	if t.MaxNetUsageWords >= math.MaxUint32/8 {
		panic("declared max_net_usage_words overflows when expanded to max net usage")
	}
}

func (t *TransactionHeader) SetReferenceBlock(referenceBlock *BlockIdType) {
	first := EndianReverseU32(uint32(referenceBlock.Hash[0]))
	t.RefBlockNum = uint16(first)
	t.RefBlockPrefix = uint32(referenceBlock.Hash[1])
}

type Transaction struct {
	TransactionHeader
	ContextFreeActionLength Vuint32
	ContextFreeActions      []*Action `json:"context_free_actions"`
	ActionLength            Vuint32
	Actions                 []*Action `json:"actions"`
	ExtensionLength         Vuint32
	TransactionExtensions   []*Extension `json:"transaction_extensions"`
}

func (t Transaction) IsEmtpy() bool {
	return len(t.ContextFreeActions) == 0 && len(t.Actions) == 0 && len(t.TransactionExtensions) == 0 && t.TransactionHeader.IsEmpty()
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

type SignedTransaction struct {
	Transaction

	Signatures      []ecc.Signature `json:"signatures"`
	ContextFreeData []HexBytes      `json:"context_free_data"`
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

type BlockStatus uint8

const (
	Irreversible BlockStatus = iota ///< this block has already been applied before by this node and is considered irreversible
	Validated                       ///< this is a complete block signed by a valid producer and has been previously applied by this node and therefore validated but it is not yet irreversible
	Complete                        ///< this is a complete block signed by a valid producer but is not yet irreversible nor has it yet been applied by this node
	Incomplete                      ///< this is an incomplete block (either being produced by a producer or speculatively produced by a node)
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

type TransactionReceiptHeader struct {
	Status        TransactionStatus `json:"status"`
	CpuUsageUs    uint32            `json:"cpu_usage_us"`
	NetUsageWords Vuint32           `json:"net_usage_words" eos:"vuint32"`
}

func (t TransactionReceiptHeader) IsEmpty() bool {
	return t.Status == 0 && t.CpuUsageUs == 0 && t.NetUsageWords == 0
}

type TransactionReceipt struct {
	TransactionReceiptHeader
	Trx TransactionWithID `json:"trx" eos:"trxID"`
}

type TransactionWithID struct {
	PackedTransaction *PackedTransaction `json:"packed_transaction" eos:"tag0"`
	TransactionID     TransactionIdType  `json:"transaction_id" eos:"tag1"`
}

// PackedTransaction represents a fully packed transaction, with
// signatures, and all. They circulate like that on the P2P net, and
// that's how they are stored.
type PackedTransaction struct {
	Signatures            []ecc.Signature `json:"signatures"`
	Compression           CompressionType `json:"compression"` // in C++, it's an enum, not sure how it Binary-marshals..
	PackedContextFreeData HexBytes        `json:"packed_context_free_data"`
	PackedTrx             HexBytes        `json:"packed_trx"`
	UnpackedTrx           *Transaction    `json:"transaction" eos:"-"`
}

func (p *PackedTransaction) GetSignedTransaction() (*SignedTransaction, error) {
	if p.Compression == CompressionNone {
		unpackContextFreeData, err := unpackContextFreeData(&p.PackedContextFreeData)

		if err != nil {
			return nil, err
		}

		unpackedTransaction, err := p.GetUnpackedTransaction()

		if err != nil {
			return nil, err
		}

		return NewSignedTransaction(unpackedTransaction, p.Signatures, unpackContextFreeData), nil
	}

	return nil, fmt.Errorf("unknown compression")
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

	return transaction, nil
}

type CompressionType uint8

const (
	CompressionNone = CompressionType(iota)
	CompressionZlib
)

func (c CompressionType) String() string {
	switch c {
	case CompressionNone:
		return "none"
	case CompressionZlib:
		return "zlib"
	default:
		return ""
	}
}

func (c CompressionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *CompressionType) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	switch s {
	case "zlib":
		*c = CompressionZlib
	default:
		*c = CompressionNone
	}
	return nil
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
