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
