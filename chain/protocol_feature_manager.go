package chain

import (
	"github.com/MetalBlockchain/antelopevm/chain/protocol"
	"github.com/MetalBlockchain/antelopevm/chain/time"
	"github.com/MetalBlockchain/antelopevm/crypto"
)

var BuiltinProtocolFeatureCodenames = map[protocol.BuiltinProtocolFeatureType]protocol.BuiltinProtocolFeatureSpec{
	protocol.PreactivateFeature: {
		CodeName:          "PREACTIVATE_FEATURE",
		DescriptionDigest: *crypto.NewSha256String("64fe7df32e9b86be2b296b3f81dfd527f84e82b98e363bc97e40bc7a83733310"),
		SubjectiveRestrictions: protocol.ProtocolFeatureSubjectiveRestrictions{
			EarliestAllowedActivationTime: time.MinTimePoint(),
			PreactivationRequired:         false,
			Enabled:                       true,
		},
	},
	protocol.OnlyLinkToExistingPermission: {
		CodeName:          "ONLY_LINK_TO_EXISTING_PERMISSION",
		DescriptionDigest: *crypto.NewSha256String("f3c3d91c4603cde2397268bfed4e662465293aab10cd9416db0d442b8cec2949"),
		SubjectiveRestrictions: protocol.ProtocolFeatureSubjectiveRestrictions{
			EarliestAllowedActivationTime: time.MinTimePoint(),
			PreactivationRequired:         true,
			Enabled:                       true,
		},
	},
}
