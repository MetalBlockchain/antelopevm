package service

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/name"
)

type GetCurrencyStatsRequest struct {
	Code   name.AccountName `json:"code"`
	Symbol string           `json:"symbol"`
}

type CurrencyStats struct {
	Supply    core.Asset       `json:"supply"`
	MaxSupply core.Asset       `json:"max_supply"`
	Issuer    name.AccountName `json:"issuer"`
}

type GetCurrencyStatsResponse map[string]CurrencyStats
