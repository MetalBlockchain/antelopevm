package protocol

import (
	"github.com/MetalBlockchain/antelopevm/chain/time"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
)

type ProtocolFeatureType uint32
type BuiltinProtocolFeatureType uint32

const (
	Builtin ProtocolFeatureType = 0
)

const (
	PreactivateFeature           BuiltinProtocolFeatureType = 0
	OnlyLinkToExistingPermission BuiltinProtocolFeatureType = 1
	ReplaceDeferred              BuiltinProtocolFeatureType = 2
	NoDuplicateDeferredId        BuiltinProtocolFeatureType = 3
)

type ProtocolFeature struct {
	FeatureDigest                 types.DigestType   `json:"feature_digest"`
	DescriptionDigest             types.DigestType   `json:"description_digest"`
	Dependencies                  []types.DigestType `json:"dependencies"`
	EarliestAllowedActivationTime time.TimePoint     `json:"earliest_allowed_activation_time"`
	PreactivationRequired         bool               `json:"preactivation_required"`
	Enabled                       bool               `json:"enabled"`
}

type ProtocolFeatureSubjectiveRestrictions struct {
	EarliestAllowedActivationTime time.TimePoint `json:"earliest_allowed_activation_time"`
	PreactivationRequired         bool           `json:"preactivation_required"`
	Enabled                       bool           `json:"enabled"`
}

type BuiltinProtocolFeatureSpec struct {
	CodeName               string
	DescriptionDigest      types.DigestType
	BuiltinDependencies    []types.DigestType
	SubjectiveRestrictions ProtocolFeatureSubjectiveRestrictions
	Type                   BuiltinProtocolFeatureType
}

type BuiltinProtocolFeature struct {
	ProtocolFeatureType    string
	DescriptionDigest      types.DigestType
	Dependencies           []types.DigestType
	SubjectiveRestrictions ProtocolFeatureSubjectiveRestrictions
	Type                   ProtocolFeatureType
	CodeName               BuiltinProtocolFeatureType
}

func (f *BuiltinProtocolFeature) Digest() []byte {
	typeEnc, _ := rlp.EncodeToBytes(f.Type)
	descriptionDigestEnc, _ := rlp.EncodeToBytes(f.DescriptionDigest)
	dependenciesEnc, _ := rlp.EncodeToBytes(f.Dependencies)
	codeNameEnc, _ := rlp.EncodeToBytes(f.CodeName)

	enc := crypto.NewSha256()
	enc.Write(typeEnc)
	enc.Write(descriptionDigestEnc)
	enc.Write(dependenciesEnc)
	enc.Write(codeNameEnc)

	return enc.Sum(nil)
}
