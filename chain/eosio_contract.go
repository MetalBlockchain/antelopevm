package chain

import (
	"fmt"
	"strings"

	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/MetalBlockchain/antelopevm/state"
	"github.com/MetalBlockchain/antelopevm/utils"
)

func applyEosioNewaccount(context *ApplyContext) error {
	create := &NewAccount{}
	if err := rlp.DecodeBytes(context.Act.Data, create); err != nil {
		return err
	}

	if err := context.RequireAuthorization(int64(create.Creator)); err != nil {
		return err
	}

	utils.Assert(create.Owner.IsValid(), "Invalid owner authority")
	utils.Assert(create.Active.IsValid(), "Invalid active authority")

	nameStr := create.Name.String()
	utils.Assert(!create.Name.IsEmpty(), "account name cannot be empty")
	utils.Assert(len(nameStr) <= 12, "account names can only be 12 chars long")

	// Check if the creator is privileged
	creator, err := context.Control.State.GetAccountByName(create.Creator)

	if err == nil && !creator.Privileged {
		utils.Assert(!strings.HasPrefix(nameStr, "eosio."), "only privileged accounts can have names that start with 'eosio.'")
	}

	_, err = context.Control.State.GetAccountByName(create.Name)
	utils.Assert(err != nil, "Cannot create account named %s, as that name is already taken", create.Name.String())

	//blockTime := context.Control.PendingBlockTime()
	newAccountObject := state.Account{Name: create.Name, CreationDate: types.BlockTimeStamp(types.Now())}
	context.Control.State.PutAccount(&newAccountObject)

	ownerPermission, _ := context.Control.Authorization.CreatePermission(create.Name, types.N("owner"), 0, create.Owner, types.TimePoint(0))
	context.Control.Authorization.CreatePermission(create.Name, types.N("active"), ownerPermission.ID, create.Active, types.TimePoint(0))

	return nil
}

func applyEosioSetCode(context *ApplyContext) error {
	act := SetCode{}
	rlp.DecodeBytes(context.Act.Data, &act)

	if err := context.RequireAuthorization(int64(act.Account)); err != nil {
		return err
	}

	utils.Assert(act.VmType == 0, "code should be 0")
	utils.Assert(act.VmVersion == 0, "version should be 0")

	codeSize := len(act.Code)
	var codeId *crypto.Sha256

	if codeSize > 0 {
		codeId = crypto.Hash256(act.Code)
		// TODO: Validate WASM code
	}

	account, err := context.Control.State.GetAccountByName(act.Account)

	if err != nil {
		panic(err)
	}

	context.Control.State.UpdateAccount(account, func(a *state.Account) {
		a.LastCodeUpdate = types.TimePoint(types.Now())
		a.CodeVersion = *codeId

		if codeSize > 0 {
			a.Code = act.Code
		}
	})

	utils.Assert(account.CodeVersion != *codeId, "contract is already running this version of code")

	return nil
}

func applyEosioSetAbi(context *ApplyContext) error {
	act := SetAbi{}
	rlp.DecodeBytes(context.Act.Data, &act)

	if err := context.RequireAuthorization(int64(act.Account)); err != nil {
		return err
	}

	abiSize := len(act.Abi)
	account, err := context.Control.State.GetAccountByName(act.Account)

	if err != nil {
		return err
	}

	context.Control.State.UpdateAccount(account, func(a *state.Account) {
		if abiSize > 0 {
			account.Abi = act.Abi
			account.AbiSequence += 1
		}
	})

	return nil
}

func applyEosioUpdateAuth(context *ApplyContext) error {
	act := UpdateAuth{}
	rlp.DecodeBytes(context.Act.Data, &act)

	if err := context.RequireAuthorization(int64(act.Account)); err != nil {
		return err
	}

	if act.Permission.IsEmpty() {
		return fmt.Errorf("cannot create authority with empty name")
	}

	if strings.HasPrefix(act.Permission.String(), "eosio.") {
		return fmt.Errorf("permission names that start with 'eosio.' are reserved")
	}

	if act.Permission == act.Parent {
		return fmt.Errorf("cannot set an authority as its own parent")
	}

	if !act.Auth.IsValid() {
		return fmt.Errorf("invalid authority")
	}

	if act.Permission == context.Control.Config.ActiveName {
		if act.Parent != context.Control.Config.OwnerName {
			return fmt.Errorf("cannot change active authority's parent from owner, update.parent %s", act.Parent.String())
		}
	} else if act.Permission == context.Control.Config.OwnerName {
		if !act.Parent.IsEmpty() {
			return fmt.Errorf("cannot change owner authority's parent")
		}
	}

	if len(act.Auth.Waits) > 0 {
		return fmt.Errorf("delayed authority is currently not supported")
	}

	permission, _ := context.Control.State.GetPermissionByOwner(act.Account, act.Permission)
	parentId := types.IdType(0)

	if act.Permission != context.Control.Config.OwnerName {
		parent, _ := context.Control.State.GetPermissionByOwner(act.Account, act.Parent)
		parentId = parent.ID
	}

	if permission != nil {
		if parentId != permission.Parent {
			return fmt.Errorf("changing parent authority is not currently supported")
		}

		if err := context.Control.Authorization.ModifyPermission(permission, &act.Auth); err != nil {
			return err
		}
	} else {
		_, err := context.Control.Authorization.CreatePermission(act.Account, act.Permission, parentId, act.Auth, types.TimePoint(0))

		return err
	}

	return nil
}

func applyEosioLinkAuth(context *ApplyContext) error {
	act := LinkAuth{}
	rlp.DecodeBytes(context.Act.Data, &act)

	if act.Requirement.IsEmpty() {
		return fmt.Errorf("required permission cannot be empty")
	}

	if err := context.RequireAuthorization(int64(act.Account)); err != nil {
		return err
	}

	if act.Requirement != types.PermissionName(context.Control.Config.EosioAnyName) {
		if _, err := context.Control.State.GetPermissionByOwner(act.Account, act.Requirement); err != nil {
			return fmt.Errorf("failed to retrieve permission %s", types.S(uint64(act.Requirement)))
		}
	}

	permissionLink, err := context.Control.State.GetPermissionLinkByActionName(act.Account, act.Code, act.Type)

	if err == nil {
		if permissionLink.RequiredPermission == act.Requirement {
			return fmt.Errorf("attempting to update required authority, but new requirement is same as old")
		}

		if err := context.Control.State.UpdatePermissionLink(permissionLink, func(pl *state.PermissionLink) {
			pl.RequiredPermission = act.Requirement
		}); err != nil {
			return err
		}
	} else {
		permissionLink = &state.PermissionLink{
			Account:            act.Account,
			Code:               act.Code,
			MessageType:        act.Type,
			RequiredPermission: act.Requirement,
		}

		if err := context.Control.State.PutPermissionLink(permissionLink); err != nil {
			return err
		}
	}

	return nil
}

func applyEosioUnlinkAuth(context *ApplyContext) error {
	act := UnLinkAuth{}
	rlp.DecodeBytes(context.Act.Data, &act)

	if err := context.RequireAuthorization(int64(act.Account)); err != nil {
		return err
	}

	permissionLink, err := context.Control.State.GetPermissionLinkByActionName(act.Account, act.Code, act.Type)

	if err != nil {
		return fmt.Errorf("attempting to unlink authority, but no link found")
	}

	if err := context.Control.State.RemovePermissionLink(permissionLink); err != nil {
		return err
	}

	return nil
}

func applyEosioDeleteAuth(context *ApplyContext) error {
	act := DeleteAuth{}
	rlp.DecodeBytes(context.Act.Data, &act)

	if err := context.RequireAuthorization(int64(act.Account)); err != nil {
		return err
	}

	if act.Permission == context.Control.Config.OwnerName {
		return fmt.Errorf("cannot delete owner authority")
	} else if act.Permission == context.Control.Config.ActiveName {
		return fmt.Errorf("cannot delete active authority")
	}

	permission, err := context.Control.State.GetPermissionByOwner(act.Account, act.Permission)

	if err != nil {
		return fmt.Errorf("cannot remove non-existing permission")
	}

	iterator := context.Control.State.GetPermissionLinksByPermissionName(act.Account, act.Permission)
	defer iterator.Release()

	if iterator.Next() {
		return fmt.Errorf("cannot delete a linked authority, remove the links first")
	}

	if err := context.Control.State.RemovePermission(permission); err != nil {
		return err
	}

	return nil
}
