package chain_test

import (
	"io/ioutil"
	"testing"

	"github.com/MetalBlockchain/antelopevm/chain"
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/authority"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/MetalBlockchain/antelopevm/state"
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
)

func TestDeploySystemContract(b *testing.T) {
	privateKey, _ := ecc.NewRandomPrivateKey()
	genesisConfig := chain.GenesisFile{
		InitialTimeStamp: core.Now(),
		InitialKey:       privateKey.PublicKey(),
	}
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	assert.NoError(b, err)
	s := state.NewState(nil, db)
	chainId := core.ChainIdType(*crypto.Hash256String("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"))
	controller := chain.NewController(chainId, s)
	session := s.CreateSession(true)
	err = controller.InitGenesis(session, &genesisConfig)
	assert.NoError(b, err)
	session.Commit()

	transaction, err := buildSetCodeTransaction(*privateKey, chainId)
	assert.NoError(b, err)
	session = s.CreateSession(true)
	block := state.Block{
		Header: core.BlockHeader{Index: 1},
	}
	trace, err := controller.PushTransaction(*transaction, &block, session)
	assert.NoError(b, err)
	assert.NotNil(b, trace)
	session.Discard()
}

func buildSetCodeTransaction(privateKey ecc.PrivateKey, chainId core.ChainIdType) (*core.PackedTransaction, error) {
	code, err := ioutil.ReadFile("../contracts/eosio.system/eosio.system.wasm")

	if err != nil {
		return nil, err
	}

	abi, err := ioutil.ReadFile("../contracts/eosio.system/eosio.system.abi")

	if err != nil {
		return nil, err
	}

	abiRlp, _ := rlp.EncodeToBytes(abi)

	setCode := chain.SetCode{
		Account:   name.StringToName("eosio"),
		VmType:    0,
		VmVersion: 0,
		Code:      code,
	}
	setCodeBuffer, err := rlp.EncodeToBytes(setCode)

	if err != nil {
		return nil, err
	}

	setAbi := chain.SetAbi{
		Account: name.StringToName("eosio"),
		Abi:     abiRlp,
	}
	setAbiBuffer, err := rlp.EncodeToBytes(setAbi)

	if err != nil {
		return nil, err
	}

	transaction := core.Transaction{
		TransactionHeader: core.TransactionHeader{
			Expiration:       core.MaxTimePointSec(),
			RefBlockNum:      0,
			RefBlockPrefix:   3832731038,
			MaxNetUsageWords: 0,
			MaxCpuUsageMS:    0,
			DelaySec:         0,
		},
		ContextFreeActions: []*core.Action{},
		Actions: []*core.Action{
			{
				Account: name.StringToName("eosio"),
				Name:    name.StringToName("setcode"),
				Data:    setCodeBuffer,
				Authorization: []authority.PermissionLevel{
					{Actor: name.StringToName("eosio"), Permission: name.StringToName("active")},
				},
			},
			{
				Account: name.StringToName("eosio"),
				Name:    name.StringToName("setabi"),
				Data:    setAbiBuffer,
				Authorization: []authority.PermissionLevel{
					{Actor: name.StringToName("eosio"), Permission: name.StringToName("active")},
				},
			},
		},
		TransactionExtensions: []*core.Extension{},
	}
	signedTransaction := core.NewSignedTransaction(&transaction, []ecc.Signature{}, []core.HexBytes{})
	signedTransaction.Sign(&privateKey, &chainId)

	return core.NewPackedTransactionFromSignedTransaction(*signedTransaction, core.CompressionNone)
}
