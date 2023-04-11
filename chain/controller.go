package chain

import (
	"fmt"
	"sync"

	"github.com/MetalBlockchain/antelopevm/config"
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/account"
	"github.com/MetalBlockchain/antelopevm/core/authority"
	"github.com/MetalBlockchain/antelopevm/core/global"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/state"
	"github.com/MetalBlockchain/antelopevm/wasm/api"
	log "github.com/inconshreveable/log15"
)

type v func(ctx ApplyContext) error

var _ api.Controller = &Controller{}

type Controller struct {
	ApplyHandlers map[string]v
	ChainId       core.ChainIdType
	State         *state.State

	contractWhitelist map[name.AccountName]bool
	contractBlacklist map[name.AccountName]bool

	transactionMutex sync.Mutex
}

func NewController(chainId core.ChainIdType, state *state.State) *Controller {
	controller := &Controller{
		ApplyHandlers: make(map[string]v),
		ChainId:       chainId,
		State:         state,

		contractWhitelist: make(map[name.AccountName]bool),
		contractBlacklist: make(map[name.AccountName]bool),
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

func (c *Controller) InitGenesis(session *state.Session, genesisConfig *GenesisFile) error {
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

func (c *Controller) CreateNativeAccount(session *state.Session, initialTimestamp core.TimePoint, name name.AccountName, owner authority.Authority, active authority.Authority, privileged bool) error {
	// Don't create account if it exists
	if existingAcc, _ := session.FindAccountByName(name); existingAcc != nil {
		return nil
	}

	account := account.Account{
		Name:         name,
		Privileged:   privileged,
		CreationDate: initialTimestamp,
	}

	if name == config.SystemAccountName {
		abi := EosioContractAbi()
		abiBytes, err := abi.Encode()

		if err != nil {
			return err
		}

		account.Abi = abiBytes
		account.AbiSequence = 1
	}

	if err := session.CreateAccount(&account); err != nil {
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
	ramDelta += 2 * int64(config.GetBillableSize("permission_object"))
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

func (c *Controller) PushTransaction(packedTrx core.PackedTransaction, block *state.Block, session *state.Session) (*core.TransactionTrace, error) {
	c.transactionMutex.Lock()
	defer c.transactionMutex.Unlock()
	trx, err := packedTrx.GetSignedTransaction()

	if err != nil {
		return nil, err
	}

	// Check authority
	keys, err := trx.GetSignatureKeys(&c.ChainId, false, true)

	if err != nil {
		return nil, err
	}

	authorizationManager := c.GetAuthorizationManager(session)

	if err := authorizationManager.CheckAuthorization(trx.Actions, keys, []authority.PermissionLevel{}, false, nil); err != nil {
		log.Error("failed to check transaction authorization", "error", err)
		return nil, err
	}

	trxContext := NewTransactionContext(c, session, &packedTrx, trx.ID(), block)

	if err := trxContext.Init(); err != nil {
		log.Error("failed to init transaction context", "error", err)
		return nil, err
	}

	if err := trxContext.Exec(); err != nil {
		log.Error("failed to exec transaction context", "error", err)
		return nil, err
	}

	if err := trxContext.Finalize(); err != nil {
		log.Error("failed to finalize transaction context", "error", err)
		return nil, err
	}

	trxContext.Trace.Hash = packedTrx.Id
	trxContext.Trace.Receipt = core.TransactionReceipt{
		TransactionReceiptHeader: core.TransactionReceiptHeader{
			Status:        core.TransactionStatusExecuted,
			CpuUsageUs:    uint32(trxContext.BilledCpuTimeUs),
			NetUsageWords: core.Vuint32(0),
		},
		Transaction: packedTrx,
	}

	return trxContext.Trace, nil
}

func (c *Controller) CheckContractList(code name.AccountName) error {
	if len(c.contractWhitelist) > 0 {
		if _, found := c.contractWhitelist[code]; !found {
			return fmt.Errorf("account %s is not on the contract whitelist", code)
		}
	} else if len(c.contractBlacklist) > 0 {
		if _, found := c.contractBlacklist[code]; found {
			return fmt.Errorf("account %s is on the contract blacklist", code)
		}
	}

	return nil
}

func (c *Controller) PendingBlockTime() core.TimePoint {
	return core.Now()
}

func (c *Controller) GetChainId() core.ChainIdType {
	return c.ChainId
}
