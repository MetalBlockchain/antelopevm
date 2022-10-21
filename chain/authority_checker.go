package chain

import (
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/state"
)

type AuthorityChecker struct {
	ProvidedKeys []ecc.PublicKey
}

func NewAuthorityChecker(keys []ecc.PublicKey) *AuthorityChecker {
	return &AuthorityChecker{
		ProvidedKeys: keys,
	}
}

// TODO: Handle wait and account based checks
func (a *AuthorityChecker) Check(permission state.Permission) bool {
	totalWeight := 0

	for _, keyWeight := range permission.Auth.Keys {
		totalWeight += a.CheckKeyWeight(keyWeight)
	}

	return totalWeight >= int(permission.Auth.Threshold)
}

func (a *AuthorityChecker) CheckKeyWeight(permission types.KeyWeight) int {
	for _, key := range a.ProvidedKeys {
		if key.Compare(permission.GetPublicKey()) {
			return 1
		}
	}

	return 0
}
