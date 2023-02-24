package api

import (
	"github.com/MetalBlockchain/antelopevm/core"
)

func GetAccountFunctions(context Context) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["require_auth"] = requireAuth(context)
	functions["has_auth"] = hasAuth(context)
	functions["require_auth2"] = requireAuth2(context)
	functions["require_recipient"] = requireRecipient(context)
	functions["is_account"] = isAccount(context)

	return functions
}

func requireAuth(context Context) func(core.AccountName) {
	return func(account core.AccountName) {
		if err := context.GetApplyContext().RequireAuthorization(account); err != nil {
			panic("missing authority of " + account.String())
		}
	}
}

func requireAuth2(context Context) func(core.AccountName, core.PermissionName) {
	return func(account core.AccountName, permission core.PermissionName) {
		if err := context.GetApplyContext().RequireAuthorizationWithPermission(account, permission); err != nil {
			panic("missing authority of " + account.String() + "/" + permission.String())
		}
	}
}

func isAccount(context Context) func(core.AccountName) uint32 {
	return func(account core.AccountName) uint32 {
		if ok := context.GetApplyContext().IsAccount(account); ok {
			return 1
		}

		return 0
	}
}

func requireRecipient(context Context) func(core.AccountName) {
	return func(recipient core.AccountName) {
		if err := context.GetApplyContext().RequireRecipient(recipient); err != nil {
			panic(err)
		}
	}
}

func hasAuth(context Context) func(core.AccountName) uint32 {
	return func(account core.AccountName) uint32 {
		if ok := context.GetApplyContext().HasAuthorization(account); ok {
			return 1
		}

		return 0
	}
}
