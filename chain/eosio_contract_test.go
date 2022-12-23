package chain

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/MetalBlockchain/antelopevm/state"
	"github.com/MetalBlockchain/metalgo/database/leveldb"
	"github.com/MetalBlockchain/metalgo/database/memdb"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/utils/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

var (
	chainId = core.ChainIdType(*crypto.NewSha256String("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"))
)

func TestRLP(t *testing.T) {
	data, _ := hex.DecodeString("0000000000ea30550000000088395564010000000100025cb529867a83bfdf743d08d58450dfc9a79d3cbf52983e77866520e0a73ca03401000000010000000100025cb529867a83bfdf743d08d58450dfc9a79d3cbf52983e77866520e0a73ca03401000000")
	hexBytes := core.HexBytes(data)
	newAccount := &NewAccount{}
	err := rlp.DecodeBytes(hexBytes, newAccount)
	assert.NoError(t, err)
	assert.Equal(t, newAccount.Active.Keys[0].Key, "EOS5bKSrKawqzzqC6izsat38dLCnQFvHd2YganNAygpYEDL7gXVdk")
}

func TestNewAccount(t *testing.T) {
	os.RemoveAll("/tmp/data/")
	db, _ := leveldb.New("/tmp/data/", nil, logging.NoLog{}, "", prometheus.NewRegistry())
	state := state.NewState(nil, db)
	controller := NewController(chainId)
	defer db.Close()

	wif := "5KdDAKdZ9G424wH8csTZFhS9GkJjwz3dFzvJtxwM4BWDX2UVUEb"
	privKey, _ := ecc.NewPrivateKey(wif)
	pubKey := privKey.PublicKey()
	creator := NewAccount{
		Creator: core.AccountName(core.StringToName("eosio")),
		Name:    core.AccountName(core.StringToName("antelope")),
		Owner: core.Authority{
			Threshold: 1,
			Keys:      []core.KeyWeight{{Key: pubKey, Weight: 1}},
		},
		Active: core.Authority{
			Threshold: 1,
			Keys:      []core.KeyWeight{{Key: pubKey, Weight: 1}},
		},
	}
	buffer, _ := rlp.EncodeToBytes(&creator)
	act := core.Action{
		Account: core.AccountName(core.StringToName("eosio")),
		Name:    core.ActionName(core.StringToName("newaccount")),
		Data:    buffer,
		Authorization: []core.PermissionLevel{
			{Actor: core.AccountName(core.StringToName("eosio")), Permission: core.PermissionName(core.StringToName("active"))},
		},
	}
	trxHeader := core.TransactionHeader{
		Expiration:       core.MaxTimePointSec(),
		RefBlockNum:      4,
		RefBlockPrefix:   3832731038,
		MaxNetUsageWords: 0,
		MaxCpuUsageMS:    0,
		DelaySec:         0,
	}
	trx := core.Transaction{
		TransactionHeader:     trxHeader,
		ContextFreeActions:    []*core.Action{},
		Actions:               []*core.Action{&act},
		TransactionExtensions: []*core.Extension{},
	}
	signedTrx := core.NewSignedTransaction(&trx, []ecc.Signature{}, []core.HexBytes{})
	privateKey, _ := ecc.NewRandomPrivateKey()
	chainIdType := core.ChainIdType(*crypto.NewSha256String("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"))
	signedTrx.Sign(privateKey, &chainIdType)
	txContext := NewTransactionContext(controller, state, signedTrx, trx.ID())
	txContext.Init()
	txContext.Exec()
	txContext.Finalize()
}

func TestSetCode(t *testing.T) {
	os.RemoveAll("/tmp/data/")
	db, _ := leveldb.New("/tmp/data/", nil, logging.NoLog{}, "", prometheus.NewRegistry())
	state := state.NewState(nil, db)
	controller := NewController(chainId)
	defer db.Close()
	account := "antelope"
	createNewAccount(controller, state, account, "5KdDAKdZ9G424wH8csTZFhS9GkJjwz3dFzvJtxwM4BWDX2UVUEb")
	code, _ := ioutil.ReadFile("../wasm/testdata/hello.wasm")
	setCode := SetCode{
		Account:   core.AccountName(core.StringToName(account)),
		VmType:    0,
		VmVersion: 0,
		Code:      code,
	}
	buffer, _ := rlp.EncodeToBytes(&setCode)
	act := core.Action{
		Account: core.AccountName(core.StringToName(account)),
		Name:    core.ActionName(core.StringToName("setcode")),
		Data:    buffer,
		Authorization: []core.PermissionLevel{
			//common.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
			{Actor: core.AccountName(core.StringToName(account)), Permission: core.PermissionName(core.StringToName("active"))},
		},
	}
	trxHeader := core.TransactionHeader{
		Expiration:       core.MaxTimePointSec(),
		RefBlockNum:      4,
		RefBlockPrefix:   3832731038,
		MaxNetUsageWords: 0,
		MaxCpuUsageMS:    0,
		DelaySec:         0,
	}
	trx := core.Transaction{
		TransactionHeader:     trxHeader,
		ContextFreeActions:    []*core.Action{},
		Actions:               []*core.Action{&act},
		TransactionExtensions: []*core.Extension{},
	}
	signedTrx := core.NewSignedTransaction(&trx, []ecc.Signature{}, []core.HexBytes{})
	wif := "5KdDAKdZ9G424wH8csTZFhS9GkJjwz3dFzvJtxwM4BWDX2UVUEb"
	privKey, _ := ecc.NewPrivateKey(wif)
	chainIdType := core.ChainIdType(*crypto.NewSha256String("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"))
	signedTrx.Sign(privKey, &chainIdType)
	packedTrx, _ := core.NewPackedTransactionFromSignedTransaction(*signedTrx, core.CompressionNone)
	_, err := controller.PushTransaction(state, *packedTrx)
	assert.Nil(t, err, "transaction should not fail")
}

func TestUpdateAuth(t *testing.T) {
	db := memdb.New()
	state := state.NewState(nil, db)
	controller := NewController(chainId)
	defer db.Close()
	account := "antelope"
	account2 := "glenn"
	createNewAccount(controller, state, account, "5KdDAKdZ9G424wH8csTZFhS9GkJjwz3dFzvJtxwM4BWDX2UVUEb")
	createNewAccount(controller, state, account2, "5KdDAKdZ9G424wH8csTZFhS9GkJjwz3dFzvJtxwM4BWDX2UVUEb")
	updateAuth := UpdateAuth{
		Account:    core.AccountName(core.StringToName(account)),
		Permission: core.PermissionName(core.StringToName("active")),
		Parent:     core.PermissionName(core.StringToName("owner")),
		Auth: core.Authority{
			Threshold: 2,
			Accounts: []core.PermissionLevelWeight{
				{
					Permission: core.PermissionLevel{Actor: core.AccountName(core.StringToName(account2)), Permission: core.PermissionName(core.StringToName("active"))},
					Weight:     2,
				},
			},
		},
	}
	buffer, _ := rlp.EncodeToBytes(&updateAuth)
	act := core.Action{
		Account: core.AccountName(core.StringToName(account)),
		Name:    core.ActionName(core.StringToName("updateauth")),
		Data:    buffer,
		Authorization: []core.PermissionLevel{
			//common.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
			{Actor: core.AccountName(core.StringToName(account)), Permission: core.PermissionName(core.StringToName("active"))},
		},
	}
	trxHeader := core.TransactionHeader{
		Expiration:       core.MaxTimePointSec(),
		RefBlockNum:      4,
		RefBlockPrefix:   3832731038,
		MaxNetUsageWords: 0,
		MaxCpuUsageMS:    0,
		DelaySec:         0,
	}
	trx := core.Transaction{
		TransactionHeader:     trxHeader,
		ContextFreeActions:    []*core.Action{},
		Actions:               []*core.Action{&act},
		TransactionExtensions: []*core.Extension{},
	}
	signedTrx := core.NewSignedTransaction(&trx, []ecc.Signature{}, []core.HexBytes{})
	wif := "5KdDAKdZ9G424wH8csTZFhS9GkJjwz3dFzvJtxwM4BWDX2UVUEb"
	privKey, _ := ecc.NewPrivateKey(wif)
	chainIdType := core.ChainIdType(*crypto.NewSha256String("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"))
	signedTrx.Sign(privKey, &chainIdType)
	packedTrx, _ := core.NewPackedTransactionFromSignedTransaction(*signedTrx, core.CompressionNone)
	_, err := controller.PushTransaction(state, *packedTrx)
	assert.Nil(t, err, "transaction should not fail")
}

func TestLinkAndUnlinkAuth(t *testing.T) {
	wif := "5KdDAKdZ9G424wH8csTZFhS9GkJjwz3dFzvJtxwM4BWDX2UVUEb"
	db := memdb.New()
	state := state.NewState(nil, db)
	controller := NewController(chainId)
	account := "antelope"
	account2 := "hello"
	createNewAccount(controller, state, account, "5KdDAKdZ9G424wH8csTZFhS9GkJjwz3dFzvJtxwM4BWDX2UVUEb")
	createNewAccount(controller, state, account2, "5KdDAKdZ9G424wH8csTZFhS9GkJjwz3dFzvJtxwM4BWDX2UVUEb")
	code, _ := ioutil.ReadFile("../wasm/testdata/hello.wasm")
	err := setCode(controller, state, "5KdDAKdZ9G424wH8csTZFhS9GkJjwz3dFzvJtxwM4BWDX2UVUEb", code, account2)
	assert.Nil(t, err, "setCode failed")
	err = linkAuth(controller, state, wif, account, account2, "hi", "active")
	assert.Nil(t, err, "linkAuth failed")
	err = unlinkAuth(controller, state, wif, account, account2, "hi")
	assert.Nil(t, err, "unlinkAuth failed")
}

func TestDeleteAuth(t *testing.T) {
	id := ids.ID{'a', 'n', 't', 'e', 'l', 'o', 'p', 'e', 'v', 'm'}
	fmt.Printf("VM ID %s", id.String())
	wif := "5KdDAKdZ9G424wH8csTZFhS9GkJjwz3dFzvJtxwM4BWDX2UVUEb"
	db := memdb.New()
	state := state.NewState(nil, db)
	controller := NewController(chainId)
	account := "antelope"
	account2 := "hello"
	createNewAccount(controller, state, account, "5KdDAKdZ9G424wH8csTZFhS9GkJjwz3dFzvJtxwM4BWDX2UVUEb")
	createNewAccount(controller, state, account2, "5KdDAKdZ9G424wH8csTZFhS9GkJjwz3dFzvJtxwM4BWDX2UVUEb")
	err := updateAuth(controller, state, "5KdDAKdZ9G424wH8csTZFhS9GkJjwz3dFzvJtxwM4BWDX2UVUEb", account, "newperm", "active", core.Authority{
		Threshold: 1,
		Keys: []core.KeyWeight{{
			Key:    *ecc.NewPublicKeyNil(),
			Weight: 1,
		}},
	})
	assert.Nil(t, err, "updateAuth failed")
	code, _ := ioutil.ReadFile("../wasm/testdata/hello.wasm")
	err = setCode(controller, state, "5KdDAKdZ9G424wH8csTZFhS9GkJjwz3dFzvJtxwM4BWDX2UVUEb", code, account2)
	assert.Nil(t, err, "setCode failed")
	err = linkAuth(controller, state, wif, account, account2, "hi", "newperm")
	assert.Nil(t, err, "linkAuth failed")
	err = deleteAuth(controller, state, wif, account, "newperm")
	assert.Error(t, err, "cannot delete a linked authority, remove the links first")
	err = unlinkAuth(controller, state, wif, account, account2, "hi")
	assert.Nil(t, err, "unlinkAuth failed")
	err = deleteAuth(controller, state, wif, "antelope", "newperm")
	assert.Nil(t, err, "deleteAuth failed")
}

func setCode(controller *Controller, st state.State, privateKeyWif string, code []byte, account string) error {
	privateKey, _ := ecc.NewPrivateKey(privateKeyWif)
	setCode := SetCode{
		Account:   core.AccountName(core.StringToName(account)),
		VmType:    0,
		VmVersion: 0,
		Code:      code,
	}
	buffer, _ := rlp.EncodeToBytes(&setCode)
	action := core.Action{
		Account: core.AccountName(core.StringToName(account)),
		Name:    core.ActionName(core.StringToName("setcode")),
		Data:    buffer,
		Authorization: []core.PermissionLevel{
			{Actor: core.AccountName(core.StringToName(account)), Permission: core.PermissionName(core.StringToName("active"))},
		},
	}
	trxHeader := core.TransactionHeader{
		Expiration:       core.MaxTimePointSec(),
		RefBlockNum:      4,
		RefBlockPrefix:   3832731038,
		MaxNetUsageWords: 0,
		MaxCpuUsageMS:    0,
		DelaySec:         0,
	}
	trx := core.Transaction{
		TransactionHeader:     trxHeader,
		ContextFreeActions:    []*core.Action{},
		Actions:               []*core.Action{&action},
		TransactionExtensions: []*core.Extension{},
	}
	signedTrx := core.NewSignedTransaction(&trx, []ecc.Signature{}, []core.HexBytes{})
	signedTrx.Sign(privateKey, &chainId)
	packedTrx, _ := core.NewPackedTransactionFromSignedTransaction(*signedTrx, core.CompressionNone)
	_, err := controller.PushTransaction(st, *packedTrx)

	if err != nil {
		return err
	}

	return nil
}

func updateAuth(controller *Controller, st state.State, privateKeyWif string, account string, permission string, parent string, auth core.Authority) error {
	privateKey, _ := ecc.NewPrivateKey(privateKeyWif)
	updateAuth := UpdateAuth{
		Account:    core.AccountName(core.StringToName(account)),
		Permission: core.PermissionName(core.StringToName(permission)),
		Parent:     core.PermissionName(core.StringToName(parent)),
		Auth:       auth,
	}
	buffer, _ := rlp.EncodeToBytes(&updateAuth)
	action := core.Action{
		Account: core.AccountName(core.StringToName("eosio")),
		Name:    core.ActionName(core.StringToName("updateauth")),
		Data:    buffer,
		Authorization: []core.PermissionLevel{
			{Actor: core.AccountName(core.StringToName(account)), Permission: core.PermissionName(core.StringToName("active"))},
		},
	}
	trxHeader := core.TransactionHeader{
		Expiration:       core.MaxTimePointSec(),
		RefBlockNum:      4,
		RefBlockPrefix:   3832731038,
		MaxNetUsageWords: 0,
		MaxCpuUsageMS:    0,
		DelaySec:         0,
	}
	trx := core.Transaction{
		TransactionHeader:     trxHeader,
		ContextFreeActions:    []*core.Action{},
		Actions:               []*core.Action{&action},
		TransactionExtensions: []*core.Extension{},
	}
	signedTrx := core.NewSignedTransaction(&trx, []ecc.Signature{}, []core.HexBytes{})
	signedTrx.Sign(privateKey, &chainId)
	packedTrx, _ := core.NewPackedTransactionFromSignedTransaction(*signedTrx, core.CompressionNone)
	_, err := controller.PushTransaction(st, *packedTrx)

	if err != nil {
		return err
	}

	return nil
}

func deleteAuth(controller *Controller, st state.State, privateKeyWif string, account string, permission string) error {
	privateKey, _ := ecc.NewPrivateKey(privateKeyWif)
	deleteAuth := DeleteAuth{
		Account:    core.AccountName(core.StringToName(account)),
		Permission: core.PermissionName(core.StringToName(permission)),
	}
	buffer, _ := rlp.EncodeToBytes(&deleteAuth)
	action := core.Action{
		Account: core.AccountName(core.StringToName("eosio")),
		Name:    core.ActionName(core.StringToName("deleteauth")),
		Data:    buffer,
		Authorization: []core.PermissionLevel{
			{Actor: core.AccountName(core.StringToName(account)), Permission: core.PermissionName(core.StringToName("active"))},
		},
	}
	trxHeader := core.TransactionHeader{
		Expiration:       core.MaxTimePointSec(),
		RefBlockNum:      4,
		RefBlockPrefix:   3832731038,
		MaxNetUsageWords: 0,
		MaxCpuUsageMS:    0,
		DelaySec:         0,
	}
	trx := core.Transaction{
		TransactionHeader:     trxHeader,
		ContextFreeActions:    []*core.Action{},
		Actions:               []*core.Action{&action},
		TransactionExtensions: []*core.Extension{},
	}
	signedTrx := core.NewSignedTransaction(&trx, []ecc.Signature{}, []core.HexBytes{})
	signedTrx.Sign(privateKey, &chainId)
	packedTrx, _ := core.NewPackedTransactionFromSignedTransaction(*signedTrx, core.CompressionNone)
	_, err := controller.PushTransaction(st, *packedTrx)

	if err != nil {
		return err
	}

	return nil
}

func linkAuth(controller *Controller, st state.State, privateKeyWif string, account string, contract string, actionName string, permission string) error {
	privateKey, _ := ecc.NewPrivateKey(privateKeyWif)
	linkAuth := LinkAuth{
		Account:     core.AccountName(core.StringToName(account)),
		Code:        core.AccountName(core.StringToName(contract)),
		Type:        core.ActionName(core.StringToName(actionName)),
		Requirement: core.PermissionName(core.StringToName(permission)),
	}
	buffer, _ := rlp.EncodeToBytes(&linkAuth)
	action := core.Action{
		Account: core.AccountName(core.StringToName("eosio")),
		Name:    core.ActionName(core.StringToName("linkauth")),
		Data:    buffer,
		Authorization: []core.PermissionLevel{
			{Actor: core.AccountName(core.StringToName(account)), Permission: core.PermissionName(core.StringToName("active"))},
		},
	}
	trxHeader := core.TransactionHeader{
		Expiration:       core.MaxTimePointSec(),
		RefBlockNum:      4,
		RefBlockPrefix:   3832731038,
		MaxNetUsageWords: 0,
		MaxCpuUsageMS:    0,
		DelaySec:         0,
	}
	trx := core.Transaction{
		TransactionHeader:     trxHeader,
		ContextFreeActions:    []*core.Action{},
		Actions:               []*core.Action{&action},
		TransactionExtensions: []*core.Extension{},
	}
	signedTrx := core.NewSignedTransaction(&trx, []ecc.Signature{}, []core.HexBytes{})
	signedTrx.Sign(privateKey, &chainId)
	packedTrx, err := core.NewPackedTransactionFromSignedTransaction(*signedTrx, core.CompressionNone)
	if err != nil {
		return err
	}
	_, err = controller.PushTransaction(st, *packedTrx)

	if err != nil {
		return err
	}

	return nil
}

func unlinkAuth(controller *Controller, st state.State, privateKeyWif string, account string, contract string, actionName string) error {
	privateKey, _ := ecc.NewPrivateKey(privateKeyWif)
	unlinkAuth := UnLinkAuth{
		Account: core.AccountName(core.StringToName(account)),
		Code:    core.AccountName(core.StringToName(contract)),
		Type:    core.ActionName(core.StringToName(actionName)),
	}
	buffer, _ := rlp.EncodeToBytes(&unlinkAuth)
	action := core.Action{
		Account: core.AccountName(core.StringToName("eosio")),
		Name:    core.ActionName(core.StringToName("unlinkauth")),
		Data:    buffer,
		Authorization: []core.PermissionLevel{
			{Actor: core.AccountName(core.StringToName(account)), Permission: core.PermissionName(core.StringToName("active"))},
		},
	}
	trxHeader := core.TransactionHeader{
		Expiration:       core.MaxTimePointSec(),
		RefBlockNum:      4,
		RefBlockPrefix:   3832731038,
		MaxNetUsageWords: 0,
		MaxCpuUsageMS:    0,
		DelaySec:         0,
	}
	trx := core.Transaction{
		TransactionHeader:     trxHeader,
		ContextFreeActions:    []*core.Action{},
		Actions:               []*core.Action{&action},
		TransactionExtensions: []*core.Extension{},
	}
	signedTrx := core.NewSignedTransaction(&trx, []ecc.Signature{}, []core.HexBytes{})
	signedTrx.Sign(privateKey, &chainId)
	packedTrx, err := core.NewPackedTransactionFromSignedTransaction(*signedTrx, core.CompressionNone)
	if err != nil {
		return err
	}
	_, err = controller.PushTransaction(st, *packedTrx)

	if err != nil {
		return err
	}

	return nil
}

func createNewAccount(controller *Controller, st state.State, name string, privateKeyWif string) {
	privKey, _ := ecc.NewPrivateKey(privateKeyWif)
	pubKey := privKey.PublicKey()
	creator := NewAccount{
		Creator: core.AccountName(core.StringToName("eosio")),
		Name:    core.AccountName(core.StringToName(name)),
		Owner: core.Authority{
			Threshold: 1,
			Keys:      []core.KeyWeight{{Key: pubKey, Weight: 1}},
		},
		Active: core.Authority{
			Threshold: 1,
			Keys:      []core.KeyWeight{{Key: pubKey, Weight: 1}},
		},
	}
	buffer, _ := rlp.EncodeToBytes(&creator)
	act := core.Action{
		Account: core.AccountName(core.StringToName("eosio")),
		Name:    core.ActionName(core.StringToName("newaccount")),
		Data:    buffer,
		Authorization: []core.PermissionLevel{
			{Actor: core.AccountName(core.StringToName("eosio")), Permission: core.PermissionName(core.StringToName("active"))},
		},
	}
	trxHeader := core.TransactionHeader{
		Expiration:       core.MaxTimePointSec(),
		RefBlockNum:      4,
		RefBlockPrefix:   3832731038,
		MaxNetUsageWords: 0,
		MaxCpuUsageMS:    0,
		DelaySec:         0,
	}
	trx := core.Transaction{
		TransactionHeader:     trxHeader,
		ContextFreeActions:    []*core.Action{},
		Actions:               []*core.Action{&act},
		TransactionExtensions: []*core.Extension{},
	}
	signedTrx := core.NewSignedTransaction(&trx, []ecc.Signature{}, []core.HexBytes{})
	privateKey, _ := ecc.NewRandomPrivateKey()
	chainIdType := core.ChainIdType(*crypto.NewSha256String("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"))
	signedTrx.Sign(privateKey, &chainIdType)
	txContext := NewTransactionContext(controller, st, signedTrx, trx.ID())
	txContext.Init()
	txContext.Exec()
	txContext.Finalize()
}
