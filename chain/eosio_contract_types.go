package chain

import "github.com/MetalBlockchain/antelopevm/core"

type NewAccount struct {
	Creator core.AccountName `json:"creator"`
	Name    core.AccountName `json:"name"`
	Owner   core.Authority   `json:"owner"`
	Active  core.Authority   `json:"active"`
}

type SetCode struct {
	Account   core.AccountName `json:"account"`
	VmType    uint8            `json:"vmtype"`
	VmVersion uint8            `json:"vmversion"`
	Code      []byte           `json:"code"`
}

type SetAbi struct {
	Account core.AccountName `json:"account"`
	Abi     []byte           `json:"abi"`
}

type UpdateAuth struct {
	Account    core.AccountName    `json:"account"`
	Permission core.PermissionName `json:"permission"`
	Parent     core.PermissionName `json:"parent"`
	Auth       core.Authority      `json:"auth"`
}

type DeleteAuth struct {
	Account    core.AccountName    `json:""`
	Permission core.PermissionName `json:""`
}

type LinkAuth struct {
	Account     core.AccountName    `json:"account"`
	Code        core.AccountName    `json:"code"`
	Type        core.ActionName     `json:"type"`
	Requirement core.PermissionName `json:"requirement"`
}

type UnLinkAuth struct {
	Account core.AccountName `json:"account"`
	Code    core.AccountName `json:"code"`
	Type    core.ActionName  `json:"type"`
}
