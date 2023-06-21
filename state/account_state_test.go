package state

import (
	"testing"

	"github.com/MetalBlockchain/antelopevm/core/account"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
)

/* func TestReadAndWrite(t *testing.T) {
	baseDb := memdb.New()
	accountState := NewAccountState(baseDb)
	account := &core.Account{
		Name:       name.StringToName("glenn"),
		Privileged: false,
	}
	err := accountState.PutAccount(account)

	assert.Nil(t, err, "error when writing account")
	assert.Equal(t, core.IdType(0), account.ID, "account should have id 0")

	// Now let's read the account
	//account, err = accountState.GetAccountByName(account.Name)

	assert.Nil(t, err, "error when reading account")
	assert.Equal(t, name.StringToName("glenn"), account.Name, "name should equal glenn")
} */

/* func TestWrite(t *testing.T) {
	os.RemoveAll("/tmp/data")
	baseDb, _ := leveldb.New("/tmp/data/", nil, logging.NoLog{}, "", prometheus.NewRegistry())
	baseDb.Put([]byte("Account__id__1"), []byte{1})
	baseDb.Put([]byte("Account__id__3"), []byte{2})
	baseDb.Put([]byte("Sttt"), []byte{3})
	iterator := baseDb.NewIteratorWithStartAndPrefix([]byte("Account__id__1"), []byte("Account__id__"))
	iterator.Next()

	for iterator.Next() {
		fmt.Printf("iterator.Value(): %v\n", iterator.Value())
	}
}

func TestState(t *testing.T) {
	os.RemoveAll("/tmp/badger")
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	assert.NoError(t, err)

	state := NewState(nil, db)
	session := NewSession(state, db.NewTransaction(true))
	account1 := &core.Account{
		Name: name.StringToName("glenn"),
	}
	session.GetAccountIndex().Create(account1)
	fmt.Printf("account: %v\n", account1)
	account2 := &core.Account{
		Name: name.StringToName("eosio.token"),
	}
	session.GetAccountIndex().Create(account2)
	fmt.Printf("account: %v\n", account2)
	assert.NoError(t, err)

	result, err := session.GetAccountIndex().Find("byName", &core.Account{Name: name.StringToName("glenn")})
	assert.NoError(t, err)
	assert.Equal(t, result.Name, name.StringToName("glenn"))

	result, err = session.GetAccountIndex().Find("byName", &core.Account{Name: name.StringToName("eeeeee")})
	assert.Error(t, err)
	assert.Nil(t, result)

	result, err = session.GetAccountIndex().Find("byName", &core.Account{Name: name.StringToName("eosio.token")})
	assert.NoError(t, err)
	assert.Equal(t, result.ID, core.IdType(1))
	assert.Equal(t, result.Name, name.StringToName("eosio.token"))
}

func TestReadState(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions("/Users/glenn/.metalgo/chainData/2JYzpqxveXMoyXvPFpT8egwRFigMtUXrrYbxT5kcncCjcDnazx/NodeID-7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg"))
	assert.NoError(t, err)
	state := NewState(nil, db)
	session := NewSession(state, db.NewTransaction(false))
	iterator := session.transaction.NewIterator(badger.DefaultIteratorOptions)

	for iterator.Rewind(); iterator.Valid(); iterator.Next() {
		item := iterator.Item()
		k := item.Key()
		fmt.Printf("k: %v\n", string(k))
	}
}

func BenchmarkUpdate(t *testing.B) {
	os.RemoveAll("/tmp/badger")
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	assert.NoError(t, err)
	state := NewState(nil, db)
	account := &core.Account{
		Name:        name.StringToName("glenn"),
		AbiSequence: uint64(0),
	}
	session := NewSession(state, db.NewTransaction(true))
	err = session.GetAccountIndex().Create(account)
	assert.NoError(t, err)
	err = session.Commit()
	assert.NoError(t, err)

	session = NewSession(state, db.NewTransaction(true))

	for i := 0; i < t.N; i++ {
		now := time.Now()
		acc, err := session.GetAccountIndex().Find("byName", &core.Account{Name: name.StringToName("glenn")})
		assert.NoError(t, err)
		assert.NotNil(t, acc)
		session.GetAccountIndex().Modify(acc, func() {
			acc.AbiSequence += 1
		})
		fmt.Printf("elapsed: %s\n", time.Since(now))
	}

	session.Commit()
}
*/

func TestAccountState(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions(t.TempDir()))
	assert.NoError(t, err)
	state := NewState(nil, db)
	session := state.CreateSession(true)
	account1 := &account.Account{
		Name:        name.StringToName("glenn"),
		AbiSequence: uint64(0),
	}
	session.CreateAccount(account1)
	account2 := &account.Account{
		Name:        name.StringToName("eosio"),
		AbiSequence: uint64(0),
	}
	session.CreateAccount(account2)
	session.Commit()
	session = state.CreateSession(false)
	account, err := session.FindAccountByName(name.StringToName("glenn"))
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, account.Name, name.StringToName("glenn"))
	session.Discard()
	session = state.CreateSession(true)
	err = session.ModifyAccount(account, func() {
		account.AbiSequence = 10
	})
	assert.NoError(t, err)
	session.Commit()
	assert.Equal(t, uint64(10), account.AbiSequence)
}

func TestAccountState2(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions("/Users/glenn/.metalgo/chainData/2JYzpqxveXMoyXvPFpT8egwRFigMtUXrrYbxT5kcncCjcDnazx/NodeID-7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg"))
	assert.NoError(t, err)
	state := NewState(nil, db)
	session := state.CreateSession(false)
	account, err := session.FindAccountByName(name.StringToName("eosio.token"))
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, account.Name, name.StringToName("glenn"))
}

func BenchmarkAccountState(t *testing.B) {
	db, err := badger.Open(badger.DefaultOptions(t.TempDir()))
	assert.NoError(t, err)
	state := NewState(nil, db)
	session := state.CreateSession(true)
	account1 := &account.Account{
		Name:        name.StringToName("glenn"),
		AbiSequence: uint64(0),
	}
	session.CreateAccount(account1)
	account2 := &account.Account{
		Name:        name.StringToName("eosio"),
		AbiSequence: uint64(0),
	}
	session.CreateAccount(account2)
	session.Commit()
	session = state.CreateSession(false)
	account, err := session.FindAccountByName(name.StringToName("glenn"))
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, account.Name, name.StringToName("glenn"))
	session.Discard()
	session = state.CreateSession(true)

	for i := 0; i < t.N; i++ {
		err = session.ModifyAccount(account, func() {
			account.AbiSequence = 10
		})
		assert.NoError(t, err)
		account, err = session.FindAccountByName(name.StringToName("glenn"))
		assert.NoError(t, err)
		assert.Equal(t, uint64(10), account.AbiSequence)
	}
}
