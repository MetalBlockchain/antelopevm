package chain

import (
	"errors"

	"github.com/MetalBlockchain/antelopevm/chain/types"
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
}

func NewAuthorizationManager(controller *Controller) *AuthorizationManager {
	return &AuthorizationManager{
		Controller: controller,
	}
}

func (a *AuthorizationManager) ModifyPermission(permission *state.Permission, auth *types.Authority) error {
	permission.Auth = *auth
	permission.LastUpdated = types.Now()

	return a.Controller.State.UpdatePermission(permission)
}

func (a *AuthorizationManager) CreatePermission(account types.AccountName, name types.PermissionName, parent types.IdType, auth types.Authority, initialCreationTime types.TimePoint) (*state.Permission, error) {
	perm := &state.Permission{
		UsageId:     0,
		Parent:      parent,
		Owner:       account,
		Name:        name,
		LastUpdated: initialCreationTime,
		Auth:        auth,
	}

	err := a.Controller.State.PutPermission(perm)

	if err != nil {
		return nil, err
	}

	return perm, nil
}

// This function determines the minimum required permission for a certain action
// If there is a linked permission we return that, if not we default to the active name
func (a *AuthorizationManager) GetMinimumPermission(authorizerAccount types.AccountName, scope types.AccountName, actName types.ActionName) (*types.PermissionName, error) {
	linkedPermission, err := a.Controller.State.GetPermissionLinkByActionName(authorizerAccount, scope, actName)

	if err != nil {
		if err == database.ErrNotFound {
			return &a.Controller.Config.ActiveName, nil
		}
	}

	return &linkedPermission.RequiredPermission, nil
}

func (a *AuthorizationManager) CheckAuthorization(keys []ecc.PublicKey, actions []*types.Action) error {
	permissionsToSatisfy := make(map[types.PermissionLevel]state.Permission)

	for _, act := range actions {
		for _, declaredAuth := range act.Authorization {
			minPermissionName, err := a.GetMinimumPermission(declaredAuth.Actor, act.Account, act.Name)

			if err != nil {
				return err
			}

			minPermission, err := a.Controller.State.GetPermissionByOwner(declaredAuth.Actor, *minPermissionName)

			if err != nil {
				return err
			}

			declaredPermission, err := a.Controller.State.GetPermissionByOwner(declaredAuth.Actor, declaredAuth.Permission)

			if err != nil {
				return err
			}

			if !declaredPermission.Satisfies(*minPermission) {
				return errIrrelevantAuthority
			}

			permissionsToSatisfy[declaredAuth] = *declaredPermission
		}
	}

	authorityChecker := NewAuthorityChecker(keys)

	for _, permission := range permissionsToSatisfy {
		if !authorityChecker.Check(permission) {
			return errPermissionNotSatisfied
		}
	}

	return nil
}
