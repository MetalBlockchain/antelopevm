package chain

import (
	"testing"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/stretchr/testify/assert"
)

func TestAuthPermissionLevelKeyWeight(t *testing.T) {
	privateKey, _ := ecc.NewRandomPrivateKey()
	keys := make([]ecc.PublicKey, 0)
	keys = append(keys, privateKey.PublicKey())
	permissions := make([]core.PermissionLevel, 0)
	permissionToAuthority := func(*core.PermissionLevel) (*core.Authority, error) {
		return &core.Authority{
			Threshold: 1,
			Keys: []core.KeyWeight{{
				Key:    privateKey.PublicKey(),
				Weight: 1,
			}},
		}, nil
	}
	checker := NewAuthorityChecker(permissionToAuthority, keys, permissions, 16)

	ok := checker.SatisfiedPermissionLevel(core.PermissionLevel{core.StringToName("glenn"), 1}, nil)
	assert.Equal(t, ok, true, "should be ok")
}

func TestAuthPermissionLevelAccountWeight(t *testing.T) {
	privateKey, _ := ecc.NewRandomPrivateKey()
	keys := make([]ecc.PublicKey, 0)
	keys = append(keys, privateKey.PublicKey())
	permissions := make([]core.PermissionLevel, 0)
	permissionToAuthority := func(a *core.PermissionLevel) (*core.Authority, error) {
		if a.Actor.String() == "joe" {
			return &core.Authority{
				Threshold: 1,
				Keys: []core.KeyWeight{{
					Key:    privateKey.PublicKey(),
					Weight: 1,
				}},
			}, nil
		}

		return &core.Authority{
			Threshold: 1,
			Accounts: []core.PermissionLevelWeight{{
				Permission: core.PermissionLevel{core.StringToName("joe"), 1},
				Weight:     1,
			}},
		}, nil
	}
	checker := NewAuthorityChecker(permissionToAuthority, keys, permissions, 16)

	ok := checker.SatisfiedPermissionLevel(core.PermissionLevel{core.StringToName("glenn"), 1}, nil)
	assert.Equal(t, ok, true, "should be ok")
}
