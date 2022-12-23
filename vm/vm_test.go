package vm

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/MetalBlockchain/metalgo/database/manager"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow"
	"github.com/MetalBlockchain/metalgo/snow/engine/common"
	"github.com/MetalBlockchain/metalgo/version"
	"github.com/stretchr/testify/assert"
)

func TestVMInit(t *testing.T) {
	assert := assert.New(t)
	ctx := context.TODO()
	vm, _, _, err := newTestVM()
	assert.NoError(err)
	ok, err := vm.state.IsInitialized()
	assert.NoError(err)
	assert.True(ok)
	lastAccepted, err := vm.LastAccepted(ctx)
	assert.NoError(err)
	assert.NotEqual(ids.Empty, lastAccepted)
}

func newTestVM() (*VM, *snow.Context, chan common.Message, error) {
	vm := &VM{}
	snowCtx := snow.DefaultContextTest()
	dbManager := manager.NewMemDB(&version.Semantic{
		Major: 1,
		Minor: 0,
		Patch: 0,
	})
	jsonFile, err := os.Open("../chain/genesis_test.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	msgChan := make(chan common.Message, 1)
	err = vm.Initialize(context.TODO(), snowCtx, dbManager, byteValue, nil, nil, msgChan, nil, nil)
	return vm, snowCtx, msgChan, err
}
