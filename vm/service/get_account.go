package service

import (
	"github.com/MetalBlockchain/antelopevm/core/authority"
)

type GetAccountRequest struct {
	AccountName string `json:"account_name"`
}

type Limit struct {
	Available           uint64 `json:"available"`
	CurrentUsed         uint64 `json:"current_used"`
	LastUsageUpdateTime string `json:"last_usage_update_time"`
	Max                 uint64 `json:"max"`
	Used                uint64 `json:"used"`
}

type Resources struct {
	CpuWeight string `json:"cpu_weight"`
	NetWeight string `json:"net_weight"`
	Owner     string `json:"owner"`
	RamBytes  uint64 `json:"ram_bytes"`
}

type Permission struct {
	Parent string              `serialize:"true" json:"parent"`
	Name   string              `serialize:"true" json:"perm_name"`
	Auth   authority.Authority `serialize:"true" json:"required_auth"`
}

type GetAccountResponse struct {
	AccountName       string       `json:"account_name"`
	CpuLimit          Limit        `json:"cpu_limit"`
	CpuWeight         uint64       `json:"cpu_weight"`
	Created           string       `json:"created"`
	CoreLiquidBalance string       `json:"core_liquid_balance"`
	HeadBlockNum      uint64       `json:"head_block_num"`
	HeadBlockTime     string       `json:"head_block_time"`
	LastCodeUpdate    string       `json:"last_code_update"`
	NetLimit          Limit        `json:"net_limit"`
	NetWeight         uint64       `json:"net_weight"`
	Permissions       []Permission `json:"permissions"`
	Privileged        bool         `json:"privileged"`
	RamQuota          uint64       `json:"ram_quota"`
	RamUsage          uint64       `json:"ram_usage"`
	TotalResources    Resources    `json:"total_resources"`
}
