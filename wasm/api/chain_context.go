package api

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
)

type AuthorizationManager interface {
	GetPermission(level core.PermissionLevel) (*core.Permission, error)
	CheckAuthorization(actions []*core.Action, keys ecc.PublicKeySet, providedPermissions []core.PermissionLevel, allowUnusedKeys bool, satisfiedAuthorizations core.PermissionLevelSet) error
	CheckAuthorizationByPermissionLevel(account core.AccountName, permission core.PermissionName, keys ecc.PublicKeySet, providedPermissions []core.PermissionLevel, allowUnusedKeys bool) error
}

type ApplyContext interface {
	RequireAuthorization(account core.AccountName) error
	RequireRecipient(recipient core.AccountName) error
	RequireAuthorizationWithPermission(account core.AccountName, permission core.PermissionName) error
	HasAuthorization(account core.AccountName) bool
	FindAccount(account core.AccountName) (*core.Account, error)
	IsAccount(account core.AccountName) bool
	GetSender() (*core.ActionName, error)

	GetAction() core.Action
	GetReceiver() core.AccountName

	// Database functions
	FindI64(code core.AccountName, scope core.ScopeName, table core.TableName, primaryKey uint64) int
	StoreI64(code core.AccountName, scope core.ScopeName, table core.TableName, payer core.AccountName, id uint64, buffer []byte) (int, error)
	GetI64(iterator int, buffer []byte, bufferSize int) (int, error)
	UpdateI64(iterator int, payer core.AccountName, buffer []byte, bufferSize int) error
	RemoveI64(iterator int) error
	NextI64(iterator int, primaryKey *uint64) (int, error)
	PreviousI64(iterator int, primaryKey *uint64) (int, error)
	LowerboundI64(code core.AccountName, scope core.ScopeName, table core.TableName, id uint64) (int, error)
	UpperboundI64(code core.AccountName, scope core.ScopeName, table core.TableName, id uint64) (int, error)
	EndI64(code core.AccountName, scope core.ScopeName, table core.TableName) (int, error)
}
