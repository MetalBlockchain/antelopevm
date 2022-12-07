package service

import "github.com/MetalBlockchain/antelopevm/state"

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
		Timestamp:      block.Created.String(),
		Producer:       "eosio",
		Confirmed:      1,
		Previous:       block.PreviousBlock.Hex(),
		ID:             block.ID().Hex(),
		BlockNum:       block.Index,
		RefBlockNum:    block.Index,
		RefBlockPrefix: block.Index,
	}
}
