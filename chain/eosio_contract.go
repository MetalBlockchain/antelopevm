package chain

import (
	"fmt"
	"strings"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
)

var (
	errDecode = fmt.Errorf("could not decode request body")
)

func applyEosioNewaccount(context *ApplyContext) error {
	create := &NewAccount{}
	if err := rlp.DecodeBytes(context.Act.Data, create); err != nil {
		return errDecode
	}

	if err := context.RequireAuthorization(int64(create.Creator)); err != nil {
		return err
	}

	if !create.Owner.IsValid() {
		return fmt.Errorf("invalid owner authority")
	}

	if !create.Active.IsValid() {
		return fmt.Errorf("invalid active authority")
	}

	nameStr := create.Name.String()
	if create.Name.IsEmpty() {
		return fmt.Errorf("account name cannot be empty")
	}
	if len(nameStr) > 12 {
		return fmt.Errorf("account names can only be 12 chars long")
	}

	// Check if the creator is privileged
	creator, err := context.State.GetAccountByName(create.Creator)

	if err == nil && !creator.Privileged {
		if strings.HasPrefix(nameStr, "eosio.") {
			return fmt.Errorf("only privileged accounts can have names that start with 'eosio.'")
		}
	}

	_, err = context.State.GetAccountByName(create.Name)
	if err == nil {
		return fmt.Errorf("cannot create account named %s, as that name is already taken", create.Name.String())
	}

	//blockTime := context.Control.PendingBlockTime()
	newAccountObject := core.Account{Name: create.Name, CreationDate: core.BlockTimeStamp(core.Now())}
	context.State.PutAccount(&newAccountObject)

	ownerPermission, _ := context.Authorization.CreatePermission(create.Name, core.StringToName("owner"), 0, create.Owner, core.TimePoint(0))
	context.Authorization.CreatePermission(create.Name, core.StringToName("active"), ownerPermission.ID, create.Active, core.TimePoint(0))

	return nil
}

func applyEosioSetCode(context *ApplyContext) error {
	act := SetCode{}
	if err := rlp.DecodeBytes(context.Act.Data, &act); err != nil {
		return errDecode
	}

	if err := context.RequireAuthorization(int64(act.Account)); err != nil {
		return err
	}

	if act.VmType != 0 {
		return fmt.Errorf("code should be 0")
	}

	if act.VmVersion != 0 {
		return fmt.Errorf("version should be 0")
	}

	codeSize := len(act.Code)
	var codeId *crypto.Sha256

	if codeSize > 0 {
		codeId = crypto.Hash256(act.Code)
		// TODO: Validate WASM code
	}

	account, err := context.State.GetAccountByName(act.Account)

	if err != nil {
		return fmt.Errorf("could not find account %s", act.Account.String())
	}

	context.State.UpdateAccount(account, func(a *core.Account) {
		a.LastCodeUpdate = core.TimePoint(core.Now())
		a.CodeVersion = *codeId

		if codeSize > 0 {
			a.Code = act.Code
		}
	})

	if account.CodeVersion == *codeId {
		return fmt.Errorf("contract is already running this version of code")
	}

	return nil
}

func applyEosioSetAbi(context *ApplyContext) error {
	act := SetAbi{}
	if err := rlp.DecodeBytes(context.Act.Data, &act); err != nil {
		return errDecode
	}

	if err := context.RequireAuthorization(int64(act.Account)); err != nil {
		return err
	}

	abiSize := len(act.Abi)
	account, err := context.State.GetAccountByName(act.Account)

	if err != nil {
		return err
	}

	context.State.UpdateAccount(account, func(a *core.Account) {
		if abiSize > 0 {
			account.Abi = act.Abi
			account.AbiSequence += 1
		}
	})

	return nil
}

func applyEosioUpdateAuth(context *ApplyContext) error {
	act := UpdateAuth{}
	if err := rlp.DecodeBytes(context.Act.Data, &act); err != nil {
		return errDecode
	}

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

	permission, _ := context.State.GetPermissionByOwner(act.Account, act.Permission)
	parentId := core.IdType(0)

	if act.Permission != context.Control.Config.OwnerName {
		parent, _ := context.State.GetPermissionByOwner(act.Account, act.Parent)
		parentId = parent.ID
	}

	if permission != nil {
		if parentId != permission.Parent {
			return fmt.Errorf("changing parent authority is not currently supported")
		}

		if err := context.Authorization.ModifyPermission(permission, &act.Auth); err != nil {
			return err
		}
	} else {
		_, err := context.Authorization.CreatePermission(act.Account, act.Permission, parentId, act.Auth, core.TimePoint(0))

		return err
	}

	return nil
}

func applyEosioLinkAuth(context *ApplyContext) error {
	act := LinkAuth{}
	if err := rlp.DecodeBytes(context.Act.Data, &act); err != nil {
		return errDecode
	}

	if act.Requirement.IsEmpty() {
		return fmt.Errorf("required permission cannot be empty")
	}

	if err := context.RequireAuthorization(int64(act.Account)); err != nil {
		return err
	}

	if act.Requirement != core.PermissionName(context.Control.Config.EosioAnyName) {
		if _, err := context.State.GetPermissionByOwner(act.Account, act.Requirement); err != nil {
			return fmt.Errorf("failed to retrieve permission %s", core.NameToString(uint64(act.Requirement)))
		}
	}

	permissionLink, err := context.State.GetPermissionLinkByActionName(act.Account, act.Code, act.Type)

	if err == nil {
		if permissionLink.RequiredPermission == act.Requirement {
			return fmt.Errorf("attempting to update required authority, but new requirement is same as old")
		}

		if err := context.State.UpdatePermissionLink(permissionLink, func(pl *core.PermissionLink) {
			pl.RequiredPermission = act.Requirement
		}); err != nil {
			return err
		}
	} else {
		permissionLink = &core.PermissionLink{
			Account:            act.Account,
			Code:               act.Code,
			MessageType:        act.Type,
			RequiredPermission: act.Requirement,
		}

		if err := context.State.PutPermissionLink(permissionLink); err != nil {
			return err
		}
	}

	return nil
}

func applyEosioUnlinkAuth(context *ApplyContext) error {
	act := UnLinkAuth{}
	if err := rlp.DecodeBytes(context.Act.Data, &act); err != nil {
		return errDecode
	}

	if err := context.RequireAuthorization(int64(act.Account)); err != nil {
		return err
	}

	permissionLink, err := context.State.GetPermissionLinkByActionName(act.Account, act.Code, act.Type)

	if err != nil {
		return fmt.Errorf("attempting to unlink authority, but no link found")
	}

	if err := context.State.RemovePermissionLink(permissionLink); err != nil {
		return err
	}

	return nil
}

func applyEosioDeleteAuth(context *ApplyContext) error {
	act := DeleteAuth{}
	if err := rlp.DecodeBytes(context.Act.Data, &act); err != nil {
		return errDecode
	}

	if err := context.RequireAuthorization(int64(act.Account)); err != nil {
		return err
	}

	if act.Permission == context.Control.Config.OwnerName {
		return fmt.Errorf("cannot delete owner authority")
	} else if act.Permission == context.Control.Config.ActiveName {
		return fmt.Errorf("cannot delete active authority")
	}

	permission, err := context.State.GetPermissionByOwner(act.Account, act.Permission)

	if err != nil {
		return fmt.Errorf("cannot remove non-existing permission")
	}

	iterator := context.State.GetPermissionLinksByPermissionName(act.Account, act.Permission)
	defer iterator.Release()

	if iterator.Next() {
		return fmt.Errorf("cannot delete a linked authority, remove the links first")
	}

	if err := context.State.RemovePermission(permission); err != nil {
		return err
	}

	return nil
}
