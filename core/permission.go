package core

import (
	"github.com/MetalBlockchain/antelopevm/core/authority"
	"github.com/MetalBlockchain/antelopevm/core/name"
)

var _ Entity = &Permission{}
var _ Entity = &PermissionLink{}

//go:generate msgp
type Permission struct {
	ID          IdType              `serialize:"true" json:"-" eos:"-"`
	Parent      IdType              `serialize:"true" json:"parent"`
	Owner       name.AccountName    `serialize:"true" json:"-" eos:"-"`
	Name        name.PermissionName `serialize:"true" json:"perm_name"`
	LastUpdated TimePoint           `serialize:"true" json:"-" eos:"-"`
	LastUsed    TimePoint           `serialize:"true" json:"-" eos:"-"`
	Auth        authority.Authority `serialize:"true" json:"required_auth"`
}

func (p Permission) GetId() []byte {
	return p.ID.ToBytes()
}

func (p Permission) GetIndexes() map[string]EntityIndex {
	return map[string]EntityIndex{
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
	return PermissionType
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
	ID                 IdType              `serialize:"true"`
	Account            name.AccountName    `serialize:"true"`
	Code               name.AccountName    `serialize:"true"`
	MessageType        name.ActionName     `serialize:"true"`
	RequiredPermission name.PermissionName `serialize:"true"`
}

func (p PermissionLink) GetId() []byte {
	return p.ID.ToBytes()
}

func (p PermissionLink) GetIndexes() map[string]EntityIndex {
	return map[string]EntityIndex{
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
	return PermissionLinkType
}
