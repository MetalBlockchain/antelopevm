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
	controller.SetApplyHandler(types.AccountName(types.N("eosio")), types.AccountName(types.N("eosio")), types.ActionName(types.N("deleteauth")), applyEosioDeleteAuth)
	controller.SetApplyHandler(types.AccountName(types.N("eosio")), types.AccountName(types.N("eosio")), types.ActionName(types.N("linkauth")), applyEosioLinkAuth)
	controller.SetApplyHandler(types.AccountName(types.N("eosio")), types.AccountName(types.N("eosio")), types.ActionName(types.N("unlinkauth")), applyEosioUnlinkAuth)

	return controller
}

func (c *Controller) InitGenesis(config *GenesisFile) error {
	if err := c.CreateNativeAccount(); err != nil {
		return err
	}

	auth := types.Authority{
		Threshold: 1,
		Keys: []types.KeyWeight{{
			Key:    config.InitialKey.String(),
			Weight: 1,
		}},
	}

	ownerPermission, err := c.Authorization.CreatePermission(c.Config.SystemAccountName, c.Config.OwnerName, types.IdType(0), auth, config.InitialTimeStamp)

	if err != nil {
		return err
	}

	if _, err := c.Authorization.CreatePermission(c.Config.SystemAccountName, c.Config.ActiveName, ownerPermission.ID, auth, config.InitialTimeStamp); err != nil {
		return err
	}

	return nil
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

func (c *Controller) PushTransaction(trx types.SignedTransaction) (*types.TransactionTrace, error) {
	// Check authority
	keys, err := trx.GetSignatureKeys(&c.ChainId, false, true)

	if err != nil {
		return nil, err
	}

	if err := c.Authorization.CheckAuthorization(keys, trx.Actions); err != nil {
		return nil, err
	}

	trxContext := NewTransactionContext(c, &trx, trx.ID())

	if err := trxContext.Init(); err != nil {
		return nil, err
	}

	if err := trxContext.Exec(); err != nil {
		return nil, err
	}

	if err := trxContext.Finalize(); err != nil {
		return nil, err
	}

	return trxContext.Trace, nil
}
