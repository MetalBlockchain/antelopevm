package chain_api_plugin

import (
	"encoding/hex"
	"encoding/json"
	"strconv"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/state"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/gin-gonic/gin"
)

// Some services send the ID as an int or a string, so we need to handle this
type BlockNumOrId string

func (fi *BlockNumOrId) UnmarshalJSON(b []byte) error {
	var value int
	if err := json.Unmarshal(b, &value); err != nil {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*fi = BlockNumOrId(s)
		return nil
	} else {
		*fi = BlockNumOrId(strconv.Itoa(value))
		return nil
	}
}

type GetBlockRequest struct {
	BlockNumOrId BlockNumOrId `json:"block_num_or_id"`
}

type GetBlockResponse struct {
	Timestamp    string                    `json:"timestamp"`
	Producer     string                    `json:"producer"`
	Confirmed    int                       `json:"confirmed"`
	Previous     string                    `json:"previous"`
	ID           string                    `json:"id"`
	BlockNum     uint64                    `json:"block_num"`
	Transactions []core.TransactionReceipt `json:"transactions"`
}

func NewGetBlockResponse(block *state.Block) GetBlockResponse {
	transactions := append([]core.TransactionReceipt{}, block.Transactions...)

	return GetBlockResponse{
		Timestamp:    block.Header.Created.String(),
		Producer:     "eosio",
		Confirmed:    1,
		Previous:     block.Header.PreviousBlockHash.Hex(),
		ID:           block.ID().Hex(),
		BlockNum:     block.Header.Index,
		Transactions: transactions,
	}
}

func init() {
	service.RegisterHandler("/v1/chain/get_block", GetBlock)
}

func GetBlock(vm service.VM) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body GetBlockRequest
		json.NewDecoder(c.Request.Body).Decode(&body)

		session := vm.GetState().CreateSession(false)
		defer session.Discard()

		if val, err := strconv.ParseUint(string(body.BlockNumOrId), 10, 64); err == nil {
			block, err := session.FindBlockByIndex(val)

			if err != nil {
				c.JSON(400, service.NewError(400, "could not parse block num"))
				return
			}

			c.JSON(200, NewGetBlockResponse(block))
			return
		}

		blockHash, err := hex.DecodeString(string(body.BlockNumOrId))

		if err != nil {
			c.JSON(400, service.NewError(400, "could not parse block num"))
			return
		}

		blockID, err := ids.ToID(blockHash)

		if err != nil {
			c.JSON(400, service.NewError(400, "could not parse block num"))
			return
		}

		block, err := session.FindBlockByHash(core.BlockHash(blockID))

		if err != nil {
			c.JSON(400, service.NewError(400, "could not parse block num"))
			return
		}

		c.JSON(200, NewGetBlockResponse(block))
	}
}
