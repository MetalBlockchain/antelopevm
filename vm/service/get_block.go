package service

import "github.com/MetalBlockchain/antelopevm/state"

var (
	blockTimeFormat = "2006-01-02T15:04:05.000"
)

type GetBlockRequest struct {
	BlockNumOrId string `json:"block_num_or_id"`
}

type GetBlockResponse struct {
	Timestamp string `json:"timestamp"`
	Producer  string `json:"producer"`
	Confirmed int    `json:"confirmed"`
	Previous  string `json:"previous"`
	ID        string `json:"id"`
	BlockNum  uint64 `json:"block_num"`
}

func NewGetBlockResponse(block *state.Block) GetBlockResponse {
	return GetBlockResponse{
		Timestamp: block.Created.String(),
		Producer:  "eosio",
		Confirmed: 1,
		Previous:  block.PreviousBlock.Hex(),
		ID:        block.ID().Hex(),
		BlockNum:  block.Index,
	}
}
