package service

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/state"
)

type ChainInfoResponse struct {
	ServerVersion            string           `json:"server_version"`
	ServerVersionString      string           `json:"server_version_string"`
	ChainId                  core.ChainIdType `json:"chain_id"`
	BlockCpuLimit            int              `json:"block_cpu_limit"`
	BlockNetLimit            int              `json:"block_net_limit"`
	VirtualBlockCpuLimit     int              `json:"virtual_block_cpu_limit"`
	VirtualBlockNetLimit     int              `json:"virtual_block_net_limit"`
	HeadBlockNum             uint64           `json:"head_block_num"`
	HeadBlockId              string           `json:"head_block_id"`
	HeadBlockTime            string           `json:"head_block_time"`
	HeadBlockProducer        string           `json:"head_block_producer"`
	LastIrreversibleBlockNum uint64           `json:"last_irreversible_block_num"`
	LastIrreversibleBlockId  string           `json:"last_irreversible_block_id"`
	ForkDbBlockNum           uint64           `json:"fork_db_head_block_num"`
	ForkDbBlockId            string           `json:"fork_db_head_block_id"`
}

func NewChainInfoResponse(version string, lastAcceptedBlock *state.Block, chainId core.ChainIdType) *ChainInfoResponse {
	return &ChainInfoResponse{
		ServerVersion:            version,
		ServerVersionString:      version,
		HeadBlockNum:             lastAcceptedBlock.Index,
		HeadBlockId:              lastAcceptedBlock.ID().Hex(),
		HeadBlockTime:            lastAcceptedBlock.Created.String(),
		HeadBlockProducer:        "eosio",
		ChainId:                  chainId,
		LastIrreversibleBlockNum: lastAcceptedBlock.Index,
		LastIrreversibleBlockId:  lastAcceptedBlock.ID().Hex(),
		ForkDbBlockNum:           lastAcceptedBlock.Index,
		ForkDbBlockId:            lastAcceptedBlock.ID().Hex(),
		VirtualBlockCpuLimit:     200000000,
		VirtualBlockNetLimit:     1048576000,
		BlockCpuLimit:            198868,
		BlockNetLimit:            1048208,
	}
}
