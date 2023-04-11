package chain

import (
	"errors"
	"fmt"

	"github.com/MetalBlockchain/antelopevm/config"
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/authority"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/state"
	wasmApi "github.com/MetalBlockchain/antelopevm/wasm/api"
	"github.com/dgraph-io/badger/v3"
)

var (
	_                      wasmApi.AuthorizationManager = &AuthorizationManager{}
	errIrrelevantAuthority                              = errors.New("action declares irrelevant authority")
)

type AuthorizationManager struct {
	Controller *Controller
	Session    *state.Session
}

func NewAuthorizationManager(controller *Controller, s *state.Session) *AuthorizationManager {
	return &AuthorizationManager{
		Controller: controller,
		Session:    s,
	}
}

func (a *AuthorizationManager) GetPermission(level authority.PermissionLevel) (*core.Permission, error) {
	if level.Actor.Empty() || level.Permission.Empty() {
		return nil, fmt.Errorf("invalid permission")
	}

	return a.Session.FindPermissionByOwner(level.Actor, level.Permission)
}

func (a *AuthorizationManager) ModifyPermission(permission *core.Permission, auth *authority.Authority) error {
	return a.Session.ModifyPermission(permission, func() {
		permission.Auth = *auth
		permission.LastUpdated = core.Now()
	})
}

func (a *AuthorizationManager) CreatePermission(account name.AccountName, name name.PermissionName, parent core.IdType, auth authority.Authority, initialCreationTime core.TimePoint) (*core.Permission, error) {
	permission := &core.Permission{
		Parent:      parent,
		Owner:       account,
		Name:        name,
		LastUpdated: initialCreationTime,
		Auth:        auth,
	}

	if err := a.Session.CreatePermission(permission); err != nil {
		return nil, err
	}

	return permission, nil
}

// This function determines the minimum required permission for a certain action
// If there is a linked permission we return that, if not we default to the active name
func (a *AuthorizationManager) GetMinimumPermission(authorizerAccount name.AccountName, scope name.AccountName, actName name.ActionName) (*name.PermissionName, error) {
	linkedPermission, err := a.Session.FindPermissionLinkByActionName(authorizerAccount, scope, actName)

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return &config.ActiveName, nil
		}

		return nil, err
	}

	return &linkedPermission.RequiredPermission, nil
}

func (a *AuthorizationManager) GetPermissionAuthority(level *authority.PermissionLevel) (*authority.Authority, error) {
	permission, err := a.Session.FindPermissionByOwner(level.Actor, level.Permission)

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, fmt.Errorf("could not find permission")
		}

		return nil, err
	}

	return &permission.Auth, nil
}

func (a *AuthorizationManager) CheckAuthorization(actions []*core.Action, keys ecc.PublicKeySet, providedPermissions []authority.PermissionLevel, allowUnusedKeys bool, satisfiedAuthorizations authority.PermissionLevelSet) error {
	authorityChecker := NewAuthorityChecker(a.GetPermissionAuthority, keys, providedPermissions, config.MaxAuthDepth)
	permissionsToSatisfy := authority.NewPermissionLevelSet(4)

	for _, action := range actions {
		specialCase := false

		if action.Account == config.SystemAccountName {
			specialCase = true

			if action.Name == name.StringToName("updateauth") {

			} else if action.Name == name.StringToName("deleteauth") {

			} else if action.Name == name.StringToName("linkauth") {

			} else if action.Name == name.StringToName("unlinkauth") {

			} else {
				specialCase = false
			}
		}

		for _, declaredAuth := range action.Authorization {
			if !specialCase {
				minimumPermissionName, err := a.GetMinimumPermission(declaredAuth.Actor, action.Account, action.Name)

				if err != nil {
					return err
				}

				minPermission, err := a.GetPermission(authority.PermissionLevel{Actor: declaredAuth.Actor, Permission: *minimumPermissionName})

				if err != nil {
					return err
				}

				providedPermission, err := a.GetPermission(declaredAuth)

				if err != nil {
					return err
				}

				if !providedPermission.Satisfies(*minPermission) {
					return fmt.Errorf("action declares irrelevant authority '%v'; minimum authority is %v", declaredAuth, authority.PermissionLevel{Actor: minPermission.Owner, Permission: minPermission.Name})
				}
			}

			if satisfiedAuthorizations == nil || !satisfiedAuthorizations.Contains(declaredAuth) {
				permissionsToSatisfy.Insert(declaredAuth)
			}
		}
	}

	// Now verify that all the declared authorizations are satisfied:
	for _, permission := range permissionsToSatisfy.Slice() {
		if !authorityChecker.SatisfiedPermissionLevel(permission, nil) {
			return fmt.Errorf("transaction declares authority %s, but does not have signatures for it", permission)
		}
	}

	if !allowUnusedKeys && !authorityChecker.AllKeysUsed() {
		return fmt.Errorf("transaction bears irrelevant signatures from these keys: %v", authorityChecker.GetUnusedKeys())
	}

	return nil
}

func (a *AuthorizationManager) CheckAuthorizationByPermissionLevel(account name.AccountName, permission name.PermissionName, keys ecc.PublicKeySet, providedPermissions []authority.PermissionLevel, allowUnusedKeys bool) error {
	authorityChecker := NewAuthorityChecker(a.GetPermissionAuthority, keys, providedPermissions, config.MaxAuthDepth)
	authority := authority.PermissionLevel{Actor: account, Permission: permission}

	if !authorityChecker.SatisfiedPermissionLevel(authority, nil) {
		return fmt.Errorf("permission %v was not satisfied, provided permissions %v, provided keys %v", authority, providedPermissions, keys)
	}

	if !allowUnusedKeys && !authorityChecker.AllKeysUsed() {
		return fmt.Errorf("irrelevant keys provided: %v", authorityChecker.GetUnusedKeys())
	}

	return nil
}

func (a *AuthorizationManager) GetRequiredKeys(transaction core.Transaction, keys ecc.PublicKeySet) ([]ecc.PublicKey, error) {
	checker := NewAuthorityChecker(a.GetPermissionAuthority, keys, []authority.PermissionLevel{}, 16)

	for _, act := range transaction.Actions {
		for _, declaredAuth := range act.Authorization {
			if !checker.SatisfiedPermissionLevel(declaredAuth, nil) {
				return nil, fmt.Errorf("transaction declares authority '%s', but does not have signatures for it", declaredAuth.Actor.String())
			}
		}
	}

	return checker.GetUsedKeys(), nil
}
