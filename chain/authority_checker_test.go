package chain

import (
	"testing"

	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/stretchr/testify/assert"
)

func TestAuthPermissionLevelKeyWeight(t *testing.T) {
	privateKey, _ := ecc.NewRandomPrivateKey()
	keys := make([]ecc.PublicKey, 0)
	keys = append(keys, privateKey.PublicKey())
	permissions := make([]types.PermissionLevel, 0)
	permissionToAuthority := func(*types.PermissionLevel) (*types.Authority, error) {
		return &types.Authority{
			Threshold: 1,
			Keys: []types.KeyWeight{{
				Key:    privateKey.PublicKey().String(),
				Weight: 1,
			}},
		}, nil
	}
	checker := NewAuthorityChecker(permissionToAuthority, keys, permissions, 16)

	ok := checker.SatisfiedPermissionLevel(types.PermissionLevel{types.N("glenn"), 1}, nil)
	assert.Equal(t, ok, true, "should be ok")
}

func TestAuthPermissionLevelAccountWeight(t *testing.T) {
	privateKey, _ := ecc.NewRandomPrivateKey()
	keys := make([]ecc.PublicKey, 0)
	keys = append(keys, privateKey.PublicKey())
	permissions := make([]types.PermissionLevel, 0)
	permissionToAuthority := func(a *types.PermissionLevel) (*types.Authority, error) {
		if a.Actor.String() == "joe" {
			return &types.Authority{
				Threshold: 1,
				Keys: []types.KeyWeight{{
					Key:    privateKey.PublicKey().String(),
					Weight: 1,
				}},
			}, nil
		}

		return &types.Authority{
			Threshold: 1,
			Accounts: []types.PermissionLevelWeight{{
				Permission: types.PermissionLevel{types.N("joe"), 1},
				Weight:     1,
			}},
		}, nil
	}
	checker := NewAuthorityChecker(permissionToAuthority, keys, permissions, 16)

	ok := checker.SatisfiedPermissionLevel(types.PermissionLevel{types.N("glenn"), 1}, nil)
	assert.Equal(t, ok, true, "should be ok")
}
