package account_test

import (
	"testing"

	"github.com/MetalBlockchain/antelopevm/chain/account"
	"github.com/stretchr/testify/assert"
)

func TestSetPrivileged(t *testing.T) {
	account := account.AccountMetaDataObject{}
	account.SetPrivileged(true)
	assert.Equal(t, true, account.IsPrivileged())
	account.SetPrivileged(false)
	assert.Equal(t, false, account.IsPrivileged())
}
