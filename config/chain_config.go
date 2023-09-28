package config

import (
	"errors"
	"fmt"
)

type ChainConfig struct {
	MaxBlockNetUsage               uint64 `json:"max_block_net_usage"`
	TargetBlockNetUsagePct         uint32 `json:"target_block_net_usage_pct"`
	MaxTransactionNetUsage         uint32 `json:"max_transaction_net_usage"`
	BasePerTransactionNetUsage     uint32 `json:"base_per_transaction_net_usage"`
	NetUsageLeeway                 uint32 `json:"net_usage_leeway"`
	ContextFreeDiscountNetUsageNum uint32 `json:"context_free_discount_net_usage_num"`
	ContextFreeDiscountNetUsageDen uint32 `json:"context_free_discount_net_usage_den"`

	MaxBlockCpuUsage       uint32 `json:"max_block_cpu_usage"`
	TargetBlockCpuUsagePct uint32 `json:"target_block_cpu_usage_pct"`
	MaxTransactionCpuUsage uint32 `json:"max_transaction_cpu_usage"`
	MinTransactionCpuUsage uint32 `json:"min_transaction_cpu_usage"`

	MaxTrxLifetime              uint32 `json:"max_transaction_lifetime"`
	DeferredTrxExpirationWindow uint32 `json:"deferred_trx_expiration_window"`
	MaxTrxDelay                 uint32 `json:"max_transaction_delay"`
	MaxInlineActionSize         uint32 `json:"max_inline_action_size"`
	MaxInlineActionDepth        uint16 `json:"max_inline_action_depth"`
	MaxAuthorityDepth           uint16 `json:"max_authority_depth"`
}

// TODO: Add validation logic
func (c ChainConfig) Validate() error {
	if c.TargetBlockNetUsagePct > uint32(Percent100) {
		return errors.New("target block net usage percentage cannot exceed 100%")
	}

	if c.TargetBlockNetUsagePct < uint32(Percent1/10) {
		return errors.New("target block net usage percentage must be at least 0.1%")
	}

	if c.TargetBlockCpuUsagePct > uint32(Percent100) {
		return errors.New("target block cpu usage percentage cannot exceed 100%")
	}

	if c.TargetBlockCpuUsagePct < uint32(Percent1/10) {
		return errors.New("target block cpu usage percentage must be at least 0.1%")
	}

	if c.MaxTransactionNetUsage >= uint32(c.MaxBlockNetUsage) {
		return errors.New("max transaction net usage must be less than max block net usage")
	}

	if c.MaxTransactionCpuUsage >= c.MaxBlockCpuUsage {
		return errors.New("max transaction cpu usage must be less than max block cpu usage")
	}

	if c.BasePerTransactionNetUsage >= c.MaxTransactionNetUsage {
		return errors.New("base net usage per transaction must be less than the max transaction net usage")
	}

	if c.MaxTransactionNetUsage-c.BasePerTransactionNetUsage < MinNetUsageDeltaBetweenBaseAndMaxForTrx {
		return fmt.Errorf("max transaction net usage must be at least %v bytes larger than base net usage per transaction", MinNetUsageDeltaBetweenBaseAndMaxForTrx)
	}

	if c.ContextFreeDiscountNetUsageDen <= 0 {
		return errors.New("net usage discount ratio for context free data cannot have a 0 denominator")
	}

	if c.ContextFreeDiscountNetUsageNum > c.ContextFreeDiscountNetUsageDen {
		return errors.New("net usage discount ratio for context free data cannot exceed 1")
	}

	if c.MinTransactionCpuUsage > c.MaxTransactionCpuUsage {
		return errors.New("min transaction cpu usage cannot exceed max transaction cpu usage")
	}

	if c.MaxTransactionCpuUsage >= (c.MaxBlockCpuUsage - c.MinTransactionCpuUsage) {
		return errors.New("max transaction cpu usage must be at less than the difference between the max block cpu usage and the min transaction cpu usage")
	}

	if c.MaxAuthorityDepth < 1 {
		return errors.New("max authority depth should be at least 1")
	}

	return nil
}
