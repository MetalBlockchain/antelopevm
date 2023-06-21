package service

import (
	"github.com/MetalBlockchain/antelopevm/core/name"
)

type GetTableRowsRequest struct {
	Code  name.AccountName `json:"code"`
	Scope name.ScopeName   `json:"scope"`
	Table name.TableName   `json:"table"`
}

type TableRow struct {
	Data  map[string]interface{} `json:"data"`
	Payer string                 `json:"payer"`
}

type GetTableRowsResponse struct {
	More bool       `json:"more"`
	Rows []TableRow `json:"rows"`
}
