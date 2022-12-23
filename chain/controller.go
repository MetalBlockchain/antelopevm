package chain

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/state"
	log "github.com/inconshreveable/log15"
)

type v func(ctx *ApplyContext) error

type Controller struct {
	Config        *Config
	ApplyHandlers map[string]v
	ChainId       core.ChainIdType
}

func NewController(chainId core.ChainIdType) *Controller {
	controller := &Controller{
		Config:        GetDefaultConfig(),
		ApplyHandlers: make(map[string]v),
		ChainId:       chainId,
	}

	// Add native functions
	controller.SetApplyHandler(core.AccountName(core.StringToName("eosio")), core.AccountName(core.StringToName("eosio")), core.ActionName(core.StringToName("newaccount")), applyEosioNewaccount)
	controller.SetApplyHandler(core.AccountName(core.StringToName("eosio")), core.AccountName(core.StringToName("eosio")), core.ActionName(core.StringToName("setcode")), applyEosioSetCode)
	controller.SetApplyHandler(core.AccountName(core.StringToName("eosio")), core.AccountName(core.StringToName("eosio")), core.ActionName(core.StringToName("setabi")), applyEosioSetAbi)
	controller.SetApplyHandler(core.AccountName(core.StringToName("eosio")), core.AccountName(core.StringToName("eosio")), core.ActionName(core.StringToName("updateauth")), applyEosioUpdateAuth)
	controller.SetApplyHandler(core.AccountName(core.StringToName("eosio")), core.AccountName(core.StringToName("eosio")), core.ActionName(core.StringToName("deleteauth")), applyEosioDeleteAuth)
	controller.SetApplyHandler(core.AccountName(core.StringToName("eosio")), core.AccountName(core.StringToName("eosio")), core.ActionName(core.StringToName("linkauth")), applyEosioLinkAuth)
	controller.SetApplyHandler(core.AccountName(core.StringToName("eosio")), core.AccountName(core.StringToName("eosio")), core.ActionName(core.StringToName("unlinkauth")), applyEosioUnlinkAuth)

	return controller
}

func (c *Controller) InitGenesis(st state.State, config *GenesisFile) error {
	if err := c.CreateNativeAccount(st); err != nil {
		return err
	}

	auth := core.Authority{
		Threshold: 1,
		Keys: []core.KeyWeight{{
			Key:    config.InitialKey,
			Weight: 1,
		}},
	}

	authorizationManager := c.GetAuthorizationManager(st)
	ownerPermission, err := authorizationManager.CreatePermission(c.Config.SystemAccountName, c.Config.OwnerName, core.IdType(0), auth, config.InitialTimeStamp)

	if err != nil {
		return err
	}

	if _, err := authorizationManager.CreatePermission(c.Config.SystemAccountName, c.Config.ActiveName, ownerPermission.ID, auth, config.InitialTimeStamp); err != nil {
		return err
	}

	return nil
}

func (c *Controller) GetAuthorizationManager(s state.State) *AuthorizationManager {
	return NewAuthorizationManager(c, s)
}

func (c *Controller) CreateNativeAccount(st state.State) error {
	_, err := st.GetAccountByName(c.Config.SystemAccountName)

	if err != nil {
		nativeAccount := &core.Account{
			Name:       c.Config.SystemAccountName,
			Privileged: true,
		}

		return st.PutAccount(nativeAccount)
	}

	return nil
}

func (c *Controller) SetApplyHandler(receiver core.AccountName, contract core.AccountName, action core.ActionName, handler func(a *ApplyContext) error) {
	handlerKey := receiver + contract + action
	c.ApplyHandlers[handlerKey.String()] = handler
}

func (c *Controller) FindApplyHandler(receiver core.AccountName, scope core.AccountName, act core.ActionName) func(*ApplyContext) error {
	handlerKey := receiver + scope + act

	if handler, ok := c.ApplyHandlers[handlerKey.String()]; ok {
		return handler
	}

	return nil
}

func (c *Controller) PushTransaction(s state.State, packedTrx core.PackedTransaction) (*core.TransactionTrace, error) {
	trx, err := packedTrx.GetSignedTransaction()

	log.Info("got signed trx", "trx", trx)

	if err != nil {
		log.Error("failed to get signed transaction", "error", err)
		return nil, err
	}

	// Check authority
	keys, err := trx.GetSignatureKeys(&c.ChainId, false, true)

	if err != nil {
		log.Error("failed to get signature keys", "error", err)

		return nil, err
	}

	authorizationManager := c.GetAuthorizationManager(s)

	if err := authorizationManager.CheckAuthorization(keys, trx.Actions); err != nil {
		log.Error("failed to get check transaction authorization", "error", err)
		return nil, err
	}

	trxContext := NewTransactionContext(c, s, trx, trx.ID())

	if err := trxContext.Init(); err != nil {
		log.Error("failed to init transaction context", "error", err)
		return nil, err
	}

	if err := trxContext.Exec(); err != nil {
		log.Error("failed to exec transaction context", "error", err)
		return nil, err
	}

	trxContext.Trace.Id = packedTrx.Id
	trxContext.Trace.Receipt = core.TransactionReceipt{
		TransactionReceiptHeader: core.TransactionReceiptHeader{
			Status:        core.TransactionStatusExecuted,
			CpuUsageUs:    0,
			NetUsageWords: core.Vuint32(0),
		},
		Transaction: packedTrx,
	}

	if err := trxContext.Finalize(); err != nil {
		log.Error("failed to finalize transaction context", "error", err)
		return nil, err
	}

	return trxContext.Trace, nil
}
