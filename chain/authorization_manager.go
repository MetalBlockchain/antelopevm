package chain

import (
	"errors"
	"fmt"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/state"
	"github.com/MetalBlockchain/metalgo/database"
)

var (
	errIrrelevantAuthority    = errors.New("action declares irrelevant authority")
	errPermissionNotSatisfied = errors.New("declared permission wasn't satisfied")
)

type AuthorizationManager struct {
	Controller *Controller
	State      state.State
}

func NewAuthorizationManager(controller *Controller, s state.State) *AuthorizationManager {
	return &AuthorizationManager{
		Controller: controller,
		State:      s,
	}
}

func (a *AuthorizationManager) ModifyPermission(permission *core.Permission, auth *core.Authority) error {
	permission.Auth = *auth
	permission.LastUpdated = core.Now()

	return a.State.UpdatePermission(permission)
}

func (a *AuthorizationManager) CreatePermission(account core.AccountName, name core.PermissionName, parent core.IdType, auth core.Authority, initialCreationTime core.TimePoint) (*core.Permission, error) {
	perm := &core.Permission{
		UsageId:     0,
		Parent:      parent,
		Owner:       account,
		Name:        name,
		LastUpdated: initialCreationTime,
		Auth:        auth,
	}

	err := a.State.PutPermission(perm)

	if err != nil {
		return nil, err
	}

	return perm, nil
}

// This function determines the minimum required permission for a certain action
// If there is a linked permission we return that, if not we default to the active name
func (a *AuthorizationManager) GetMinimumPermission(authorizerAccount core.AccountName, scope core.AccountName, actName core.ActionName) (*core.PermissionName, error) {
	linkedPermission, err := a.State.GetPermissionLinkByActionName(authorizerAccount, scope, actName)

	if err != nil {
		if err == database.ErrNotFound {
			return &a.Controller.Config.ActiveName, nil
		}
	}

	return &linkedPermission.RequiredPermission, nil
}

func (a *AuthorizationManager) GetPermissionAuthority(level *core.PermissionLevel) (*core.Authority, error) {
	permission, err := a.State.GetPermissionByOwner(level.Actor, level.Permission)

	if err != nil {
		return nil, err
	}

	return &permission.Auth, nil
}

func (a *AuthorizationManager) CheckAuthorization(keys []ecc.PublicKey, actions []*core.Action) error {
	permissionsToSatisfy := make(map[core.PermissionLevel]core.Permission)

	for _, act := range actions {
		for _, declaredAuth := range act.Authorization {
			minPermissionName, err := a.GetMinimumPermission(declaredAuth.Actor, act.Account, act.Name)

			if err != nil {
				return err
			}

			minPermission, err := a.State.GetPermissionByOwner(declaredAuth.Actor, *minPermissionName)

			if err != nil {
				return err
			}

			declaredPermission, err := a.State.GetPermissionByOwner(declaredAuth.Actor, declaredAuth.Permission)

			if err != nil {
				return err
			}

			if !declaredPermission.Satisfies(*minPermission) {
				return errIrrelevantAuthority
			}

			permissionsToSatisfy[declaredAuth] = *declaredPermission
		}
	}

	//authorityChecker := NewAuthorityChecker(keys, nil, 16)

	/* for _, permission := range permissionsToSatisfy {
		if !authorityChecker.Check(permission) {
			return errPermissionNotSatisfied
		}
	} */

	return nil
}

func (a *AuthorizationManager) GetRequiredKeys(transaction core.Transaction, keys []ecc.PublicKey) ([]ecc.PublicKey, error) {
	checker := NewAuthorityChecker(a.GetPermissionAuthority, keys, []core.PermissionLevel{}, 16)

	for _, act := range transaction.Actions {
		for _, declaredAuth := range act.Authorization {
			if !checker.SatisfiedPermissionLevel(declaredAuth, nil) {
				return nil, fmt.Errorf("transaction declares authority '%s', but does not have signatures for it", declaredAuth.Actor.String())
			}
		}
	}

	return checker.GetUsedKeys(), nil
}
