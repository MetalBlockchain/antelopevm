package chain_api_plugin

import (
	"encoding/json"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/authority"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/gin-gonic/gin"
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

func init() {
	service.RegisterHandler("/v1/chain/get_account", GetAccount)
}

func GetAccount(vm service.VM) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body GetAccountRequest
		json.NewDecoder(c.Request.Body).Decode(&body)
		session := vm.GetState().CreateSession(false)
		defer session.Discard()
		acc, err := session.FindAccountByName(name.StringToName(body.AccountName))

		if err != nil {
			c.JSON(404, service.NewError(404, "account not found"))
			return
		}

		permissions := make([]*core.Permission, 0)
		iterator := session.FindPermissionsByOwner(acc.Name)
		defer iterator.Close()

		for iterator.Rewind(); iterator.Valid(); iterator.Next() {
			if item, err := iterator.Item(); err == nil {
				permissions = append(permissions, item)
			} else {
				c.JSON(404, service.NewError(404, "account not found"))
				return
			}
		}

		if err != nil {
			c.JSON(404, service.NewError(404, "account not found"))
			return
		}

		response := GetAccountResponse{
			AccountName: body.AccountName,
			CpuLimit: Limit{
				Available:           16346040,
				CurrentUsed:         0,
				LastUsageUpdateTime: core.Now().String(),
				Max:                 16346421,
				Used:                381,
			},
			CpuWeight:         500000,
			Created:           acc.CreationDate.String(),
			CoreLiquidBalance: "1000.0000 SYS",
			HeadBlockNum:      0,
			HeadBlockTime:     core.Now().String(),
			LastCodeUpdate:    acc.LastCodeUpdate.String(),
			NetLimit: Limit{
				Available:           88094630,
				CurrentUsed:         0,
				LastUsageUpdateTime: core.Now().String(),
				Max:                 88094878,
				Used:                248,
			},
			NetWeight:   500000,
			Permissions: make([]Permission, 0),
			Privileged:  acc.Privileged,
			RamQuota:    525686,
			RamUsage:    5544,
			TotalResources: Resources{
				CpuWeight: "50.0000 SYS",
				NetWeight: "50.0000 SYS",
				Owner:     body.AccountName,
				RamBytes:  524286,
			},
		}

		for _, permission := range permissions {
			parent := ""

			for _, p := range permissions {
				if p.ID != permission.ID && permission.Parent == p.ID {
					parent = p.Name.String()
				}
			}

			response.Permissions = append(response.Permissions, Permission{
				Parent: parent,
				Name:   permission.Name.String(),
				Auth:   permission.Auth,
			})
		}

		c.JSON(200, response)
	}
}
