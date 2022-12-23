package core

type Permission struct {
	ID          IdType         `serialize:"true" json:"-" eos:"-"`
	Parent      IdType         `serialize:"true" json:"parent"`
	UsageId     IdType         `serialize:"true" json:"-" eos:"-"`
	Owner       AccountName    `serialize:"true" json:"-" eos:"-"`
	Name        PermissionName `serialize:"true" json:"perm_name"`
	LastUpdated TimePoint      `serialize:"true" json:"-" eos:"-"`
	Auth        Authority      `serialize:"true" json:"required_auth"`
}

func (p *Permission) Satisfies(other Permission) bool {
	if p.Owner != other.Owner {
		return false
	}

	if p.ID == other.ID || p.ID == other.Parent {
		return true
	}

	return false
}

type PermissionUsage struct {
	ID       IdType    `serialize:"true"`
	LastUsed TimePoint `serialize:"true"`
}

type PermissionLink struct {
	ID                 IdType         `serialize:"true"`
	Account            AccountName    `serialize:"true"`
	Code               AccountName    `serialize:"true"`
	MessageType        ActionName     `serialize:"true"`
	RequiredPermission PermissionName `serialize:"true"`
}

type PermissionLevel struct {
	Actor      AccountName    `serialize:"true" json:"actor"`
	Permission PermissionName `serialize:"true" json:"permission"`
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
