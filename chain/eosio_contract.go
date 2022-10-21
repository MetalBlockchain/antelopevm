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

	//oldSize := len(account.Code) * int(context.Control.Config.SetCodeRamBytesMultiplier)
	//newSize := codeSize * int(context.Control.Config.SetCodeRamBytesMultiplier)
	account.LastCodeUpdate = types.TimePoint(types.Now())
	account.CodeVersion = *codeId

	if codeSize > 0 {
		account.Code = act.Code
	}

	context.Control.State.UpdateAccount(account)

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

	if abiSize > 0 {
		account.Abi = act.Abi
		account.AbiSequence += 1
	}

	context.Control.State.UpdateAccount(account)

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
