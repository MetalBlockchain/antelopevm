package state

import (
	"github.com/MetalBlockchain/antelopevm/chain/types"
)

type Permission struct {
	ID          types.IdType         `serialize:"true"`
	Parent      types.IdType         `serialize:"true"`
	UsageId     types.IdType         `serialize:"true"`
	Owner       types.AccountName    `serialize:"true"`
	Name        types.PermissionName `serialize:"true"`
	LastUpdated types.TimePoint      `serialize:"true"`
	Auth        types.Authority      `serialize:"true"`
}

func (p *Permission) Satisfies(other Permission) bool {
	if p.Owner != other.Owner {
		return false
	}

	if p.ID == other.ID {
		return true
	}

	return false
}

type PermissionUsage struct {
	ID       types.IdType    `serialize:"true"`
	LastUsed types.TimePoint `serialize:"true"`
}

type PermissionLink struct {
	ID                 types.IdType         `serialize:"true"`
	Account            types.AccountName    `serialize:"true"`
	Code               types.AccountName    `serialize:"true"`
	MessageType        types.ActionName     `serialize:"true"`
	RequiredPermission types.PermissionName `serialize:"true"`
}
