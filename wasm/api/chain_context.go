package api

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/account"
	"github.com/MetalBlockchain/antelopevm/core/authority"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
)

type Controller interface {
	GetActiveProducers() ([]name.Name, error)
}

type AuthorizationManager interface {
	GetPermission(level authority.PermissionLevel) (*core.Permission, error)
	CheckAuthorization(actions []*core.Action, keys ecc.PublicKeySet, providedPermissions []authority.PermissionLevel, allowUnusedKeys bool, satisfiedAuthorizations authority.PermissionLevelSet) error
	CheckAuthorizationByPermissionLevel(account name.AccountName, permission name.PermissionName, keys ecc.PublicKeySet, providedPermissions []authority.PermissionLevel, allowUnusedKeys bool) error
}

type ApplyContext interface {
	RequireAuthorization(account name.AccountName) error
	RequireRecipient(recipient name.AccountName) error
	RequireAuthorizationWithPermission(account name.AccountName, permission name.PermissionName) error
	HasAuthorization(account name.AccountName) bool
	FindAccount(account name.AccountName) (*account.Account, error)
	IsAccount(account name.AccountName) bool
	GetSender() (*name.ActionName, error)

	GetAction() core.Action
	GetReceiver() name.AccountName

	// Database functions
	FindI64(code name.AccountName, scope name.ScopeName, table name.TableName, primaryKey uint64) int
	StoreI64(code name.AccountName, scope name.ScopeName, table name.TableName, payer name.AccountName, id uint64, buffer []byte) (int, error)
	GetI64(iterator int, buffer []byte, bufferSize int) (int, error)
	UpdateI64(iterator int, payer name.AccountName, buffer []byte, bufferSize int) error
	RemoveI64(iterator int) error
	NextI64(iterator int, primaryKey *uint64) (int, error)
	PreviousI64(iterator int, primaryKey *uint64) (int, error)
	LowerboundI64(code name.AccountName, scope name.ScopeName, table name.TableName, id uint64) (int, error)
	UpperboundI64(code name.AccountName, scope name.ScopeName, table name.TableName, id uint64) (int, error)
	EndI64(code name.AccountName, scope name.ScopeName, table name.TableName) (int, error)

	// Console functions
	ConsoleAppend(value string)

	SetActionReturnValue(value []byte)
	GetPackedTransaction() *core.PackedTransaction

	// Transaction functions
	ExecuteInline(action core.Action) error

	IsContextPrivileged() bool
	IsPrivileged(name name.AccountName) (bool, error)
	SetPrivileged(name name.AccountName, privileged bool) error
}

type MultiIndex[S any] interface {
	Store(scope name.ScopeName, tableName name.TableName, payer name.AccountName, primaryKey uint64, secondaryKey S) (int, error)
	Remove(iterator int) error
	Update(iterator int, payer name.AccountName, secondaryKey S) error
	FindSecondary(code name.AccountName, scope name.ScopeName, tableName name.TableName, secondaryKey *S, primaryKey *uint64) int
	LowerboundSecondary(code name.AccountName, scope name.ScopeName, tableName name.TableName, secondaryKey *S, primaryKey *uint64) int
	UpperboundSecondary(code name.AccountName, scope name.ScopeName, tableName name.TableName, secondaryKey *S, primaryKey *uint64) int
	EndSecondary(code name.AccountName, scope name.ScopeName, tableName name.TableName) int
	NextSecondary(iterator int, primaryKey *uint64) (int, error)
	PreviousSecondary(iterator int, primaryKey *uint64) (int, error)
	FindPrimary(code name.AccountName, scope name.ScopeName, tableName name.TableName, secondaryKey *S, primaryKey uint64) int
	LowerboundPrimary(code name.AccountName, scope name.ScopeName, tableName name.TableName, primaryKey uint64) int
	UpperboundPrimary(code name.AccountName, scope name.ScopeName, tableName name.TableName, primaryKey uint64) int
	NextPrimary(iterator int, primaryKey *uint64) (int, error)
	PreviousPrimary(iterator int, primaryKey *uint64) (int, error)
}

type ResourceLimitsManager interface {
	GetAccountLimits(account name.AccountName, ramBytes *int64, netWeight *int64, cpuWeight *int64) error
	SetAccountLimits(account name.AccountName, ramBytes int64, netWeight int64, cpuWeight int64) (bool, error)
}

func eosAssert(condition bool, msg string) {
	if !condition {
		panic(msg)
	}
}
