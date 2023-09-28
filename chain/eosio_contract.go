package chain

import (
	"fmt"
	"strings"

	"github.com/MetalBlockchain/antelopevm/chain/account"
	"github.com/MetalBlockchain/antelopevm/chain/authority"
	"github.com/MetalBlockchain/antelopevm/chain/block"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/time"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/config"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/dgraph-io/badger/v3"
)

var (
	errDecode = fmt.Errorf("could not decode request body")
)

func validateAuthorityPrecondition(context ApplyContext, auth authority.Authority) error {
	for _, a := range auth.Accounts {
		if acct, _ := context.GetSession().FindAccountByName(a.Permission.Actor); acct == nil {
			return fmt.Errorf("account '%s' does not exist", a.Permission.Actor)
		}

		if a.Permission.Permission == config.OwnerName || a.Permission.Permission == config.ActiveName {
			continue // account was already checked to exist, so its owner and active permissions should exist
		}

		if a.Permission.Permission == config.EosioCodeName {
			continue // virtual eosio.code permission does not really exist but is allowed
		}

		if _, err := context.GetAuthorizationManager().GetPermission(authority.PermissionLevel{Actor: a.Permission.Actor, Permission: a.Permission.Permission}); err != nil {
			return fmt.Errorf("permission '%s' does not exist", a.Permission)
		}
	}

	return nil
}

func applyEosioNewaccount(context ApplyContext) error {
	create := &NewAccount{}
	if err := rlp.DecodeBytes(context.GetAction().Data, create); err != nil {
		return errDecode
	}

	if err := context.RequireAuthorization(create.Creator); err != nil {
		return err
	} else if !create.Owner.IsValid() {
		return fmt.Errorf("invalid owner authority")
	} else if !create.Active.IsValid() {
		return fmt.Errorf("invalid active authority")
	}

	nameStr := create.Name.String()
	if create.Name.IsEmpty() {
		return fmt.Errorf("account name cannot be empty")
	} else if len(nameStr) > 12 {
		return fmt.Errorf("account names can only be 12 chars long")
	}

	// Check if the creator is privileged
	if creator, err := context.GetSession().FindAccountMetaDataByName(create.Creator); err != nil {
		return err
	} else if strings.HasPrefix(nameStr, "eosio.") && !creator.IsPrivileged() {
		return fmt.Errorf("only privileged accounts can have names that start with 'eosio.'")
	} else if _, err := context.GetSession().FindAccountByName(create.Name); err == nil {
		return fmt.Errorf("cannot create account named %s, as that name is already taken", create.Name.String())
	}

	//blockTime := context.Control.PendingBlockTime()
	newAccountObject := account.Account{Name: create.Name, CreationDate: block.NewBlockTimeStampFromTimePoint(time.Now())}
	if err := context.GetSession().CreateAccount(&newAccountObject); err != nil {
		return err
	}

	newAccountMetaDataObject := account.AccountMetaDataObject{Name: create.Name}
	if err := context.GetSession().CreateAccountMetaData(&newAccountMetaDataObject); err != nil {
		return err
	}

	for _, auth := range []authority.Authority{create.Owner, create.Active} {
		if err := validateAuthorityPrecondition(context, auth); err != nil {
			return err
		}
	}

	ownerPermission, err := context.GetAuthorizationManager().CreatePermission(create.Name, name.StringToName("owner"), 0, create.Owner, time.TimePoint(0))
	if err != nil {
		return err
	}

	activePermission, err := context.GetAuthorizationManager().CreatePermission(create.Name, name.StringToName("active"), ownerPermission.ID, create.Active, time.TimePoint(0))
	if err != nil {
		return err
	}

	// TODO: initialize account for resources
	ramDelta := int64(config.OverheadPerAccountRamBytes)
	ramDelta += int64(2 * authority.PermissionObjectBillableSize)
	ramDelta += int64(ownerPermission.Auth.GetBillableSize())
	ramDelta += int64(activePermission.Auth.GetBillableSize())
	context.AddRamUsage(create.Name, ramDelta)

	return nil
}

func applyEosioSetCode(context ApplyContext) error {
	act := &SetCode{}
	if err := rlp.DecodeBytes(context.GetAction().Data, act); err != nil {
		return errDecode
	} else if err := context.RequireAuthorization(act.Account); err != nil {
		return err
	} else if act.VmType != 0 {
		return fmt.Errorf("code should be 0")
	} else if act.VmVersion != 0 {
		return fmt.Errorf("version should be 0")
	}

	codeHash := crypto.NewSha256Nil()
	codeSize := len(act.Code)

	if codeSize > 0 {
		codeHash = *crypto.Hash256(act.Code)
		// TODO: Validate WASM code
	}

	existingAccount, err := context.GetSession().FindAccountMetaDataByName(act.Account)
	if err != nil {
		return fmt.Errorf("could not find account %s", act.Account.String())
	}

	existingCode := (!existingAccount.CodeHash.IsEmpty() && !existingAccount.CodeHash.IsZero())

	if codeSize == 0 && !existingCode {
		return fmt.Errorf("contract is already cleared")
	}

	var oldSize int64
	var newSize int64 = int64(codeSize) * int64(config.SetCodeRamBytesMultiplier)

	if existingCode {
		if oldCodeEntry, err := context.GetSession().FindCodeObjectByCodeHash(existingAccount.CodeHash, existingAccount.VmType, existingAccount.VmVersion); err != nil {
			return err
		} else {
			if oldCodeEntry.CodeHash.Equals(codeHash) {
				return fmt.Errorf("contract is already running this version of code")
			}

			oldSize = (int64)(oldCodeEntry.Code.Size()) * int64(config.SetCodeRamBytesMultiplier)

			if oldCodeEntry.CodeRefCount == 1 {
				if err := context.GetSession().RemoveCodeObject(oldCodeEntry); err != nil {
					return err
				}
			} else {
				if err := context.GetSession().ModifyCodeObject(oldCodeEntry, func() {
					oldCodeEntry.CodeRefCount -= 1
				}); err != nil {
					return err
				}
			}
		}
	}

	if codeSize > 0 {
		newCodeEntry, err := context.GetSession().FindCodeObjectByCodeHash(codeHash, act.VmType, act.VmVersion)

		if err != nil {
			if err == badger.ErrKeyNotFound {
				newCodeEntry = &account.CodeObject{
					CodeHash:     codeHash,
					Code:         act.Code,
					CodeRefCount: 1,
					// TODO: Add FirstBlockUsed
					VmType:    act.VmType,
					VmVersion: act.VmVersion,
				}

				if err := context.GetSession().CreateCodeObject(newCodeEntry); err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			if err := context.GetSession().ModifyCodeObject(newCodeEntry, func() {
				newCodeEntry.CodeRefCount += 1
			}); err != nil {
				return err
			}
		}
	}

	if err := context.GetSession().ModifyAccountMetaData(existingAccount, func() {
		existingAccount.CodeSequence += 1
		existingAccount.CodeHash = codeHash
		existingAccount.VmType = act.VmType
		existingAccount.LastCodeUpdate = time.Now()
	}); err != nil {
		return err
	}

	if newSize != oldSize {
		context.AddRamUsage(act.Account, newSize-oldSize)
	}

	return nil
}

func applyEosioSetAbi(context ApplyContext) error {
	act := &SetAbi{}
	if err := rlp.DecodeBytes(context.GetAction().Data, act); err != nil {
		return errDecode
	} else if err := context.RequireAuthorization(act.Account); err != nil {
		return err
	}

	existingAccount, err := context.GetSession().FindAccountByName(act.Account)
	if err != nil {
		return fmt.Errorf("could not find account %s", act.Account.String())
	}
	abiSize := int64(len(act.Abi))
	oldSize := int64(existingAccount.Abi.Size())
	newSize := abiSize

	if err := context.GetSession().ModifyAccount(existingAccount, func() {
		existingAccount.Abi = act.Abi
	}); err != nil {
		return err
	}

	accountMetaData, err := context.GetSession().FindAccountMetaDataByName(act.Account)
	if err != nil {
		return err
	}

	if err := context.GetSession().ModifyAccountMetaData(accountMetaData, func() {
		accountMetaData.AbiSequence += 1
	}); err != nil {
		return err
	}

	if newSize != oldSize {
		context.AddRamUsage(act.Account, newSize-oldSize)
	}

	return nil
}

func applyEosioUpdateAuth(context ApplyContext) error {
	update := UpdateAuth{}
	if err := rlp.DecodeBytes(context.GetAction().Data, &update); err != nil {
		return errDecode
	}

	if err := context.RequireAuthorization(update.Account); err != nil {
		return err
	}

	if update.Permission.IsEmpty() {
		return fmt.Errorf("cannot create authority with empty name")
	}

	if strings.HasPrefix(update.Permission.String(), "eosio.") {
		return fmt.Errorf("permission names that start with 'eosio.' are reserved")
	}

	if update.Permission == update.Parent {
		return fmt.Errorf("cannot set an authority as its own parent")
	}

	if !update.Auth.IsValid() {
		return fmt.Errorf("invalid authority")
	}

	if update.Permission == config.ActiveName {
		if update.Parent != config.OwnerName {
			return fmt.Errorf("cannot change active authority's parent from owner, update.parent %s", update.Parent.String())
		}
	} else if update.Permission == config.OwnerName {
		if !update.Parent.IsEmpty() {
			return fmt.Errorf("cannot change owner authority's parent")
		}
	} else {
		if update.Parent.IsEmpty() {
			return fmt.Errorf("only owner permission can have empty parent")
		}
	}

	if len(update.Auth.Waits) > 0 {
		return fmt.Errorf("delayed authority is currently not supported")
	}

	if err := validateAuthorityPrecondition(context, update.Auth); err != nil {
		return err
	}

	permission, err := context.GetSession().FindPermissionByOwner(update.Account, update.Permission)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	parentId := types.IdType(0)

	if update.Permission != config.OwnerName {
		if parent, err := context.GetSession().FindPermissionByOwner(update.Account, update.Parent); err != nil {
			return err
		} else {
			parentId = parent.ID
		}
	}

	if permission != nil {
		if parentId != permission.Parent {
			return fmt.Errorf("changing parent authority is not currently supported")
		}

		oldSize := int64(authority.PermissionObjectBillableSize + permission.Auth.GetBillableSize())
		if err := context.GetAuthorizationManager().ModifyPermission(permission, &update.Auth); err != nil {
			return err
		}
		newSize := int64(authority.PermissionObjectBillableSize + permission.Auth.GetBillableSize())
		context.AddRamUsage(permission.Owner, newSize-oldSize)
	} else {
		if p, err := context.GetAuthorizationManager().CreatePermission(update.Account, update.Permission, parentId, update.Auth, time.TimePoint(0)); err != nil {
			return err
		} else {
			newSize := int64(authority.PermissionObjectBillableSize + p.Auth.GetBillableSize())
			context.AddRamUsage(permission.Owner, newSize)
		}
	}

	return nil
}

func applyEosioDeleteAuth(context ApplyContext) error {
	remove := DeleteAuth{}
	if err := rlp.DecodeBytes(context.GetAction().Data, &remove); err != nil {
		return errDecode
	} else if err := context.RequireAuthorization(remove.Account); err != nil {
		return err
	} else if remove.Permission == config.OwnerName {
		return fmt.Errorf("cannot delete owner authority")
	} else if remove.Permission == config.ActiveName {
		return fmt.Errorf("cannot delete active authority")
	}

	permission, err := context.GetSession().FindPermissionByOwner(remove.Account, remove.Permission)

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
		permissionLink = &authority.PermissionLink{
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
