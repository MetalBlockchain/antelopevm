package chain_api_plugin

import (
	"context"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/state"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/gin-gonic/gin"
)

func init() {
	service.RegisterHandler("/v1/chain/get_info", GetInfo)
}

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
		HeadBlockNum:             lastAcceptedBlock.Header.Index,
		HeadBlockId:              lastAcceptedBlock.ID().Hex(),
		HeadBlockTime:            lastAcceptedBlock.Header.Created.String(),
		HeadBlockProducer:        "eosio",
		ChainId:                  chainId,
		LastIrreversibleBlockNum: lastAcceptedBlock.Header.Index,
		LastIrreversibleBlockId:  lastAcceptedBlock.ID().Hex(),
		ForkDbBlockNum:           lastAcceptedBlock.Header.Index,
		ForkDbBlockId:            lastAcceptedBlock.ID().Hex(),
		VirtualBlockCpuLimit:     200000000,
		VirtualBlockNetLimit:     1048576000,
		BlockCpuLimit:            198868,
		BlockNetLimit:            1048208,
	}
}

func GetInfo(vm service.VM) gin.HandlerFunc {
	return func(c *gin.Context) {
		version := "aaa"
		lastAcceptedId, _ := vm.LastAccepted(context.Background())
		lastAccepted, _ := vm.GetState().CreateSession(false).FindBlockByHash(core.BlockHash(lastAcceptedId))
		info := NewChainInfoResponse(version, lastAccepted, *crypto.Hash256("lll"))

		c.JSON(200, info)
	}
}
