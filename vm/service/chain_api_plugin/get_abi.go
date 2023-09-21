package chain_api_plugin

import (
	"encoding/json"

	"github.com/MetalBlockchain/antelopevm/abi"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/gin-gonic/gin"
)

type GetAbiRequest struct {
	AccountName string `json:"account_name"`
}

type GetAbiResponse struct {
	AccountName string          `json:"account_name"`
	Abi         abi.ContractAbi `json:"abi,omitempty"`
}

func init() {
	service.RegisterHandler("/v1/chain/get_abi", GetAbi)
}

func GetAbi(vm service.VM) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body GetAbiRequest
		json.NewDecoder(c.Request.Body).Decode(&body)
		session := vm.GetState().CreateSession(false)
		defer session.Discard()
		acc, err := session.FindAccountByName(name.StringToName(body.AccountName))

		if err != nil {
			c.JSON(404, nil)
			return
		}

		if len(acc.Abi) == 0 {
			response := GetAbiResponse{
				AccountName: acc.Name.String(),
			}

			c.JSON(200, response)
			return
		}

		abi, err := abi.NewABI(acc.Abi)

		if err != nil {
			c.JSON(404, nil)
			return
		}

		response := GetAbiResponse{
			AccountName: acc.Name.String(),
			Abi:         *abi,
		}

		c.JSON(200, response)
	}
}
