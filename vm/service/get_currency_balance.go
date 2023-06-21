package service

import "github.com/MetalBlockchain/antelopevm/core/name"

type GetCurrencyBalanceRequest struct {
	Code    name.AccountName `json:"code"`
	Account name.AccountName `json:"account"`
	Symbol  string           `json:"symbol"`
}
