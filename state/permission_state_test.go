package state

import (
	"testing"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/authority"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
)

func TestPermissionState(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions(t.TempDir()))
	assert.NoError(t, err)
	state := NewState(nil, db)
	session := state.CreateSession(true)
	permission := &core.Permission{
		Parent:      0,
		Owner:       name.StringToName("glenn"),
		Name:        name.StringToName("owner"),
		LastUpdated: core.Now(),
		LastUsed:    core.Now(),
		Auth: authority.Authority{
			Threshold: 1,
			Keys: []authority.KeyWeight{
				{
					Key:    *ecc.NewPublicKeyNil(),
					Weight: 1,
				},
			},
		},
	}
	session.CreatePermission(permission)
	permission2 := &core.Permission{
		Parent:      permission.ID,
		LastUsed:    core.Now(),
		Owner:       name.StringToName("glenn"),
		Name:        name.StringToName("active"),
		LastUpdated: core.Now(),
		Auth: authority.Authority{
			Threshold: 1,
			Keys: []authority.KeyWeight{
				{
					Key:    *ecc.NewPublicKeyNil(),
					Weight: 1,
				},
			},
		},
	}
	session.CreatePermission(permission2)
	session.Commit()
	session = state.CreateSession(false)
	permission, err = session.FindPermissionByOwner(name.StringToName("glenn"), name.StringToName("active"))
	assert.NoError(t, err)
	assert.NotNil(t, permission)
	assert.Equal(t, permission.Name, name.StringToName("active"))
}
