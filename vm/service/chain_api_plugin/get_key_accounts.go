package chain_api_plugin

import (
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/MetalBlockchain/antelopevm/vm/service/history_api_plugin"
	"github.com/gin-gonic/gin"
)

func init() {
	service.RegisterHandler("/v1/chain/get_key_accounts", GetKeyAccounts)
}

func GetKeyAccounts(vm service.VM) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := history_api_plugin.GetKeyAccountsResponse{
			AccountNames: []name.AccountName{name.StringToName("joe"), name.StringToName("eosio")},
		}
		c.JSON(200, response)
	}
}
