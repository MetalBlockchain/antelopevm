package chain

import (
	"github.com/MetalBlockchain/antelopevm/chain/authority"
	"github.com/MetalBlockchain/antelopevm/chain/name"
)

type NewAccount struct {
	Creator name.AccountName    `json:"creator"`
	Name    name.AccountName    `json:"name"`
	Owner   authority.Authority `json:"owner"`
	Active  authority.Authority `json:"active"`
}

type SetCode struct {
	Account   name.AccountName `json:"account"`
	VmType    uint8            `json:"vmtype"`
	VmVersion uint8            `json:"vmversion"`
	Code      []byte           `json:"code"`
}

type SetAbi struct {
	Account name.AccountName `json:"account"`
	Abi     []byte           `json:"abi"`
}

type UpdateAuth struct {
	Account    name.AccountName    `json:"account"`
	Permission name.PermissionName `json:"permission"`
	Parent     name.PermissionName `json:"parent"`
	Auth       authority.Authority `json:"auth"`
}

type DeleteAuth struct {
	Account    name.AccountName    `json:""`
	Permission name.PermissionName `json:""`
}

type LinkAuth struct {
	Account     name.AccountName    `json:"account"`
	Code        name.AccountName    `json:"code"`
	Type        name.ActionName     `json:"type"`
	Requirement name.PermissionName `json:"requirement"`
}

type UnLinkAuth struct {
	Account name.AccountName `json:"account"`
	Code    name.AccountName `json:"code"`
	Type    name.ActionName  `json:"type"`
}
