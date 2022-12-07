package chain

import (
	"encoding/json"

	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
)

type GenesisFile struct {
	InitialTimeStamp     types.TimePoint `json:"initial_timestamp"`
	InitialKey           ecc.PublicKey   `json:"initial_key"`
	InitialConfiguration GenesisConfig   `json:"initial_configuration"`
}

type GenesisConfig struct {
	MaxBlockNetUsage                    uint32 `json:"max_block_net_usage"`
	TargetBlockNetUsagePct              uint32 `json:"target_block_net_usage_pct"`
	MaxTransactionNetUsage              uint32 `json:"max_transaction_net_usage"`
	BasePerTransactionNetUsage          uint32 `json:"base_per_transaction_net_usage"`
	NetUsageLeeway                      uint32 `json:"net_usage_leeway"`
	ContextFreeDiscountNetUsageNum      uint32 `json:"context_free_discount_net_usage_num"`
	ContextFreeDiscountNetUsageDen      uint32 `json:"context_free_discount_net_usage_den"`
	MaxBlockCpuUsage                    uint32 `json:"max_block_cpu_usage"`
	TargetBlockCpuUsagePct              uint32 `json:"target_block_cpu_usage_pct"`
	MaxTransactionCpuUsage              uint32 `json:"max_transaction_cpu_usage"`
	MinTransactionCpuUsage              uint32 `json:"min_transaction_cpu_usage"`
	MaxTransactionLifetime              uint32 `json:"max_transaction_lifetime"`
	DeferredTransactionExpirationWindow uint32 `json:"deferred_trx_expiration_window"`
	MaxTransactionDelay                 uint32 `json:"max_transaction_delay"`
	MaxInlineActionSize                 uint32 `json:"max_inline_action_size"`
	MaxInlineActionDepth                uint32 `json:"max_inline_action_depth"`
	MaxAuthorityDepth                   uint32 `json:"max_authority_depth"`
	MaxRamSize                          uint32 `json:"max_ram_size"`
}

func (g *GenesisFile) GetChainId() types.ChainIdType {
	return types.ChainIdType(*crypto.Hash256(g))
}

func ParseGenesisData(data []byte) *GenesisFile {
	genesis := &GenesisFile{}
	json.Unmarshal(data, genesis)
	return genesis
}
