package chain_api_plugin

import (
	"encoding/json"

	"github.com/MetalBlockchain/antelopevm/state"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/gin-gonic/gin"
)

type GetBlockInfoRequest struct {
	BlockNum uint64 `json:"block_num"`
}

type GetBlockInfoResponse struct {
	Timestamp      string `json:"timestamp"`
	Producer       string `json:"producer"`
	Confirmed      int    `json:"confirmed"`
	Previous       string `json:"previous"`
	ID             string `json:"id"`
	BlockNum       uint64 `json:"block_num"`
	RefBlockNum    uint64 `json:"ref_block_num"`
	RefBlockPrefix uint64 `json:"ref_block_prefix"`
}

func NewGetBlockInfoResponse(block *state.Block) GetBlockInfoResponse {
	return GetBlockInfoResponse{
		Timestamp:      block.Header.Created.String(),
		Producer:       "eosio",
		Confirmed:      1,
		Previous:       block.Header.PreviousBlockHash.Hex(),
		ID:             block.Hash.Hex(),
		BlockNum:       block.Header.Index,
		RefBlockNum:    block.Header.Index,
		RefBlockPrefix: block.Header.Index,
	}
}

func init() {
	service.RegisterHandler("/v1/chain/get_block_info", GetBlockInfo)
}

func GetBlockInfo(vm service.VM) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body GetBlockInfoRequest
		json.NewDecoder(c.Request.Body).Decode(&body)

		session := vm.GetState().CreateSession(false)
		defer session.Discard()

		block, err := session.FindBlockByIndex(body.BlockNum)

		if err != nil {
			c.JSON(404, service.NewError(404, "block not found"))
			return
		}

		c.JSON(200, NewGetBlockInfoResponse(block))
	}
}
