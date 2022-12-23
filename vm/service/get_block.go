package service

import (
	"encoding/json"
	"strconv"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/state"
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
	return GetBlockResponse{
		Timestamp:    block.Created.String(),
		Producer:     "eosio",
		Confirmed:    1,
		Previous:     block.PreviousBlock.Hex(),
		ID:           block.ID().Hex(),
		BlockNum:     block.Index,
		Transactions: block.Transactions,
	}
}
