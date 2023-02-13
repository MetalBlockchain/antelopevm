package api

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/hashicorp/go-set"
	log "github.com/inconshreveable/log15"
)

func GetPermissionFunctions(context Context) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["check_transaction_authorization"] = checkTransactionAuthorization(context)
	functions["check_permission_authorization"] = checkPermissionAuthorization(context)
	functions["get_permission_last_used"] = getPermissionLastUsed(context)
	functions["get_account_creation_time"] = getAccountCreationTime(context)

	return functions
}

func checkTransactionAuthorization(context Context) func(uint32, uint32, uint32, uint32, uint32, uint32) int32 {
	return func(trxPtr uint32, trxSize uint32, pubKeys uint32, pubKeysSize uint32, perms uint32, permsSize uint32) int32 {
		log.Info("check_transaction_authorization", "trxPtr", trxPtr, "trxSize", trxSize, "pubKeys", pubKeys, "pubKeysSize", pubKeysSize, "perms", perms, "permsSize", permsSize)

		trxDataBytes := context.ReadMemory(trxPtr, trxSize)
		pubKeysBytes := context.ReadMemory(pubKeys, pubKeysSize)
		permsBytes := context.ReadMemory(perms, permsSize)

		// Parse trx data
		trx := core.Transaction{}
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

func checkPermissionAuthorization(context Context) func(core.AccountName, core.PermissionName, uint32, uint32, uint32, uint32, uint64) int32 {
	return func(account core.AccountName, permission core.PermissionName, pubKeys, pubKeysSize, perms, permsSize uint32, delay uint64) int32 {
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

func getPermissionLastUsed(context Context) func(core.AccountName, core.PermissionName) int64 {
	return func(account core.AccountName, permissionName core.PermissionName) int64 {
		permission, err := context.GetAuthorizationManager().GetPermission(core.PermissionLevel{Actor: account, Permission: permissionName})

		if err != nil {
			panic(err)
		}

		return permission.LastUsed.TimeSinceEpoch().Count()
	}
}

func getAccountCreationTime(context Context) func(core.AccountName) int64 {
	return func(accountName core.AccountName) int64 {
		account, err := context.GetApplyContext().FindAccount(accountName)

		if err != nil {
			panic(err)
		}

		return account.CreationDate.TimeSinceEpoch().Count()
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

func buildPermissionLevelSet(data []byte) *set.HashSet[core.PermissionLevel, string] {
	v := []core.PermissionLevel{}
	if err := rlp.DecodeBytes(data, v); err != nil {
		panic("could not decode permission level array")
	}
	set := set.NewHashSet[core.PermissionLevel, string](len(v))

	for i := range v {
		set.Insert(v[i])
	}

	return set
}
