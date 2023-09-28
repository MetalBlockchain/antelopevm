package chain

import (
	"fmt"
	"sync"

	"github.com/MetalBlockchain/antelopevm/chain/account"
	"github.com/MetalBlockchain/antelopevm/chain/authority"
	"github.com/MetalBlockchain/antelopevm/chain/block"
	"github.com/MetalBlockchain/antelopevm/chain/fc"
	"github.com/MetalBlockchain/antelopevm/chain/global"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/time"
	"github.com/MetalBlockchain/antelopevm/chain/transaction"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/config"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/state"
	"github.com/MetalBlockchain/antelopevm/utils"
	"github.com/MetalBlockchain/antelopevm/wasm/api"
	log "github.com/inconshreveable/log15"
)

type v func(ctx ApplyContext) error

var _ api.Controller = &Controller{}

type Controller struct {
	ApplyHandlers map[string]v
	ChainId       types.ChainIdType
	State         *state.State

	ActorWhitelist    name.NameSet
	ActorBlacklist    name.NameSet
	ContractWhitelist name.NameSet
	ContractBlacklist name.NameSet
	KeyBlackist       ecc.PublicKeySet
	ReadOnly          bool

	transactionMutex sync.Mutex
}

func NewController(chainId types.ChainIdType, state *state.State) *Controller {
	controller := &Controller{
		ApplyHandlers: make(map[string]v),
		ChainId:       chainId,
		State:         state,

		ActorWhitelist:    name.NewNameSet(0),
		ActorBlacklist:    name.NewNameSet(0),
		ContractWhitelist: name.NewNameSet(0),
		ContractBlacklist: name.NewNameSet(0),
		KeyBlackist:       ecc.NewPublicKeySet(0),
		ReadOnly:          false,
	}

	// Add native functions
	controller.SetApplyHandler(name.AccountName(name.StringToName("eosio")), name.AccountName(name.StringToName("eosio")), name.ActionName(name.StringToName("newaccount")), applyEosioNewaccount)
	controller.SetApplyHandler(name.AccountName(name.StringToName("eosio")), name.AccountName(name.StringToName("eosio")), name.ActionName(name.StringToName("setcode")), applyEosioSetCode)
	controller.SetApplyHandler(name.AccountName(name.StringToName("eosio")), name.AccountName(name.StringToName("eosio")), name.ActionName(name.StringToName("setabi")), applyEosioSetAbi)
	controller.SetApplyHandler(name.AccountName(name.StringToName("eosio")), name.AccountName(name.StringToName("eosio")), name.ActionName(name.StringToName("updateauth")), applyEosioUpdateAuth)
	controller.SetApplyHandler(name.AccountName(name.StringToName("eosio")), name.AccountName(name.StringToName("eosio")), name.ActionName(name.StringToName("deleteauth")), applyEosioDeleteAuth)
	controller.SetApplyHandler(name.AccountName(name.StringToName("eosio")), name.AccountName(name.StringToName("eosio")), name.ActionName(name.StringToName("linkauth")), applyEosioLinkAuth)
	controller.SetApplyHandler(name.AccountName(name.StringToName("eosio")), name.AccountName(name.StringToName("eosio")), name.ActionName(name.StringToName("unlinkauth")), applyEosioUnlinkAuth)

	return controller
}

func (c *Controller) InitializeBlockchainState(genesis *GenesisState) error {
	log.Info("initializing new blockchain with genesis state")

	return nil
}

func (c *Controller) InitGenesis(session *state.Session, genesisConfig *GenesisState) error {
	// Validate the genesis configuration
	if err := genesisConfig.InitialConfiguration.Validate(); err != nil {
		return err
	}

	gpo := global.GlobalPropertyObject{
		Configuration:     genesisConfig.InitialConfiguration,
		WasmConfiguration: config.DefaultInitialWasmConfiguration(),
		ChainId:           c.ChainId,
	}

	if err := session.CreateGlobalPropertyObject(&gpo); err != nil {
		return err
	}

	systemAuthority := authority.Authority{
		Threshold: 1,
		Keys: []authority.KeyWeight{{
			Key:    genesisConfig.InitialKey,
			Weight: 1,
		}},
	}

	// Create eosio account
	if err := c.CreateNativeAccount(session, genesisConfig.InitialTimeStamp, config.SystemAccountName, systemAuthority, systemAuthority, true); err != nil {
		return fmt.Errorf("could not initialize account %s: %v", config.SystemAccountName, err)
	}

	emptyAuthority := authority.Authority{
		Threshold: 1,
		Keys:      make([]authority.KeyWeight, 0),
		Accounts:  make([]authority.PermissionLevelWeight, 0),
	}
	activeProducersAuthority := authority.Authority{
		Threshold: 1,
		Keys:      make([]authority.KeyWeight, 0),
		Accounts:  make([]authority.PermissionLevelWeight, 0),
	}
	activeProducersAuthority.Accounts = append(activeProducersAuthority.Accounts, authority.PermissionLevelWeight{
		Permission: authority.PermissionLevel{
			Actor:      config.SystemAccountName,
			Permission: config.ActiveName,
		},
		Weight: 1,
	})

	// Create eosio.null account
	if err := c.CreateNativeAccount(session, genesisConfig.InitialTimeStamp, config.NullAccountName, emptyAuthority, emptyAuthority, false); err != nil {
		return fmt.Errorf("could not initialize account %s: %v", config.NullAccountName, err)
	}

	// Create eosio.prods account
	if err := c.CreateNativeAccount(session, genesisConfig.InitialTimeStamp, config.ProducersAccountName, emptyAuthority, activeProducersAuthority, false); err != nil {
		return fmt.Errorf("could not initialize account %s: %v", config.NullAccountName, err)
	}

	authorization := c.GetAuthorizationManager(session)
	activePermission, err := authorization.GetPermission(authority.PermissionLevel{Actor: config.ProducersAccountName, Permission: config.ActiveName})

	if err != nil {
		return err
	}

	majorityPermission, err := authorization.CreatePermission(config.ProducersAccountName, config.MajorityProducersPermissionName, activePermission.ID, activeProducersAuthority, genesisConfig.InitialTimeStamp)

	if err != nil {
		return err
	}

	_, err = authorization.CreatePermission(config.ProducersAccountName, config.MinorityProducersPermissionName, majorityPermission.ID, activeProducersAuthority, genesisConfig.InitialTimeStamp)

	return err
}

func (c *Controller) GetAuthorizationManager(s *state.Session) *AuthorizationManager {
	return NewAuthorizationManager(c, s)
}

func (c *Controller) GetResourceLimitsManager(s *state.Session) *ResourceLimitsManager {
	return NewResourceLimitsManager(s)
}

func (c *Controller) CreateNativeAccount(session *state.Session, initialTimestamp time.TimePoint, name name.AccountName, owner authority.Authority, active authority.Authority, privileged bool) error {
	// Don't create account if it exists
	if existingAcc, _ := session.FindAccountByName(name); existingAcc != nil {
		return nil
	}

	acc := account.Account{
		Name:         name,
		CreationDate: block.NewBlockTimeStampFromTimePoint(initialTimestamp),
	}

	if name == config.SystemAccountName {
		abi := EosioContractAbi()
		abiBytes, err := abi.Encode()

		if err != nil {
			return err
		}

		acc.Abi = abiBytes
	}

	if err := session.CreateAccount(&acc); err != nil {
		return err
	}

	accountMetaData := account.AccountMetaDataObject{
		Name: name,
	}
	accountMetaData.SetPrivileged(privileged)
	if err := session.CreateAccountMetaData(&accountMetaData); err != nil {
		return err
	}

	authorizationManager := c.GetAuthorizationManager(session)
	ownerPermission, err := authorizationManager.CreatePermission(name, config.OwnerName, 0, owner, initialTimestamp)

	if err != nil {
		return err
	}

	activePermission, err := authorizationManager.CreatePermission(name, config.ActiveName, ownerPermission.ID, active, initialTimestamp)

	if err != nil {
		return err
	}

	resourceLimitsManager := c.GetResourceLimitsManager(session)

	if err := resourceLimitsManager.InitializeAccount(name); err != nil {
		return err
	}

	var ramDelta int64 = int64(config.OverheadPerAccountRamBytes)
	ramDelta += 2 * int64(authority.PermissionObjectBillableSize)
	ramDelta += int64(ownerPermission.Auth.GetBillableSize())
	ramDelta += int64(activePermission.Auth.GetBillableSize())

	if err := resourceLimitsManager.AddPendingRamUsage(name, ramDelta); err != nil {
		return err
	}

	if err := resourceLimitsManager.VerifyAccountRamUsage(name); err != nil {
		return err
	}

	return nil
}

func (c *Controller) SetApplyHandler(receiver name.AccountName, contract name.AccountName, action name.ActionName, handler func(a ApplyContext) error) {
	handlerKey := receiver + contract + action
	c.ApplyHandlers[handlerKey.String()] = handler
}

func (c *Controller) FindApplyHandler(receiver name.AccountName, scope name.AccountName, act name.ActionName) func(ApplyContext) error {
	handlerKey := receiver + scope + act

	if handler, ok := c.ApplyHandlers[handlerKey.String()]; ok {
		return handler
	}

	return nil
}

func (c *Controller) PushTransaction(trx transaction.TransactionMetaData, block *state.Block, session *state.Session) (*transaction.TransactionTrace, error) {
	c.transactionMutex.Lock()
	defer c.transactionMutex.Unlock()
	//start := core.Now()
	checkAuth := !trx.Implicit()
	signedTransaction, err := trx.PackedTrx().GetSignedTransaction()

	if err != nil {
		return nil, err
	}

	trxContext := NewTransactionContext(c, session, trx.PackedTrx(), *trx.Id(), block)

	if trx.Implicit() {
		if err := trxContext.InitForImplicitTransaction(0); err != nil {
			return nil, err
		}
	} else {
		if err := trxContext.InitForInputTransaction(uint64(trx.PackedTrx().GetUnprunableSize()), uint64(trx.PackedTrx().GetPrunableSize())); err != nil {
			return nil, err
		}
	}

	// Check authority
	if checkAuth {
		authorizationManager := c.GetAuthorizationManager(session)

		if err := authorizationManager.CheckAuthorization(signedTransaction.Actions, trx.RecoveredKeys(), []authority.PermissionLevel{}, false, nil); err != nil {
			log.Error("failed to check transaction authorization", "error", err)
			return nil, err
		}
	}

	if err := trxContext.Exec(); err != nil {
		log.Error("failed to exec transaction context", "error", err)
		return nil, err
	}

	if err := trxContext.Finalize(); err != nil {
		log.Error("failed to finalize transaction context", "error", err)
		return nil, err
	}

	trxContext.Trace.Hash = *trx.Id()
	trxContext.Trace.Receipt = transaction.TransactionReceipt{
		TransactionReceiptHeader: transaction.TransactionReceiptHeader{
			Status:        transaction.TransactionStatusExecuted,
			CpuUsageUs:    uint32(trxContext.BilledCpuTimeUs),
			NetUsageWords: fc.UnsignedInt(0),
		},
		Transaction: *trx.PackedTrx(),
	}

	return trxContext.Trace, nil
}

func (c *Controller) CheckContractList(code name.AccountName) error {
	if c.ContractWhitelist.Size() > 0 {
		if !c.ContractWhitelist.Contains(code) {
			return fmt.Errorf("account %s is not on the contract whitelist", code)
		}
	} else if c.ContractBlacklist.Size() > 0 {
		if c.ContractBlacklist.Contains(code) {
			return fmt.Errorf("account %s is on the contract blacklist", code)
		}
	}

	return nil
}

func (c *Controller) PendingBlockTime() time.TimePoint {
	return time.Now()
}

func (c *Controller) GetChainId() types.ChainIdType {
	return c.ChainId
}

// Only eosio for now
func (c *Controller) GetActiveProducers() ([]name.Name, error) {
	return []name.Name{config.SystemAccountName}, nil
}

func (c *Controller) CalculateTransactionMerkle(trxs []transaction.TransactionReceipt) (*crypto.Sha256, error) {
	digests := make([]crypto.Sha256, 0)

	for _, trx := range trxs {
		if digest, err := trx.Digest(); err != nil {
			return nil, err
		} else {
			digests = append(digests, *digest)
		}
	}

	merkle := utils.Merkle(digests)

	return &merkle, nil
}
