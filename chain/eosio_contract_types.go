package chain

import (
	"github.com/MetalBlockchain/antelopevm/chain/types"
)

type NewAccount struct {
	Creator types.AccountName `json:"creator"`
	Name    types.AccountName `json:"name"`
	Owner   types.Authority   `json:"owner"`
	Active  types.Authority   `json:"active"`
}

type SetCode struct {
	Account   types.AccountName `json:"account"`
	VmType    uint8             `json:"vmtype"`
	VmVersion uint8             `json:"vmversion"`
	Code      []byte            `json:"code"`
}

type SetAbi struct {
	Account types.AccountName `json:"account"`
	Abi     []byte            `json:"abi"`
}

type UpdateAuth struct {
	Account    types.AccountName    `json:"account"`
	Permission types.PermissionName `json:"permission"`
	Parent     types.PermissionName `json:"parent"`
	Auth       types.Authority      `json:"auth"`
}

type DeleteAuth struct {
	Account    types.AccountName    `json:""`
	Permission types.PermissionName `json:""`
}

type LinkAuth struct {
	Account     types.AccountName    `json:"account"`
	Code        types.AccountName    `json:"code"`
	Type        types.ActionName     `json:"type"`
	Requirement types.PermissionName `json:"requirement"`
}

type UnLinkAuth struct {
	Account types.AccountName `json:"account"`
	Code    types.AccountName `json:"code"`
	Type    types.ActionName  `json:"type"`
}
