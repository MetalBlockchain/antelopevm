package core_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/MetalBlockchain/antelopevm/core/protocol"
	"github.com/MetalBlockchain/antelopevm/crypto"
)

func TestX(t *testing.T) {
	a := protocol.BuiltinProtocolFeature{
		DescriptionDigest: *crypto.NewSha256String("64fe7df32e9b86be2b296b3f81dfd527f84e82b98e363bc97e40bc7a83733310"),
		Type:              protocol.Builtin,
		CodeName:          protocol.PreactivateFeature,
		Dependencies:      make([]crypto.Sha256, 0),
	}

	fmt.Printf("hex.EncodeToString(a.Digest()): %v\n", hex.EncodeToString(a.Digest()))
}
