package protocol

import (
	"github.com/MetalBlockchain/antelopevm/core"
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
	FeatureDigest                 core.DigestType   `json:"feature_digest"`
	DescriptionDigest             core.DigestType   `json:"description_digest"`
	Dependencies                  []core.DigestType `json:"dependencies"`
	EarliestAllowedActivationTime core.TimePoint    `json:"earliest_allowed_activation_time"`
	PreactivationRequired         bool              `json:"preactivation_required"`
	Enabled                       bool              `json:"enabled"`
}

type ProtocolFeatureSubjectiveRestrictions struct {
	EarliestAllowedActivationTime core.TimePoint `json:"earliest_allowed_activation_time"`
	PreactivationRequired         bool           `json:"preactivation_required"`
	Enabled                       bool           `json:"enabled"`
}

type BuiltinProtocolFeatureSpec struct {
	CodeName               string
	DescriptionDigest      core.DigestType
	BuiltinDependencies    []core.DigestType
	SubjectiveRestrictions ProtocolFeatureSubjectiveRestrictions
	Type                   BuiltinProtocolFeatureType
}

type BuiltinProtocolFeature struct {
	ProtocolFeatureType    string
	DescriptionDigest      core.DigestType
	Dependencies           []core.DigestType
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
