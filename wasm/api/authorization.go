package api

import (
	"github.com/MetalBlockchain/antelopevm/core/name"
)

func init() {
	Functions["require_auth"] = requireAuth
	Functions["has_auth"] = hasAuth
	Functions["require_auth2"] = requireAuth2
	Functions["require_recipient"] = requireRecipient
	Functions["is_account"] = isAccount
}

func GetAccountFunctions(context Context) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["require_auth"] = requireAuth(context)
	functions["has_auth"] = hasAuth(context)
	functions["require_auth2"] = requireAuth2(context)
	functions["require_recipient"] = requireRecipient(context)
	functions["is_account"] = isAccount(context)

	return functions
}

func requireAuth(context Context) interface{} {
	return func(account name.AccountName) {
		if err := context.GetApplyContext().RequireAuthorization(account); err != nil {
			panic("missing authority of " + account.String())
		}
	}
}

func requireAuth2(context Context) interface{} {
	return func(account name.AccountName, permission name.PermissionName) {
		if err := context.GetApplyContext().RequireAuthorizationWithPermission(account, permission); err != nil {
			panic("missing authority of " + account.String() + "/" + permission.String())
		}
	}
}

func isAccount(context Context) interface{} {
	return func(account name.AccountName) uint32 {
		if ok := context.GetApplyContext().IsAccount(account); ok {
			return 1
		}

		return 0
	}
}

func requireRecipient(context Context) interface{} {
	return func(recipient name.AccountName) {
		if err := context.GetApplyContext().RequireRecipient(recipient); err != nil {
			panic(err)
		}
	}
}

func hasAuth(context Context) interface{} {
	return func(account name.AccountName) uint32 {
		if ok := context.GetApplyContext().HasAuthorization(account); ok {
			return 1
		}

		return 0
	}
}
