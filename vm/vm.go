package vm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MetalBlockchain/antelopevm/chain"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/mempool"
	"github.com/MetalBlockchain/antelopevm/state"
	"github.com/MetalBlockchain/metalgo/database/manager"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow"
	"github.com/MetalBlockchain/metalgo/snow/choices"
	"github.com/MetalBlockchain/metalgo/snow/consensus/snowman"
	"github.com/MetalBlockchain/metalgo/snow/engine/common"
	"github.com/MetalBlockchain/metalgo/snow/engine/snowman/block"
	"github.com/MetalBlockchain/metalgo/utils"
	"github.com/MetalBlockchain/metalgo/version"

	log "github.com/inconshreveable/log15"
)

const (
	dataLen = 32
	Name    = "antelopevm"
)

var (
	errNoPendingBlocks = errors.New("there is no block to propose")
	Version            = &version.Semantic{
		Major: 0,
		Minor: 0,
		Patch: 1,
	}

	_ block.ChainVM = &VM{}
	_ state.VM      = &VM{}
)

type VM struct {
	// The context of this vm
	ctx       *snow.Context
	dbManager manager.Manager

	// State of this VM
	state state.State

	// ID of the preferred block
	preferred ids.ID

	// channel to send messages to the consensus engine
	toEngine chan<- common.Message

	// Proposed pieces of data that haven't been put into a block and proposed yet
	mempool *mempool.Mempool

	// Block ID --> Block
	// Each element is a block that passed verification but
	// hasn't yet been accepted/rejected
	verifiedBlocks map[ids.ID]*state.Block

	// Indicates that this VM has finised bootstrapping for the chain
	bootstrapped utils.AtomicBool

	controller *chain.Controller
	builder    BlockBuilder

	stop chan struct{}

	builderStop chan struct{}
	doneBuild   chan struct{}
	doneGossip  chan struct{}
}

// Initialize this vm
// [ctx] is this vm's context
// [dbManager] is the manager of this vm's database
// [toEngine] is used to notify the consensus engine that new blocks are
//
//	ready to be added to consensus
//
// The data in the genesis block is [genesisData]
func (vm *VM) Initialize(
	ctx context.Context,
	chainCtx *snow.Context,
	dbManager manager.Manager,
	genesisData []byte,
	upgradeData []byte,
	configData []byte,
	toEngine chan<- common.Message,
	_ []*common.Fx,
	_ common.AppSender,
) error {
	log.Info("Initializing Antelope VM")

	vm.dbManager = dbManager
	vm.ctx = chainCtx
	vm.toEngine = toEngine
	vm.verifiedBlocks = make(map[ids.ID]*state.Block)

	// Create new state and controller
	chainId := types.ChainIdType(*crypto.NewSha256String("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"))
	vm.state = state.NewState(vm, vm.dbManager.Current().Database)
	vm.controller = chain.NewController(vm.state, chainId)
	vm.mempool = mempool.New(100)
	vm.builder = vm.NewBlockBuilder()

	// Init channels
	vm.stop = make(chan struct{})
	vm.builderStop = make(chan struct{})
	vm.doneBuild = make(chan struct{})
	vm.doneGossip = make(chan struct{})

	// Initialize genesis
	if err := vm.initGenesis(genesisData); err != nil {
		return err
	}

	// Get last accepted
	lastAccepted, err := vm.state.GetLastAccepted()
	if err != nil {
		return fmt.Errorf("failed to get last accepted block: %s", err)
	}

	log.Info("initializing last accepted block", "lastAccepted", lastAccepted)

	// Build off the most recently accepted block
	if err := vm.SetPreference(ctx, lastAccepted); err != nil {
		return err
	}

	go vm.builder.Build()

	return nil
}

// Initializes Genesis if required
func (vm *VM) initGenesis(genesisData []byte) error {
	stateInitialized, err := vm.state.IsInitialized()

	if err != nil {
		return err
	}

	if stateInitialized {
		return nil
	}

	genesisFile := chain.ParseGenesisData(genesisData)

	// Initialize the genesis state
	if err := vm.controller.InitGenesis(genesisFile); err != nil {
		return err
	}

	// Create the genesis block
	// Timestamp of genesis block is 0. It has no parent.
	genesisBlock, err := vm.NewBlock(ids.Empty, 1, genesisData, genesisFile.InitialTimeStamp)

	if err != nil {
		log.Error("error while creating genesis block: %v", err)
		return err
	}

	// Put genesis block to state
	if err := vm.state.PutBlock(genesisBlock); err != nil {
		log.Error("error while saving genesis block: %v", err)
		return err
	}

	// Accept the genesis block
	// Sets [vm.lastAccepted] and [vm.preferred]
	if err := genesisBlock.Accept(context.Background()); err != nil {
		return fmt.Errorf("error accepting genesis block: %w", err)
	}

	// Mark this vm's state as initialized, so we can skip initGenesis in further restarts
	if err := vm.state.SetInitialized(); err != nil {
		return fmt.Errorf("error while setting db to initialized: %w", err)
	}

	// Flush VM's database to underlying db
	return vm.state.Commit()
}

// CreateHandlers returns a map where:
// Keys: The path extension for this VM's API (empty in this case)
// Values: The handler for the API
func (vm *VM) CreateHandlers(ctx context.Context) (map[string]*common.HTTPHandler, error) {
	service := &Service{vm}

	return map[string]*common.HTTPHandler{
		"/v1/chain/get_info": {
			Handler: NewRequestHandler(service.GetInfo),
		},
		"/v1/chain/get_block": {
			Handler: NewRequestHandler(service.GetBlock),
		},
		"/v1/chain/get_block_info": {
			Handler: NewRequestHandler(service.GetBlockInfo),
		},
		"/v1/chain/get_required_keys": {
			Handler: NewRequestHandler(service.GetRequiredKeys),
		},
		"/v1/chain/send_transaction": {
			Handler: NewRequestHandler(service.PushTransaction),
		},
	}, nil
}

// CreateStaticHandlers returns a map where:
// Keys: The path extension for this VM's static API
// Values: The handler for that static API
func (vm *VM) CreateStaticHandlers(ctx context.Context) (map[string]*common.HTTPHandler, error) {
	return nil, nil
}

// Health implements the common.VM interface
func (vm *VM) HealthCheck(ctx context.Context) (interface{}, error) {
	return nil, nil
}

// BuildBlock returns a block that this vm wants to add to consensus
func (vm *VM) BuildBlock(ctx context.Context) (snowman.Block, error) {
	if vm.mempool.Len() == 0 { // There is no block to be built
		return nil, errNoPendingBlocks
	}

	// Get transaction
	tx := vm.mempool.Pop()

	// Try to run the transaction
	signedTx, _ := tx.GetSignedTransaction()

	_, err := vm.controller.PushTransaction(*signedTx)

	if err != nil {
		return nil, fmt.Errorf("couldn't push transaction to controller: %w", err)
	}

	// Gets Preferred Block
	preferredBlock, err := vm.getBlock(vm.preferred)

	if err != nil {
		return nil, fmt.Errorf("couldn't get preferred block: %w", err)
	}

	preferredHeight := preferredBlock.Height()

	// Build the block with preferred height
	newBlock, err := vm.NewBlock(vm.preferred, preferredHeight+1, []byte{1}, types.Now())

	if err != nil {
		return nil, fmt.Errorf("couldn't build block: %w", err)
	}

	// Verifies block
	if err := newBlock.Verify(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to verify block: %s", err)
	}

	return newBlock, nil
}

// GetBlock implements the snowman.ChainVM interface
func (vm *VM) GetBlock(ctx context.Context, blkID ids.ID) (snowman.Block, error) {
	block, err := vm.getBlock(blkID)

	if err != nil {
		return nil, fmt.Errorf("failed to get block %s: %s", blkID, err)
	}

	return block, nil
}

func (vm *VM) getBlock(blkID ids.ID) (*state.Block, error) {
	// If block is in memory, return it.
	if blk, exists := vm.verifiedBlocks[blkID]; exists {
		return blk, nil
	}

	return vm.state.GetBlock(blkID)
}

// LastAccepted returns the block most recently accepted
func (vm *VM) LastAccepted(ctx context.Context) (ids.ID, error) {
	return vm.state.GetLastAccepted()
}

// ParseBlock parses [bytes] to a snowman.Block
// This function is used by the vm's state to unmarshal blocks saved in state
// and by the consensus layer when it receives the byte representation of a block
// from another node
func (vm *VM) ParseBlock(ctx context.Context, bytes []byte) (snowman.Block, error) {
	// A new empty block
	block := &state.Block{}

	// Unmarshal the byte repr. of the block into our empty block
	_, err := state.Codec.Unmarshal(bytes, block)
	if err != nil {
		return nil, err
	}

	// Initialize the block
	block.Initialize(vm, choices.Processing)

	if blk, err := vm.getBlock(block.ID()); err == nil {
		// If we have seen this block before, return it with the most up-to-date
		// info
		return blk, nil
	}

	// Return the block
	return block, nil
}

// NewBlock returns a new Block where:
// - the block's parent is [parentID]
// - the block's data is [data]
// - the block's timestamp is [timestamp]
func (vm *VM) NewBlock(parentID ids.ID, height uint64, data []byte, timestamp types.TimePoint) (*state.Block, error) {
	block := &state.Block{
		BlockHeader: state.BlockHeader{
			Created:       timestamp,
			Producer:      types.N("eosio"),
			Confirmed:     1,
			PreviousBlock: parentID,
			Index:         height,
		},
		Transactions: []types.TransactionReceipt{},
	}

	// Initialize the block by providing it with its byte representation
	// and a reference to this VM
	block.Initialize(vm, choices.Processing)

	return block, nil
}

// Shutdown this vm
func (vm *VM) Shutdown(ctx context.Context) error {
	if vm.state == nil {
		return nil
	}

	return vm.state.Close() // close versionDB
}

// SetPreference sets the block with ID [ID] as the preferred block
func (vm *VM) SetPreference(ctx context.Context, id ids.ID) error {
	vm.preferred = id
	return nil
}

func (vm *VM) Verified(block *state.Block) error {
	vm.verifiedBlocks[block.ID()] = block
	return nil
}

func (vm *VM) Accepted(block *state.Block) error {
	block.SetStatus(choices.Accepted) // Change state of this block
	blkID := block.ID()

	// Persist data
	if err := vm.state.PutBlock(block); err != nil {
		return fmt.Errorf("failed to insert block: %s", err)
	}

	if err := vm.state.SetLastAccepted(blkID); err != nil {
		return fmt.Errorf("failed to set last accepted: %s", err)
	}

	// Delete this block from verified blocks as it's accepted
	delete(vm.verifiedBlocks, block.ID())

	return vm.state.Commit()
}

// SetState sets this VM state according to given snow.State
func (vm *VM) SetState(ctx context.Context, state snow.State) error {
	switch state {
	// Engine reports it's bootstrapping
	case snow.Bootstrapping:
		return vm.onBootstrapStarted()
	case snow.NormalOp:
		// Engine reports it can start normal operations
		return vm.onNormalOperationsStarted()
	default:
		return snow.ErrUnknownState
	}
}

// onBootstrapStarted marks this VM as bootstrapping
func (vm *VM) onBootstrapStarted() error {
	vm.bootstrapped.SetValue(false)
	return nil
}

// onNormalOperationsStarted marks this VM as bootstrapped
func (vm *VM) onNormalOperationsStarted() error {
	// No need to set it again
	if vm.bootstrapped.GetValue() {
		return nil
	}
	vm.bootstrapped.SetValue(true)
	return nil
}

// Returns this VM's version
func (vm *VM) Version(ctx context.Context) (string, error) {
	return Version.String(), nil
}

func (vm *VM) Connected(ctx context.Context, id ids.NodeID, nodeVersion *version.Application) error {
	return nil // noop
}

func (vm *VM) Disconnected(ctx context.Context, id ids.NodeID) error {
	return nil // noop
}

// This VM doesn't (currently) have any app-specific messages
func (vm *VM) AppGossip(ctx context.Context, nodeID ids.NodeID, msg []byte) error {
	return nil
}

// This VM doesn't (currently) have any app-specific messages
func (vm *VM) AppRequest(ctx context.Context, nodeID ids.NodeID, requestID uint32, time time.Time, request []byte) error {
	return nil
}

// This VM doesn't (currently) have any app-specific messages
func (vm *VM) AppResponse(ctx context.Context, nodeID ids.NodeID, requestID uint32, response []byte) error {
	return nil
}

// This VM doesn't (currently) have any app-specific messages
func (vm *VM) AppRequestFailed(ctx context.Context, nodeID ids.NodeID, requestID uint32) error {
	return nil
}

// implements "snowmanblock.ChainVM.commom.VM.AppHandler"
func (vm *VM) CrossChainAppRequest(ctx context.Context, chainID ids.ID, requestID uint32, deadline time.Time, request []byte) error {
	// (currently) no app-specific messages
	return nil
}

// implements "snowmanblock.ChainVM.commom.VM.AppHandler"
func (vm *VM) CrossChainAppRequestFailed(ctx context.Context, chainID ids.ID, requestID uint32) error {
	// (currently) no app-specific messages
	return nil
}

// implements "snowmanblock.ChainVM.commom.VM.AppHandler"
func (vm *VM) CrossChainAppResponse(ctx context.Context, chainID ids.ID, requestID uint32, request []byte) error {
	// (currently) no app-specific messages
	return nil
}
