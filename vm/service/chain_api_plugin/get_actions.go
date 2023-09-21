package chain_api_plugin

import (
	"encoding/json"
	"net/http"

	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/gin-gonic/gin"
)

type GetActionsRequest struct {
	AccountName string `json:"account_name"`
}

type GetActionsResponse struct {
	Actions                  []string `json:"actions"`
	HeadBlockNum             int      `json:"head_block_num"`
	LastIrreversibleBlockNum int      `json:"last_irreversible_block"`
}

func init() {
	service.RegisterHandler("/v1/chain/get_actions", service.Handler{
		Methods:     []string{http.MethodPost},
		HandlerFunc: GetActions,
	})
}

func GetActions(vm service.VM) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body GetActionsRequest
		json.NewDecoder(c.Request.Body).Decode(&body)

		c.JSON(200, GetActionsResponse{
			Actions:                  []string{},
			HeadBlockNum:             0,
			LastIrreversibleBlockNum: 0,
		})
	}
}
