package chain

import (
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/state"
)

type v func(ctx *ApplyContext) error

type Controller struct {
	Config        *Config
	ApplyHandlers map[string]v
	State         state.State
	Authorization *AuthorizationManager
	ChainId       types.ChainIdType
}

func NewController(state state.State, chainId types.ChainIdType) *Controller {
	controller := &Controller{
		Config:        GetDefaultConfig(),
		State:         state,
		ApplyHandlers: make(map[string]v),
		ChainId:       chainId,
	}

	controller.Authorization = NewAuthorizationManager(controller)

	// Add native functions
	controller.CreateNativeAccount()
	controller.SetApplyHandler(types.AccountName(types.N("eosio")), types.AccountName(types.N("eosio")), types.ActionName(types.N("newaccount")), applyEosioNewaccount)
	controller.SetApplyHandler(types.AccountName(types.N("eosio")), types.AccountName(types.N("eosio")), types.ActionName(types.N("setcode")), applyEosioSetCode)
	controller.SetApplyHandler(types.AccountName(types.N("eosio")), types.AccountName(types.N("eosio")), types.ActionName(types.N("setabi")), applyEosioSetAbi)
	controller.SetApplyHandler(types.AccountName(types.N("eosio")), types.AccountName(types.N("eosio")), types.ActionName(types.N("updateauth")), applyEosioUpdateAuth)

	return controller
}

func (c *Controller) CreateNativeAccount() error {
	_, err := c.State.GetAccountByName(c.Config.SystemAccountName)

	if err != nil {
		nativeAccount := &state.Account{
			Name:       c.Config.SystemAccountName,
			Privileged: true,
		}

		return c.State.PutAccount(nativeAccount)
	}

	return nil
}

func (c *Controller) SetApplyHandler(receiver types.AccountName, contract types.AccountName, action types.ActionName, handler func(a *ApplyContext) error) {
	handlerKey := receiver + contract + action
	c.ApplyHandlers[handlerKey.String()] = handler
}

func (c *Controller) FindApplyHandler(receiver types.AccountName, scope types.AccountName, act types.ActionName) func(*ApplyContext) error {
	handlerKey := receiver + scope + act

	if handler, ok := c.ApplyHandlers[handlerKey.String()]; ok {
		return handler
	}

	return nil
}

func (c *Controller) PushTransaction(trx types.SignedTransaction) error {
	// Check authority
	keys, err := trx.GetSignatureKeys(&c.ChainId, false, true)

	if err != nil {
		return err
	}

	if err := c.Authorization.CheckAuthorization(keys, trx.Actions); err != nil {
		return err
	}

	trxContext := NewTransactionContext(c, &trx, trx.ID())
	trxContext.Init()
	trxContext.Exec()
	trxContext.Finalize()

	return nil
}
