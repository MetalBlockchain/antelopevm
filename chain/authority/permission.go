package authority

import (
	"github.com/MetalBlockchain/antelopevm/chain/entity"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/resource"
	"github.com/MetalBlockchain/antelopevm/chain/time"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/config"
)

var _ entity.Entity = &Permission{}
var _ entity.Entity = &PermissionLink{}

var PermissionObjectBillableSize = resource.NewBillableSize((AuthorityBillableSize + 64) + uint64(5*config.OverheadPerRowPerIndexRamBytes))

type Permission struct {
	ID          types.IdType        `serialize:"true" json:"-" eos:"-"`
	Parent      types.IdType        `serialize:"true" json:"parent"`
	Owner       name.AccountName    `serialize:"true" json:"-" eos:"-"`
	Name        name.PermissionName `serialize:"true" json:"perm_name"`
	LastUpdated time.TimePoint      `serialize:"true" json:"-" eos:"-"`
	LastUsed    time.TimePoint      `serialize:"true" json:"-" eos:"-"`
	Auth        Authority           `serialize:"true" json:"required_auth"`
}

func (p Permission) GetId() []byte {
	return p.ID.ToBytes()
}

func (p Permission) GetIndexes() map[string]entity.EntityIndex {
	return map[string]entity.EntityIndex{
		"id": {
			Fields: []string{"ID"},
		},
		"byParent": {
			Fields: []string{"Parent", "ID"},
		},
		"byOwner": {
			Fields: []string{"Owner", "Name"},
		},
		"byName": {
			Fields: []string{"Name", "ID"},
		},
	}
}

func (p Permission) GetObjectType() uint8 {
	return entity.PermissionType
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
	ID       types.IdType   `serialize:"true"`
	LastUsed time.TimePoint `serialize:"true"`
}

type PermissionLink struct {
	ID                 types.IdType        `serialize:"true"`
	Account            name.AccountName    `serialize:"true"`
	Code               name.AccountName    `serialize:"true"`
	MessageType        name.ActionName     `serialize:"true"`
	RequiredPermission name.PermissionName `serialize:"true"`
}

func (p PermissionLink) GetId() []byte {
	return p.ID.ToBytes()
}

func (p PermissionLink) GetIndexes() map[string]entity.EntityIndex {
	return map[string]entity.EntityIndex{
		"id": {
			Fields: []string{"ID"},
		},
		"byActionName": {
			Fields: []string{"Account", "Code", "MessageType"},
		},
		"byPermissionName": {
			Fields: []string{"Account", "RequiredPermission", "ID"},
		},
	}
}

func (p PermissionLink) GetObjectType() uint8 {
	return entity.PermissionLinkType
}
