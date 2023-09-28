package api

import (
	"github.com/MetalBlockchain/antelopevm/chain/authority"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/transaction"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/hashicorp/go-set"
)

func init() {
	Functions["check_transaction_authorization"] = checkTransactionAuthorization
	Functions["check_permission_authorization"] = checkPermissionAuthorization
	Functions["get_permission_last_used"] = getPermissionLastUsed
	Functions["get_account_creation_time"] = getAccountCreationTime
}

func checkTransactionAuthorization(context Context) interface{} {
	return func(trxPtr uint32, trxSize uint32, pubKeys uint32, pubKeysSize uint32, perms uint32, permsSize uint32) int32 {
		trxDataBytes := context.ReadMemory(trxPtr, trxSize)
		pubKeysBytes := context.ReadMemory(pubKeys, pubKeysSize)
		permsBytes := context.ReadMemory(perms, permsSize)

		// Parse trx data
		trx := transaction.Transaction{}
		if err := rlp.DecodeBytes(trxDataBytes, trx); err != nil {
			panic("could not decode transaction data")
		}

		// Parse public keys
		pubKeySet := buildPublicKeySet(pubKeysBytes)

		// Parse permission levels
		permLevelSet := buildPermissionLevelSet(permsBytes)

		if err := context.GetAuthorizationManager().CheckAuthorization(trx.Actions, pubKeySet, permLevelSet.Slice(), false, nil); err == nil {
			return 1
		}

		return 0
	}
}

func checkPermissionAuthorization(context Context) interface{} {
	return func(account name.AccountName, permission name.PermissionName, pubKeys, pubKeysSize, perms, permsSize uint32, delay uint64) int32 {
		pubKeysBytes := context.ReadMemory(pubKeys, pubKeysSize)
		permsBytes := context.ReadMemory(perms, permsSize)

		// Parse public keys
		pubKeySet := buildPublicKeySet(pubKeysBytes)

		// Parse permission levels
		permLevelSet := buildPermissionLevelSet(permsBytes)

		if err := context.GetAuthorizationManager().CheckAuthorizationByPermissionLevel(account, permission, pubKeySet, permLevelSet.Slice(), false); err == nil {
			return 1
		}

		return 0
	}
}

func getPermissionLastUsed(context Context) interface{} {
	return func(account name.AccountName, permissionName name.PermissionName) int64 {
		permission, err := context.GetAuthorizationManager().GetPermission(authority.PermissionLevel{Actor: account, Permission: permissionName})

		if err != nil {
			panic(err)
		}

		return permission.LastUsed.TimeSinceEpoch().Count()
	}
}

func getAccountCreationTime(context Context) interface{} {
	return func(accountName name.AccountName) int64 {
		account, err := context.GetApplyContext().FindAccount(accountName)

		if err != nil {
			panic(err)
		}

		return account.CreationDate.ToTimePoint().TimeSinceEpoch().Count()
	}
}

func buildPublicKeySet(data []byte) *set.HashSet[ecc.PublicKey, string] {
	v := []ecc.PublicKey{}
	if err := rlp.DecodeBytes(data, v); err != nil {
		panic("could not decode public key array")
	}
	set := set.NewHashSet[ecc.PublicKey, string](len(v))

	for i := range v {
		set.Insert(v[i])
	}

	return set
}

func buildPermissionLevelSet(data []byte) *set.HashSet[authority.PermissionLevel, string] {
	v := []authority.PermissionLevel{}
	if err := rlp.DecodeBytes(data, v); err != nil {
		panic("could not decode permission level array")
	}
	set := set.NewHashSet[authority.PermissionLevel, string](len(v))

	for i := range v {
		set.Insert(v[i])
	}

	return set
}
