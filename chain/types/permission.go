package types

type PermissionLevel struct {
	Actor      AccountName    `json:"actor"`
	Permission PermissionName `json:"permission"`
}

func ComparePermissionLevel(first interface{}, second interface{}) int {
	if first.(PermissionLevel).Actor > second.(PermissionLevel).Actor {
		return 1
	} else if first.(PermissionLevel).Actor < second.(PermissionLevel).Actor {
		return -1
	}
	if first.(PermissionLevel).Permission > second.(PermissionLevel).Permission {
		return 1
	} else if first.(PermissionLevel).Permission < second.(PermissionLevel).Permission {
		return -1
	} else {
		return 0
	}
}
