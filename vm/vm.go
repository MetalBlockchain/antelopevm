package vm

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"

	"github.com/MetalBlockchain/antelopevm/chain"
	chainBlock "github.com/MetalBlockchain/antelopevm/chain/block"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	chainTime "github.com/MetalBlockchain/antelopevm/chain/time"
	"github.com/MetalBlockchain/antelopevm/chain/transaction"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/mempool"
	"github.com/MetalBlockchain/antelopevm/state"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/MetalBlockchain/metalgo/database/manager"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow"
	"github.com/MetalBlockchain/metalgo/snow/choices"
	"github.com/MetalBlockchain/metalgo/snow/consensus/snowman"
	"github.com/MetalBlockchain/metalgo/snow/engine/common"
	"github.com/MetalBlockchain/metalgo/snow/engine/snowman/block"
	"github.com/MetalBlockchain/metalgo/version"
	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"

	// Initializes service plugins
	_ "github.com/MetalBlockchain/antelopevm/vm/service/chain_api_plugin"

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
	_ service.VM    = &VM{}
)

type VM struct {
	// The context of this vm
	ctx       *snow.Context
	dbManager manager.Manager

	// State of this VM
	db         *badger.DB
	dbPath     string
	state      *state.State
	controller *chain.Controller

	// ID of the preferred block
	preferred ids.ID

	// channel to send messages to the consensus engine
	toEngine chan<- common.Message

	// Proposed pieces of data that haven't been put into a block and proposed yet
	mempool *mempool.Mempool

	// Block ID --> Block
	// Each element is a block that passed verification but
	// hasn't yet been accepted/rejected
	verifiedBlocks map[chainBlock.BlockHash]*state.Block

	// Indicates that this VM has finised bootstrapping for the chain
	bootstrapped bool
	builder      BlockBuilder

	chainId crypto.Sha256

	stop chan struct{}

	builderStop chan struct{}
	doneBuild   chan struct{}
	doneGossip  chan struct{}

	cpuProfiler *os.File
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
	log.Info("initializing Antelope VM", "version", "0.0.1")
	vm.dbManager = dbManager
	vm.ctx = chainCtx
	vm.toEngine = toEngine
	vm.verifiedBlocks = make(map[chainBlock.BlockHash]*state.Block)

	// Create new state and controller
	vm.chainId = types.ChainIdType(*crypto.NewSha256String("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"))
	vm.dbPath = filepath.Join(vm.ctx.ChainDataDir, chainCtx.NodeID.String())

	if db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true)); err == nil {
		vm.db = db
	} else {
		return err
	}

	vm.state = state.NewState(vm, vm.db)
	vm.mempool = mempool.New(100)
	vm.builder = vm.NewBlockBuilder()
	vm.controller = chain.NewController(vm.chainId, vm.state)

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
	lastAccepted, err := vm.LastAccepted(ctx)

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
	session := vm.state.CreateSession(true)
	defer session.Discard()
	stateInitialized, err := vm.state.IsInitialized()

	if err != nil {
		return err
	}

	if stateInitialized {
		return nil
	}

	genesisFile, err := chain.ParseGenesisData(genesisData)
	if err != nil {
		return err
	}

	if err := vm.controller.InitGenesis(session, genesisFile); err != nil {
		return err
	}

	// Create the genesis block
	// Timestamp of genesis block is 0. It has no parent.
	genesisBlock, err := vm.NewBlock(chainBlock.BlockHash(ids.Empty), 0, []transaction.TransactionReceipt{}, genesisFile.InitialTimeStamp)

	if err != nil {
		log.Error("error while creating genesis block: %v", err)
		return err
	}

	// Put genesis block to state
	if err := session.CreateBlock(genesisBlock); err != nil {
		log.Error("error while saving genesis block: %v", err)
		return err
	}

	if err := session.Commit(); err != nil {
		return fmt.Errorf("could not commit session: %v", err)
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
	return nil
}

// CreateHandlers returns a map where:
// Keys: The path extension for this VM's API (empty in this case)
// Values: The handler for the API
func (vm *VM) CreateHandlers(ctx context.Context) (map[string]*common.HTTPHandler, error) {
	handlers := make(map[string]*common.HTTPHandler)
	router := gin.Default()

	for path, handler := range service.GetHandlers() {
		for _, method := range handler.Methods {
			if method == http.MethodPost {
				router.POST("/ext/bc/"+vm.ctx.ChainID.String()+path, handler.HandlerFunc(vm))
			} else if method == http.MethodGet {
				router.GET("/ext/bc/"+vm.ctx.ChainID.String()+path, handler.HandlerFunc(vm))
			}
		}

		handlers[path] = &common.HTTPHandler{
			LockOptions: common.NoLock,
			Handler:     router.Handler(),
		}
	}

	return handlers, nil
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
	defer vm.builder.HandleGenerateBlock()

	newBlock, err := state.BuildBlock(vm, vm.preferred)

	if err != nil {
		return nil, fmt.Errorf("couldn't build block: %w", err)
	}

	log.Debug("block built successfully", "block", newBlock.ID())

	return newBlock, nil
}

// GetBlock implements the snowman.ChainVM interface
func (vm *VM) GetBlock(ctx context.Context, blkID ids.ID) (snowman.Block, error) {
	block, err := vm.getBlock(chainBlock.BlockHash(blkID))

	if err != nil {
		return nil, fmt.Errorf("failed to get block %s: %s", blkID, err)
	}

	return block, nil
}

func (vm *VM) getBlock(blkID chainBlock.BlockHash) (*state.Block, error) {
	// If block is in memory, return it.
	if blk, exists := vm.verifiedBlocks[blkID]; exists {
		return blk, nil
	}

	session := vm.state.CreateSession(false)
	defer session.Discard()

	return session.FindBlockByHash(blkID)
}

// LastAccepted returns the block most recently accepted
func (vm *VM) LastAccepted(ctx context.Context) (ids.ID, error) {
	session := vm.state.CreateSession(false)
	defer session.Discard()

	return session.GetLastAccepted()
}

// ParseBlock parses [bytes] to a snowman.Block
// This function is used by the vm's state to unmarshal blocks saved in state
// and by the consensus layer when it receives the byte representation of a block
// from another node
func (vm *VM) ParseBlock(ctx context.Context, bytes []byte) (snowman.Block, error) {
	// A new empty block
	block := &state.Block{}
	// Unmarshal the byte repr. of the block into our empty block
	if _, err := state.Codec.Unmarshal(bytes, block); err != nil {
		log.Error("couldn't parse block", "error", err)
		return nil, err
	}

	// Initialize the block
	block.Initialize(vm)
	block.SetStatus(choices.Processing)

	if blk, err := vm.getBlock(block.Hash); err == nil {
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
func (vm *VM) NewBlock(parentID chainBlock.BlockHash, height uint64, receipts []transaction.TransactionReceipt, timestamp chainTime.TimePoint) (*state.Block, error) {
	block := &state.Block{
		Header: chainBlock.BlockHeader{
			Timestamp: chainBlock.NewBlockTimeStampFromTimePoint(timestamp),
			Producer:  name.StringToName("eosio"),
			Confirmed: 1,
		},
		Transactions: receipts,
	}

	// Initialize the block by providing it with its byte representation
	// and a reference to this VM
	block.Initialize(vm)
	block.SetStatus(choices.Processing)
	block.Finalize()

	return block, nil
}

// Shutdown this vm
func (vm *VM) Shutdown(ctx context.Context) error {
	pprof.StopCPUProfile()
	vm.cpuProfiler.Close()
	if vm.state == nil {
		return nil
	}

	return vm.state.Close() // close versionDB
}

func (vm *VM) State() *state.State {
	return vm.state
}

// SetPreference sets the block with ID [ID] as the preferred block
func (vm *VM) SetPreference(ctx context.Context, id ids.ID) error {
	vm.preferred = id
	return nil
}

func (vm *VM) Verified(block *state.Block) error {
	vm.verifiedBlocks[block.Hash] = block
	return nil
}

func (vm *VM) Accepted(block *state.Block) error {
	session := vm.state.CreateSession(true)
	defer session.Discard()
	block.SetStatus(choices.Accepted) // Change state of this block
	blkID := block.ID()

	// Persist data
	if err := session.CreateBlock(block); err != nil {
		return fmt.Errorf("failed to insert block: %s", err)
	}

	if err := session.SetLastAccepted(blkID); err != nil {
		return fmt.Errorf("failed to set last accepted: %s", err)
	}

	if err := session.Commit(); err != nil {
		return fmt.Errorf("failed to commit session: %s", err)
	}

	// Delete this block from verified blocks as it's accepted
	delete(vm.verifiedBlocks, block.Hash)

	return nil
}

func (vm *VM) Rejected(block *state.Block) error {
	delete(vm.verifiedBlocks, block.Hash)

	return nil
}

func (vm *VM) GetMempool() *mempool.Mempool {
	return vm.mempool
}

func (vm *VM) GetStoredBlock(context context.Context, blkID ids.ID) (*state.Block, error) {
	if blk, exists := vm.verifiedBlocks[chainBlock.BlockHash(blkID)]; exists {
		return blk, nil
	}

	stBlk, err := vm.getBlock(chainBlock.BlockHash(blkID))

	if err != nil {
		log.Error("could not get stored block from DB", "id", blkID)
		return nil, fmt.Errorf("could not get stored block")
	}

	return stBlk, nil
}

func (vm *VM) ExecuteTransaction(trx *transaction.PackedTransaction, block *state.Block, session *state.Session) (*transaction.TransactionTrace, error) {
	if err := trx.UnpackTransaction(); err != nil {
		return nil, err
	}

	trxMeta, err := transaction.RecoverKeys(trx, vm.chainId, chainTime.MaxMicroseconds(), transaction.Input, 0)
	trace, err := vm.controller.PushTransaction(*trxMeta, block, session)

	if err != nil {
		log.Error("failed to execute trx", "error", err)
		return nil, err
	}

	log.Info("done processing transaction")

	return trace, nil
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
	vm.bootstrapped = false
	return nil
}

// onNormalOperationsStarted marks this VM as bootstrapped
func (vm *VM) onNormalOperationsStarted() error {
	// No need to set it again
	if vm.bootstrapped {
		return nil
	}
	vm.bootstrapped = true
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

func (vm *VM) GetState() *state.State {
	return vm.state
}

func (vm *VM) GetController() *chain.Controller {
	return vm.controller
}
