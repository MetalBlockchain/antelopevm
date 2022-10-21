package state

import (
	"testing"

	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/metalgo/database/memdb"
	"github.com/stretchr/testify/assert"
)

func TestReadAndWrite(t *testing.T) {
	baseDb := memdb.New()
	accountState := NewAccountState(baseDb)
	account := &Account{
		Name:       types.N("glenn"),
		Privileged: false,
	}
	err := accountState.PutAccount(account)

	assert.Nil(t, err, "error when writing account")
	assert.Equal(t, types.IdType(0), account.ID, "account should have id 0")

	// Now let's read the account
	account, err = accountState.GetAccountByName(account.Name)

	assert.Nil(t, err, "error when reading account")
	assert.Equal(t, types.N("glenn"), account.Name, "name should equal glenn")
}

func TestUpdate(t *testing.T) {
	baseDb := memdb.New()
	accountState := NewAccountState(baseDb)
	account := &Account{
		Name:       types.N("glenn"),
		Privileged: false,
	}
	err := accountState.PutAccount(account)

	assert.Nil(t, err, "error when writing account")
	assert.Equal(t, types.IdType(0), account.ID, "account should have id 0")

	// Now let's read the account
	account.Privileged = true
	err = accountState.UpdateAccount(account)
	assert.Nil(t, err, "error when updating account")
	account, err = accountState.GetAccountByName(account.Name)

	assert.Nil(t, err, "error when reading account")
	assert.Equal(t, true, account.Privileged, "privileged should equal true")
}
