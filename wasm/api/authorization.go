package api

import (
	"github.com/MetalBlockchain/antelopevm/core/name"
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

func requireAuth(context Context) func(name.AccountName) {
	return func(account name.AccountName) {
		if err := context.GetApplyContext().RequireAuthorization(account); err != nil {
			panic("missing authority of " + account.String())
		}
	}
}

func requireAuth2(context Context) func(name.AccountName, name.PermissionName) {
	return func(account name.AccountName, permission name.PermissionName) {
		if err := context.GetApplyContext().RequireAuthorizationWithPermission(account, permission); err != nil {
			panic("missing authority of " + account.String() + "/" + permission.String())
		}
	}
}

func isAccount(context Context) func(name.AccountName) uint32 {
	return func(account name.AccountName) uint32 {
		if ok := context.GetApplyContext().IsAccount(account); ok {
			return 1
		}

		return 0
	}
}

func requireRecipient(context Context) func(name.AccountName) {
	return func(recipient name.AccountName) {
		if err := context.GetApplyContext().RequireRecipient(recipient); err != nil {
			panic(err)
		}
	}
}

func hasAuth(context Context) func(name.AccountName) uint32 {
	return func(account name.AccountName) uint32 {
		if ok := context.GetApplyContext().HasAuthorization(account); ok {
			return 1
		}

		return 0
	}
}
