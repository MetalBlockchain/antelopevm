package api

import (
	"github.com/MetalBlockchain/antelopevm/core"
	log "github.com/inconshreveable/log15"
)

func GetAccountFunctions(context Context) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["require_auth"] = requireAuth(context)
	functions["require_auth2"] = requireAuth2(context)
	functions["current_receiver"] = currentReceiver(context)
	functions["is_account"] = isAccount(context)
	functions["require_recipient"] = requireRecipient(context)
	functions["has_auth"] = hasAuth(context)

	return functions
}

func requireAuth(context Context) func(uint64) {
	return func(arg uint64) {
		account := core.Name(arg)
		log.Info("require_auth", "name", account.String())

		if err := context.GetApplyContext().RequireAuthorization(account); err != nil {
			panic("missing authority of " + account.String())
		}
	}
}

func requireAuth2(context Context) func(uint64, uint64) {
	return func(arg1 uint64, arg2 uint64) {
		account := core.Name(arg1)
		permission := core.Name(arg2)
		log.Info("require_auth2", "name", account.String(), "permission", permission.String())

		if err := context.GetApplyContext().RequireAuthorizationWithPermission(account, permission); err != nil {
			panic("missing authority of " + account.String() + "/" + permission.String())
		}
	}
}

func currentReceiver(context Context) func() uint64 {
	return func() uint64 {
		log.Info("current_receiver", "receiver", context.GetApplyContext().GetReceiver().String())
		return uint64(context.GetApplyContext().GetReceiver())
	}
}

func isAccount(context Context) func(uint64) uint32 {
	return func(arg uint64) uint32 {
		account := core.Name(arg)
		log.Info("is_account", "name", account.String())

		if ok := context.GetApplyContext().IsAccount(account); ok {
			return 1
		}

		return 0
	}
}

func requireRecipient(context Context) func(core.AccountName) {
	return func(recipient core.AccountName) {
		log.Info("require_receipient", "recipient", recipient)

		if err := context.GetApplyContext().RequireRecipient(recipient); err != nil {
			panic(err)
		}
	}
}

func hasAuth(context Context) func(uint64) uint32 {
	return func(arg uint64) uint32 {
		account := core.Name(arg)
		log.Info("has_auth", "name", account.String())

		if ok := context.GetApplyContext().HasAuthorization(account); ok {
			return 1
		}

		return 0
	}
}
