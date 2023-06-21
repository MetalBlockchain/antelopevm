package chain

import (
	"fmt"
	"strings"

	"github.com/MetalBlockchain/antelopevm/abi"
	"github.com/MetalBlockchain/antelopevm/config"
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/account"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
)

var (
	errDecode = fmt.Errorf("could not decode request body")
)

func applyEosioNewaccount(context ApplyContext) error {
	create := &NewAccount{}
	if err := rlp.DecodeBytes(context.GetAction().Data, create); err != nil {
		return errDecode
	}

	if err := context.RequireAuthorization(create.Creator); err != nil {
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
	creator, err := context.GetSession().FindAccountByName(create.Creator)

	if err == nil && !creator.Privileged {
		if strings.HasPrefix(nameStr, "eosio.") {
			return fmt.Errorf("only privileged accounts can have names that start with 'eosio.'")
		}
	}

	_, err = context.GetSession().FindAccountByName(create.Name)
	if err == nil {
		return fmt.Errorf("cannot create account named %s, as that name is already taken", create.Name.String())
	}

	//blockTime := context.Control.PendingBlockTime()
	newAccountObject := account.Account{Name: create.Name, CreationDate: core.Now()}

	if err := context.GetSession().CreateAccount(&newAccountObject); err != nil {
		return err
	}

	if ownerPermission, err := context.GetAuthorizationManager().CreatePermission(create.Name, name.StringToName("owner"), 0, create.Owner, core.TimePoint(0)); err == nil {
		if _, err := context.GetAuthorizationManager().CreatePermission(create.Name, name.StringToName("active"), ownerPermission.ID, create.Active, core.TimePoint(0)); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func applyEosioSetCode(context ApplyContext) error {
	act := &SetCode{}
	if err := rlp.DecodeBytes(context.GetAction().Data, act); err != nil {
		return errDecode
	}

	if err := context.RequireAuthorization(act.Account); err != nil {
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

	account, err := context.GetSession().FindAccountByName(act.Account)

	if err != nil {
		return fmt.Errorf("could not find account %s", act.Account.String())
	}

	if account.CodeVersion == *codeId {
		return fmt.Errorf("contract is already running this version of code")
	}

	if err := context.GetSession().ModifyAccount(account, func() {
		account.LastCodeUpdate = core.Now()
		account.CodeVersion = *codeId

		if codeSize > 0 {
			account.Code = act.Code
		}
	}); err != nil {
		return err
	}

	return nil
}

func applyEosioSetAbi(context ApplyContext) error {
	act := &SetAbi{}
	if err := rlp.DecodeBytes(context.GetAction().Data, act); err != nil {
		return errDecode
	}

	if err := context.RequireAuthorization(act.Account); err != nil {
		return err
	}

	abiSize := len(act.Abi)
	var abiId *crypto.Sha256

	if abiSize > 0 {
		abiId = crypto.Hash256(act.Abi)
		// TODO: Validate WASM code
	}

	// Validate ABI
	_, err := abi.NewABI(act.Abi)

	if err != nil {
		return err
	}

	account, err := context.GetSession().FindAccountByName(act.Account)

	if err != nil {
		return err
	}

	if err := context.GetSession().ModifyAccount(account, func() {
		if abiSize > 0 {
			account.Abi = act.Abi
			account.AbiVersion = *abiId
			account.AbiSequence += 1
		}
	}); err != nil {
		return err
	}

	return nil
}

func applyEosioUpdateAuth(context ApplyContext) error {
	act := UpdateAuth{}
	if err := rlp.DecodeBytes(context.GetAction().Data, &act); err != nil {
		return errDecode
	}

	if err := context.RequireAuthorization(act.Account); err != nil {
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

	if act.Permission == config.ActiveName {
		if act.Parent != config.OwnerName {
			return fmt.Errorf("cannot change active authority's parent from owner, update.parent %s", act.Parent.String())
		}
	} else if act.Permission == config.OwnerName {
		if !act.Parent.IsEmpty() {
			return fmt.Errorf("cannot change owner authority's parent")
		}
	}

	if len(act.Auth.Waits) > 0 {
		return fmt.Errorf("delayed authority is currently not supported")
	}

	permission, err := context.GetSession().FindPermissionByOwner(act.Account, act.Permission)

	if err != nil {
		return err
	}

	parentId := core.IdType(0)

	if act.Permission != config.OwnerName {
		parent, err := context.GetSession().FindPermissionByOwner(act.Account, act.Parent)

		if err != nil {
			return err
		}

		parentId = parent.ID
	}

	if permission != nil {
		if parentId != permission.Parent {
			return fmt.Errorf("changing parent authority is not currently supported")
		}

		if err := context.GetAuthorizationManager().ModifyPermission(permission, &act.Auth); err != nil {
			return err
		}
	} else {
		_, err := context.GetAuthorizationManager().CreatePermission(act.Account, act.Permission, parentId, act.Auth, core.TimePoint(0))

		return err
	}

	return nil
}

func applyEosioLinkAuth(context ApplyContext) error {
	act := LinkAuth{}
	if err := rlp.DecodeBytes(context.GetAction().Data, &act); err != nil {
		return errDecode
	}

	if act.Requirement.IsEmpty() {
		return fmt.Errorf("required permission cannot be empty")
	}

	if err := context.RequireAuthorization(act.Account); err != nil {
		return err
	}

	if act.Requirement != name.PermissionName(config.EosioAnyName) {
		if _, err := context.GetSession().FindPermissionByOwner(act.Account, act.Requirement); err != nil {
			return fmt.Errorf("failed to retrieve permission %s", name.NameToString(uint64(act.Requirement)))
		}
	}

	permissionLink, err := context.GetSession().FindPermissionLinkByActionName(act.Account, act.Code, act.Type)

	if err == nil {
		if permissionLink.RequiredPermission == act.Requirement {
			return fmt.Errorf("attempting to update required authority, but new requirement is same as old")
		}

		if err := context.GetSession().ModifyPermissionLink(permissionLink, func() {
			permissionLink.RequiredPermission = act.Requirement
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

		if err := context.GetSession().CreatePermissionLink(permissionLink); err != nil {
			return err
		}
	}

	return nil
}

func applyEosioUnlinkAuth(context ApplyContext) error {
	act := UnLinkAuth{}
	if err := rlp.DecodeBytes(context.GetAction().Data, &act); err != nil {
		return errDecode
	}

	if err := context.RequireAuthorization(act.Account); err != nil {
		return err
	}

	permissionLink, err := context.GetSession().FindPermissionLinkByActionName(act.Account, act.Code, act.Type)

	if err != nil {
		return fmt.Errorf("attempting to unlink authority, but no link found")
	}

	if err := context.GetSession().RemovePermissionLink(permissionLink); err != nil {
		return err
	}

	return nil
}

func applyEosioDeleteAuth(context ApplyContext) error {
	act := DeleteAuth{}
	if err := rlp.DecodeBytes(context.GetAction().Data, &act); err != nil {
		return errDecode
	}

	if err := context.RequireAuthorization(act.Account); err != nil {
		return err
	}

	if act.Permission == config.OwnerName {
		return fmt.Errorf("cannot delete owner authority")
	} else if act.Permission == config.ActiveName {
		return fmt.Errorf("cannot delete active authority")
	}

	permission, err := context.GetSession().FindPermissionByOwner(act.Account, act.Permission)

	if err != nil {
		return fmt.Errorf("cannot remove non-existing permission")
	}

	// TODO: Fix
	/* iterator := context.GetSession().GetPermissionLinksByPermissionName(act.Account, act.Permission)
	defer iterator.Release()

	if iterator.Next() {
		return fmt.Errorf("cannot delete a linked authority, remove the links first")
	} */

	if err := context.GetSession().RemovePermission(permission); err != nil {
		return err
	}

	return nil
}
